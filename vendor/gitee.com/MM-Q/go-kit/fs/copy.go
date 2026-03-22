// Package fs 提供文件系统操作工具，支持文件和目录的复制、移动、检查等功能。
//
// 复制功能支持的场景
//
// Copy 和 CopyEx 函数支持以下复制场景：
//
// 【文件复制】
//   - Copy("a.txt", "b.txt")           → 创建 b.txt（精确路径模式）
//   - Copy("a.txt", "existingDir")     → 创建 existingDir/a.txt（自动追加文件名）
//   - Copy("a.txt", "existingDir/")    → 创建 existingDir/a.txt（自动追加文件名）
//   - Copy("a.txt", "newDir/b.txt")    → 创建 newDir/b.txt（自动创建父目录）
//
// 【目录复制】
//   - Copy("dirA", "dirB")             → 创建 dirB/（dirB 不存在时）
//   - Copy("dirA", "existingDir")      → 创建 existingDir/dirA/（自动追加目录名）
//   - Copy("dirA", "existingDir/")     → 创建 existingDir/dirA/（自动追加目录名）
//   - Copy("dirA", "newDir/subDir")    → 创建 newDir/subDir/（自动创建父目录）
//
// 【特殊类型】
//   - 符号链接：Linux/macOS 保留链接，Windows 当作普通文件复制
//   - 特殊文件：设备文件、命名管道等（仅 Unix 系统）
//
// 移动功能
//
// Move 和 MoveEx 函数（见 move.go）支持文件和目录移动，移动规则与复制相同。
// 注意：移动操作通过复制+删除实现，支持跨文件系统移动。
//
// 智能路径处理规则：
//   - 如果目标路径是已存在的目录，自动追加源文件名/目录名
//   - 如果目标路径不存在或不是目录，使用精确路径模式
//
// 覆盖控制：
//   - Copy() / Move() 函数默认不允许覆盖已存在的目标
//   - CopyEx() / MoveEx() 函数可通过 overwrite 参数控制是否允许覆盖
//
// 原子性保证：
//   - 文件复制使用临时文件 + os.Rename 保证原子性
//   - 覆盖时先备份原文件，失败时自动恢复
//
// 【已知限制】
//   - Windows 上复制指向目录的符号链接会失败
//     原因：Windows 符号链接需要管理员权限，当符号链接指向目录时，
//     内部调用 copyFile 会因目标不是普通文件而返回错误。
//     建议：Windows 用户使用快捷方式而非符号链接，或手动处理此类特殊情况。
package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gitee.com/MM-Q/go-kit/pool"
)

// Copy 通用复制函数，自动判断源路径类型并调用相应的复制函数
// 支持复制普通文件、目录、符号链接和特殊文件 (设备文件、命名管道等)
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径（支持文件、目录，自动创建父目录）
//
// 返回:
//   - error: 复制失败时返回错误，如果目标已存在则返回错误
func Copy(src, dst string) error {
	return CopyEx(src, dst, false)
}

// CopyEx 通用复制函数 (可控制是否覆盖)，自动判断源路径类型并调用相应的复制函数
// 支持复制普通文件、目录、符号链接和特殊文件 (设备文件、命名管道等)
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径（支持文件、目录，自动创建父目录）
//   - overwrite: 是否允许覆盖已存在的目标文件/目录
//
// 返回:
//   - error: 复制失败时返回错误
//
// 智能路径处理:
//   - 如果 dst 是已存在的目录，会自动追加源文件名/目录名
//   - 例如: Copy("a.txt", "existingDir") → 创建 existingDir/a.txt
//   - 例如: Copy("dirA", "existingDir") → 创建 existingDir/dirA/
func CopyEx(src, dst string, overwrite bool) (err error) {
	// 捕获 panic 并转换为错误
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("copy operation panicked: %v", r)
		}
	}()

	// 验证路径并获取绝对路径
	srcAbs, dstAbs, err := validateAndResolvePaths(src, dst)
	if err != nil {
		return err
	}

	// 验证路径之间的关系（检查路径不相同）
	if err := validatePathRelations(srcAbs, dstAbs, false); err != nil {
		return err
	}

	// 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名
	dstAbs = resolveDestinationPathAbs(srcAbs, dstAbs)

	// 调用内部复制函数
	return copyExInternal(srcAbs, dstAbs, overwrite)
}

// copyExInternal 内部复制函数，接受已验证的绝对路径
// 用于避免重复验证，供 MoveEx 等内部函数调用
//
// 参数:
//   - srcAbs: 已验证的源绝对路径
//   - dstAbs: 已验证的目标绝对路径（已通过智能路径处理）
//   - overwrite: 是否允许覆盖已存在的目标文件/目录
//
// 返回:
//   - error: 复制失败时返回错误
func copyExInternal(srcAbs, dstAbs string, overwrite bool) error {
	// 获取源路径信息（使用 Lstat 避免跟随符号链接）
	srcInfo, localErr := os.Lstat(srcAbs)
	if localErr != nil {
		return fmt.Errorf("failed to get source info '%s': %w", srcAbs, localErr)
	}

	// 根据源路径类型调用相应的复制函数
	if srcInfo.IsDir() {
		return copyDir(srcAbs, dstAbs, overwrite)
	} else {
		// 处理所有文件类型（普通文件、符号链接、特殊文件等）
		return copyFileRouter(srcAbs, dstAbs, srcInfo, overwrite)
	}
}

// resolveDestinationPathAbs 解析目标路径，实现智能路径追加
// 如果目标是已存在的目录，自动追加源文件名/目录名
//
// 参数:
//   - srcAbs: 源绝对路径
//   - dstAbs: 目标绝对路径
//
// 返回:
//   - string: 调整后的目标绝对路径
func resolveDestinationPathAbs(srcAbs, dstAbs string) string {
	// 清理路径，移除末尾的路径分隔符
	dstAbs = strings.TrimRight(dstAbs, string(filepath.Separator))

	// 检查目标是否是已存在的目录
	dstInfo, err := os.Lstat(dstAbs)
	if err != nil || !dstInfo.IsDir() {
		// 目标不存在或不是目录，保持原样（精确路径模式）
		return dstAbs
	}

	// 目标是已存在的目录，自动追加源文件名/目录名
	srcBase := filepath.Base(srcAbs)
	return filepath.Join(dstAbs, srcBase)
}

// validateAndResolvePaths 验证路径并返回绝对路径
// 检查路径是否为空，获取并清理绝对路径
//
// 参数:
//   - src: 源路径
//   - dst: 目标路径
//
// 返回:
//   - srcAbs: 源绝对路径
//   - dstAbs: 目标绝对路径
//   - error: 验证失败时返回错误
func validateAndResolvePaths(src, dst string) (srcAbs, dstAbs string, err error) {
	// 检查路径是否为空
	if src == "" || dst == "" {
		return "", "", fmt.Errorf("source and destination paths cannot be empty")
	}

	// 获取并清理绝对路径
	srcAbs, err = filepath.Abs(filepath.Clean(src))
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path for source '%s': %w", src, err)
	}
	dstAbs, err = filepath.Abs(filepath.Clean(dst))
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path for destination '%s': %w", dst, err)
	}

	return srcAbs, dstAbs, nil
}

// validatePathRelations 验证路径之间的关系
// 检查源路径和目标路径不相同、是否复制到子目录
//
// 参数:
//   - srcAbs: 源绝对路径
//   - dstAbs: 目标绝对路径
//   - checkSubdir: 是否检查目录复制到子目录的情况
//
// 返回:
//   - error: 验证失败时返回错误
func validatePathRelations(srcAbs, dstAbs string, checkSubdir bool) error {
	// 检查源路径和目标路径不相同
	if srcAbs == dstAbs {
		return fmt.Errorf("source and destination paths cannot be the same")
	}

	// 检查是否尝试将目录复制到自己的子目录中
	if checkSubdir {
		if strings.HasPrefix(dstAbs+string(filepath.Separator), srcAbs+string(filepath.Separator)) {
			return fmt.Errorf("cannot copy directory '%s' to its subdirectory '%s'", srcAbs, dstAbs)
		}
	}

	return nil
}

// handleBackupAndRestore 统一的备份恢复处理函数
// 用于在覆盖操作前创建备份，失败时恢复备份
//
// 参数:
//   - dst: 目标路径
//   - overwrite: 是否允许覆盖
//
// 返回:
//   - backupPath: 备份文件路径（如果创建了备份）
//   - error: 处理失败时返回错误
func handleBackupAndRestore(dst string, overwrite bool) (string, error) {
	// 检查目标是否存在
	if _, err := os.Lstat(dst); err != nil {
		// 目标不存在，无需备份
		return "", nil
	}

	// 目标存在
	if !overwrite {
		return "", fmt.Errorf("destination '%s' already exists", dst)
	}

	// 允许覆盖，创建备份（使用 pid + 纳秒时间戳确保并发安全）
	backupPath := dst + ".backup." + fmt.Sprintf("%d.%d", os.Getpid(), time.Now().UnixNano())
	if err := os.Rename(dst, backupPath); err != nil {
		return "", fmt.Errorf("failed to backup existing '%s': %w", dst, err)
	}

	return backupPath, nil
}

// restoreBackup 恢复备份文件
// 在操作失败时调用，尽力恢复原始文件
//
// 参数:
//   - dst: 目标路径
//   - backupPath: 备份文件路径
func restoreBackup(dst, backupPath string) {
	if backupPath != "" {
		_ = os.Rename(backupPath, dst) // 忽略恢复错误，尽力而为
	}
}

// cleanupBackup 清理备份文件/目录
// 在操作成功时调用，删除不再需要的备份
//
// 参数:
//   - backupPath: 备份文件或目录路径
func cleanupBackup(backupPath string) {
	if backupPath != "" {
		_ = os.RemoveAll(backupPath) // 忽略删除错误，支持文件和目录
	}
}

// copyFile 内部复制文件逻辑
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - srcAbs: 源文件绝对路径
//   - dstAbs: 目标文件绝对路径
//   - srcInfo: 源文件信息（包含 Mode、Size 等）
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyFile(srcAbs, dstAbs string, srcInfo os.FileInfo, overwrite bool) error {
	// 注意：路径验证已在 CopyEx 入口处统一完成，此处无需重复验证

	// 安全覆盖机制：处理已存在的目标文件
	backupPath, err := handleBackupAndRestore(dstAbs, overwrite)
	if err != nil {
		return err
	}

	// 打开源文件
	in, err := os.Open(srcAbs)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s': %w", srcAbs, err)
	}
	defer func() { _ = in.Close() }()

	// 检查是否为普通文件
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("source '%s' is not a regular file", srcAbs)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dstAbs)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建临时文件（与目标同目录，保证 rename 原子性）
	// 使用 pid + 纳秒时间戳确保同一进程内并发安全
	tmp := dstAbs + ".tmp." + fmt.Sprintf("%d.%d", os.Getpid(), time.Now().UnixNano())
	out, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_EXCL, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create temporary file '%s': %w", tmp, err)
	}

	// 统一清理资源
	success := false
	defer func() {
		if out != nil {
			_ = out.Close()
			out = nil // 防止重复关闭
		}
		if !success {
			_ = os.Remove(tmp) // 清理临时文件
		}
	}()

	// 根据文件大小选择处理方式
	if srcInfo.Size() == 0 {
		// 空文件：跳过数据复制，直接进行后续操作
	} else {
		// 非空文件：使用缓冲区进行数据拷贝
		bufSize := pool.CalculateBufferSize(srcInfo.Size())
		buf := pool.GetByteCap(bufSize)
		defer pool.PutByte(buf)

		if _, err := io.CopyBuffer(out, in, buf); err != nil {
			return fmt.Errorf("failed to copy data from '%s' to '%s': %w", srcAbs, tmp, err)
		}

		// 强制刷盘，确保数据持久化（仅对非空文件）
		if err := out.Sync(); err != nil {
			return fmt.Errorf("failed to sync temporary file '%s': %w", tmp, err)
		}

	}

	// 在重命名前关闭文件句柄(Windows要求)
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file '%s': %w", tmp, err)
	}
	out = nil // 标记为已关闭

	// 原子重命名
	if err := os.Rename(tmp, dstAbs); err != nil {
		// 复制失败，恢复备份文件
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to rename temporary file '%s' to '%s': %w", tmp, dstAbs, err)
	}

	// 复制成功，删除备份文件
	cleanupBackup(backupPath)

	success = true
	return nil
}

// copySymlink 复制符号链接
// Windows 平台：当作普通文件复制（因为 Windows 主要使用快捷方式而非符号链接）
// 非 Windows 平台：创建相同的符号链接
//
// 参数:
//   - srcAbs: 源符号链接绝对路径
//   - dstAbs: 目标绝对路径
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copySymlink(srcAbs, dstAbs string, overwrite bool) error {
	// Windows 平台：当作普通文件复制
	// 注意：Windows 符号链接可能指向目录，需要先获取目标信息
	if runtime.GOOS == "windows" {
		// 获取符号链接目标信息
		targetInfo, err := os.Stat(srcAbs)
		if err != nil {
			return fmt.Errorf("failed to get symlink target info '%s': %w", srcAbs, err)
		}
		return copyFile(srcAbs, dstAbs, targetInfo, overwrite)
	}

	// 非 Windows 平台：创建符号链接
	// 安全覆盖机制：处理已存在的目标符号链接
	backupPath, err := handleBackupAndRestore(dstAbs, overwrite)
	if err != nil {
		return err
	}

	// 读取符号链接的目标
	target, err := os.Readlink(srcAbs)
	if err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to read symlink '%s': %w", srcAbs, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dstAbs)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建符号链接
	if err := os.Symlink(target, dstAbs); err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to create symlink '%s' -> '%s': %w", dstAbs, target, err)
	}

	// 创建成功，删除备份
	cleanupBackup(backupPath)
	return nil
}

// copySpecialFile 复制特殊文件（设备文件、命名管道、套接字等）
// 对于特殊文件，只创建一个具有相同权限模式的空文件
//
// 参数:
//   - srcAbs: 源特殊文件绝对路径
//   - dstAbs: 目标文件绝对路径
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copySpecialFile(srcAbs, dstAbs string, overwrite bool) error {
	// 安全覆盖机制：处理已存在的目标文件
	backupPath, err := handleBackupAndRestore(dstAbs, overwrite)
	if err != nil {
		return err
	}

	// 获取源文件信息
	srcInfo, err := os.Lstat(srcAbs)
	if err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to get source file info '%s': %w", srcAbs, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dstAbs)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建一个空文件，保持相同的权限模式
	file, err := os.OpenFile(dstAbs, os.O_CREATE|os.O_WRONLY, srcInfo.Mode())
	if err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to create special file '%s': %w", dstAbs, err)
	}

	// 立即关闭文件
	if err := file.Close(); err != nil {
		// 关闭失败，删除创建的文件并恢复备份
		_ = os.Remove(dstAbs)
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to close special file '%s': %w", dstAbs, err)
	}

	// 创建成功，删除备份
	cleanupBackup(backupPath)
	return nil
}

// copyFileRouter 文件复制路由函数
// 根据文件类型调用相应的复制函数
//
// 参数:
//   - srcAbs: 源文件绝对路径
//   - dstAbs: 目标文件绝对路径
//   - srcInfo: 源文件信息（包含 Mode、Size 等）
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copyFileRouter(srcAbs, dstAbs string, srcInfo os.FileInfo, overwrite bool) error {
	switch {
	case srcInfo.Mode().IsRegular():
		// 普通文件
		return copyFile(srcAbs, dstAbs, srcInfo, overwrite)

	case srcInfo.Mode()&os.ModeSymlink != 0:
		// 符号链接
		return copySymlink(srcAbs, dstAbs, overwrite)

	default:
		// 其他特殊文件（设备文件、命名管道、套接字等）
		return copySpecialFile(srcAbs, dstAbs, overwrite)
	}
}

// copyDir 内部复制目录逻辑
// 用于递归复制整个目录，保持文件权限和目录结构
//
// 参数:
//   - srcAbs: 源目录绝对路径
//   - dstAbs: 目标目录绝对路径
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyDir(srcAbs, dstAbs string, overwrite bool) error {
	// 单独调用子目录检查（避免重复基础验证）
	if err := validatePathRelations(srcAbs, dstAbs, true); err != nil {
		return err
	}

	// 获取源目录信息
	srcInfo, err := os.Stat(srcAbs)
	if err != nil {
		return fmt.Errorf("failed to get source directory info '%s': %w", srcAbs, err)
	}

	// 检查源路径是否为目录
	if !srcInfo.IsDir() {
		return fmt.Errorf("source '%s' is not a directory", srcAbs)
	}

	// 安全覆盖机制：处理已存在的目标目录
	backupPath, err := handleBackupAndRestore(dstAbs, overwrite)
	if err != nil {
		return err
	}

	// 创建目标目录，使用合适的权限（至少需要写权限以便后续操作）
	dirMode := srcInfo.Mode()
	if dirMode&0o200 == 0 {
		// 如果源目录没有写权限，临时添加写权限以便复制操作
		dirMode |= 0o200
	}
	if err := os.MkdirAll(dstAbs, dirMode); err != nil {
		restoreBackup(dstAbs, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dstAbs, err)
	}

	// 遍历源目录
	copyErr := filepath.WalkDir(srcAbs, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path '%s': %w", path, err)
		}

		// 计算相对路径
		relPath, err := filepath.Rel(srcAbs, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for '%s': %w", path, err)
		}

		// 跳过根目录本身
		if relPath == "." {
			return nil
		}

		// 构建目标路径
		dstPath := filepath.Join(dstAbs, relPath)

		// 处理目录
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("failed to get directory info '%s': %w", path, err)
			}
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", dstPath, err)
			}
			return nil
		}

		// 处理所有文件类型（普通文件、符号链接、特殊文件等）
		// 获取文件信息传递给 copyFileRouter
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info '%s': %w", path, err)
		}
		return copyFileRouter(path, dstPath, info, overwrite)
	})

	// 处理复制结果
	if copyErr != nil {
		// 复制失败，清理已复制的内容并恢复备份
		_ = os.RemoveAll(dstAbs) // 清理部分复制的内容
		restoreBackup(dstAbs, backupPath)
		return copyErr
	}

	// 复制成功，恢复目录的原始权限（如果之前临时修改了权限）
	if srcInfo.Mode() != dirMode {
		_ = os.Chmod(dstAbs, srcInfo.Mode()) // 忽略权限恢复错误
	}

	// 复制成功，删除备份目录
	cleanupBackup(backupPath)
	return nil
}
