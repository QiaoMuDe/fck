import os
import subprocess
import sys
import argparse
import time
import zipfile
from datetime import datetime, timezone
import platform
import threading
import concurrent.futures
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass


############################### 以下为可配置的变量 #################################
# 基础输出文件名(指定时无需包含后缀)同时也是注入的appName
BASE_OUTPUT_NAME = "myapp"
# 默认输出目录
DEFAULT_OUTPUT_DIR = "output"
# 默认入口文件的位置
DEFAULT_ENTRY_FILE = "./main.go"
# 默认的 Go 编译器, 使用全局 PATH 中的 go
DEFAULT_GO_COMPILER = "go"
# 默认不使用 vendor 克隆依赖
DEFAULT_USE_VENDOR = False
# 默认在构建阶段使用 vendor 目录
DEFAULT_USE_VENDOR_IN_BUILD = False
# 是否注入git信息, 默认为True
DEFAULT_INJECT_GIT_INFO = True
# 是否使用简单文件名格式, 默认为False
DEFAULT_SIMPLE_NAME = False
# 是否仅构建当前平台, 默认为False
DEFAULT_CURRENT_PLATFORM_ONLY = False
# 是否将构建成功的可执行文件打包为zip文件, 默认为False
DEFAULT_PACKAGE_ZIP = False
# 批量构建时默认的并发线程数
DEFAULT_CONCURRENCY = max(1, os.cpu_count() - 1)  # 使用CPU核心数-1
# 批量构建时默认的超时时间(秒)
DEFAULT_TIMEOUT = 1800
# 默认环境变量字典
DEFAULT_ENV_VARS = {
    "GOPROXY": "https://goproxy.cn,https://goproxy.io,direct",  # Go 代理地址, 默认为 goproxy.cn 和 goproxy.io
    "CGO_ENABLED": "0",  # 是否启用 CGO 编译, 0为禁用, 1为启用
    "CC": "gcc",  # 默认使用gcc编译器
    "CXX": "g++",  # 默认使用g++编译器
}
####################################################################################


############################### 以下为内部使用的变量 ###############################
# Git信息缓存字典
_git_info_cache = {"version": None, "commit": None, "commit_time": None, "status": None}
# 支持的平台列表
SUPPORTED_PLATFORMS = ["windows", "linux", "darwin"]
# 平台简写映射
PLATFORM_SHORTCUTS = {"w": "windows", "l": "linux", "d": "darwin"}
# 支持的架构列表
SUPPORTED_ARCHITECTURES = ["amd64"]
# 架构简写映射
ARCHITECTURE_SHORTCUTS = {"a64": "amd64"}
# 定义颜色转义字符
RED_BOLD = "\033[1;31m"  # 红色加粗
GREEN_BOLD = "\033[1;32m"  # 绿色加粗
RESET = "\033[0m"  # 重置颜色
# 默认构建时的链接器标志
DEFAULT_LDFLAGS = "-s -w"
# 启用git信息注入时的链接器标志模板
LD_FLAGS_TEMPLATE = "-X 'gitee.com/MM-Q/verman.appName={app_name}' -X 'gitee.com/MM-Q/verman.gitVersion={git_version}' -X 'gitee.com/MM-Q/verman.gitCommit={git_commit}' -X 'gitee.com/MM-Q/verman.gitCommitTime={commit_time}' -X 'gitee.com/MM-Q/verman.buildTime={build_time}' -X 'gitee.com/MM-Q/verman.gitTreeState={tree_state}' -s -w"
####################################################################################


# 函数定义 #
def print_success(message):
    """打印成功信息"""
    print(f"{GREEN_BOLD}ok: {message}{RESET}")


def print_error(message):
    """打印错误信息"""
    print(f"{RED_BOLD}error: {message}{RESET}")


def check_go_installed(go_compiler):
    """检查指定的 Go 编译器是否可用"""
    try:
        print_success("正在检查 Go 编译器是否可用...")
        subprocess.run(
            [go_compiler, "version"],
            capture_output=True,
            text=True,
            check=True,
            encoding="utf-8",
        )
        return True
    except subprocess.CalledProcessError:
        print_error(
            f"未检测到 {go_compiler} 编译器, 请确保已安装并添加到 PATH 中, 或者指定正确的路径。"
        )
        return False


def check_go_mod_file():
    """检查当前目录是否存在 go.mod 文件"""
    print_success("正在检查当前目录是否存在 go.mod 文件...")
    if os.path.exists("go.mod"):
        return True
    else:
        print_error("当前目录下未找到 go.mod 文件。")
        return False


def check_entry_file(entry_file):
    """检查指定的入口文件是否存在"""
    print_success(f"正在检查指定的入口文件 {entry_file} 是否存在...")
    if os.path.exists(entry_file):
        return True
    else:
        print_error(f"未找到入口文件 {entry_file}。")
        return False


def run_go_mod_vendor(go_compiler):
    """执行 go mod vendor 克隆依赖"""
    try:
        print_success("正在执行 go mod vendor 克隆依赖...")
        subprocess.run(
            [go_compiler, "mod", "vendor"],
            capture_output=True,
            text=True,
            check=True,
            encoding="utf-8",
        )
        print_success("go mod vendor 执行成功, 依赖已克隆到 vendor 目录。")
    except subprocess.CalledProcessError as e:
        print_error("go mod vendor 执行失败：")
        print_error(e.stderr.strip())
        sys.exit(1)


def run_go_mod_tidy(go_compiler, use_vendor):
    """执行 go mod tidy 并处理输出"""
    command = [go_compiler, "mod", "tidy"]
    if use_vendor:
        command.extend(["-v"])
    try:
        print_success("正在执行 go mod tidy...")
        subprocess.run(
            command, capture_output=True, text=True, check=True, encoding="utf-8"
        )
    except subprocess.CalledProcessError as e:
        print_error("go mod tidy 执行失败：")
        print_error(e.stderr.strip())
        sys.exit(1)


def run_code_check(go_compiler):
    """执行代码检查, 使用go vet"""
    command = [go_compiler, "vet", "./..."]
    try:
        print_success("正在执行 go vet 检查代码...")
        subprocess.run(
            command, capture_output=True, text=True, check=True, encoding="utf-8"
        )
    except subprocess.CalledProcessError as e:
        print_error("go vet 检查失败：")
        print_error(e.stderr.strip())
        sys.exit(1)


def run_gofmt(go_compiler):
    """执行 gofmt -w . 格式化代码"""
    try:
        print_success("正在执行 gofmt 格式化代码...")
        subprocess.run(
            [go_compiler, "fmt", "./..."],
            capture_output=True,
            text=True,
            check=True,
            encoding="utf-8",
        )
    except subprocess.CalledProcessError as e:
        print_error("代码格式化失败：")
        print_error(e.stderr.strip())
        sys.exit(1)


# 数据类封装构建配置参数
@dataclass
class BuildConfig:
    go_compiler: str
    output_file: str
    entry_file: str
    ldflags: str
    use_vendor_in_build: bool
    is_batch: bool = False
    args: argparse.Namespace = None


def build_go_app(
    config: BuildConfig,
):
    """组装并执行构建命令"""
    command = [
        config.go_compiler,
        "build",
        "-o",
        config.output_file,
        "-ldflags",
        config.ldflags,
    ]
    if config.use_vendor_in_build:
        # 检查 vendor 目录是否存在
        if not os.path.exists("vendor"):
            print_error("vendor 目录不存在, 无法使用 -mod=vendor 选项。")
            return False
        command.extend(["-mod=vendor"])
    command.append(config.entry_file)
    try:
        # 强制设置环境变量
        env = os.environ.copy()
        # 添加默认环境变量
        env.update(DEFAULT_ENV_VARS)

        # 为arm64架构设置特定的交叉编译工具链
        if env.get("GOARCH") == "arm64":
            env["CC"] = "aarch64-linux-gnu-gcc"
            env["CXX"] = "aarch64-linux-gnu-g++"
        # 添加自定义环境变量
        if config.args and hasattr(config.args, "env") and config.args.env:
            for env_var in config.args.env:
                if "=" in env_var:
                    key, value = env_var.split("=", 1)
                    env[key] = value
        # 从输出文件名中提取平台信息
        if "_windows_" in config.output_file:
            env["GOOS"] = "windows"
        elif "_linux_" in config.output_file:
            env["GOOS"] = "linux"
        elif "_darwin_" in config.output_file:
            env["GOOS"] = "darwin"
        else:
            env["GOOS"] = platform.system().lower()

        # 从输出文件名中提取架构信息
        if "_amd64" in config.output_file:
            env["GOARCH"] = "amd64"
        elif "_arm64" in config.output_file:
            env["GOARCH"] = "arm64"
        elif "_386" in config.output_file:
            env["GOARCH"] = "386"
        elif "_arm" in config.output_file:
            env["GOARCH"] = "arm"
        else:
            machine = platform.machine().lower()
            # 自动转换x86_64为amd64
            if machine == "x86_64":
                machine = "amd64"
            env["GOARCH"] = machine

        # 使用指定的链接器标志和环境变量进行构建
        subprocess.run(
            command,
            capture_output=True,
            text=True,
            check=True,
            env=env,
            encoding="utf-8",
        )

        if not config.is_batch:
            print_success(f"构建成功, 输出文件：{config.output_file}")
        return True
    except subprocess.CalledProcessError as e:
        print_error("构建失败：")
        print_error(e.stderr.strip())
        return False


def zip_executable(output_file, zip_file, is_batch=False):
    """将构建成功的可执行文件打包到 ZIP 文件中"""
    try:
        if not is_batch:
            print_success(f"{output_file} --> {zip_file}")
        with zipfile.ZipFile(zip_file, "w", zipfile.ZIP_DEFLATED) as zipf:
            zipf.write(output_file)

        try:
            os.remove(output_file)
        except Exception as e:
            print_error(f"删除源文件 {output_file} 失败: {str(e)}")
    except Exception as e:
        print_error(f"打包到 {zip_file} 失败：{str(e)}")


def batch_build(args):
    """批量构建所有支持的平台和架构组合"""
    print_success("开始批量构建所有支持的平台和架构组合...")
    total_start_time = time.time()
    success_count = 0
    fail_count = 0
    skip_count = 0
    total_tasks = len(SUPPORTED_PLATFORMS) * len(SUPPORTED_ARCHITECTURES)
    print_success(f"总任务数: {total_tasks}")
    lock = threading.Lock()

    # 执行构建前的检查工作
    print_success("开始检查构建环境...")
    if not pre_build_checks(args.go_compiler, args.entry, args.use_vendor):
        sys.exit(1)

    # 根据参数注入 Git 信息
    if args.git:
        print_success("正在获取 Git 信息...")
        if _git_info_cache["version"] is None:
            print_error("Git 信息获取失败, 请检查 Git 环境是否正确配置。")
            sys.exit(1)

    # 创建临时args对象用于批量构建
    batch_args = argparse.Namespace(**vars(args))

    def build_task(system, architecture):
        nonlocal success_count, fail_count, skip_count
        # 跳过不支持的darwin/386和darwin/arm组合
        if system == "darwin" and architecture in ("386", "arm"):
            with lock:
                skip_count += 1
                print_success(f"跳过不支持的平台/架构组合: {system}/{architecture}")
            return

        # 如果启用了仅构建当前平台且平台不一致则跳过
        if args.current_platform_only and system != platform.system().lower():
            with lock:
                skip_count += 1
                print_success(f"跳过非当前平台: {system}/{architecture}")
            return

        batch_args.platform = system
        batch_args.arch = architecture

        # 检查架构组合是否支持
        if parse_arguments() is None:
            return

        # 如果使用简单文件名格式, 则生成输出文件名时不包含系统和架构信息
        if args.simple_name:
            base_name = f"{BASE_OUTPUT_NAME}"
            git_version = _git_info_cache["version"] if args.git else None
            output_file = generate_output_file_name(base_name, system, git_version)
        else:
            base_name = f"{BASE_OUTPUT_NAME}_{system}_{architecture}"
            git_version = _git_info_cache["version"] if args.git else None
            output_file = generate_output_file_name(base_name, system, git_version)

        # 生成zip文件名
        if args.zip:
            if args.zip_file is None:
                git_version = _git_info_cache["version"] if args.git else None
                zip_file = generate_zip_file_name(
                    BASE_OUTPUT_NAME, system, architecture, git_version
                )
            else:
                zip_file = args.zip_file
        else:
            zip_file = None

        try:
            # 执行构建
            build_result = single_build(
                args, system, architecture, output_file, zip_file
            )

            with lock:
                if build_result:
                    success_count += 1
                else:
                    fail_count += 1
                completed_count = success_count + fail_count
                print_success(
                    f"已完成 {completed_count}/{total_tasks} 个任务 (成功 {success_count} 个, 失败 {fail_count} 个, 跳过 {skip_count} 个)"
                )
                pass
        except Exception as e:
            with lock:
                fail_count += 1
                completed_count = success_count + fail_count
                print_success(
                    f"已完成 {completed_count}/{total_tasks} 个任务 (成功 {success_count} 个, 失败 {fail_count} 个, 跳过 {skip_count} 个)"
                )
                print_error(f"构建 {system}/{architecture} 时发生异常: {str(e)}")

    # 创建线程池
    with ThreadPoolExecutor(max_workers=args.max_workers) as executor:
        # 提交所有构建任务
        futures = []
        for system in SUPPORTED_PLATFORMS:
            for architecture in SUPPORTED_ARCHITECTURES:
                futures.append(executor.submit(build_task, system, architecture))

        # 等待所有任务完成, 设置超时时间为30分钟
        try:
            for future in futures:
                future.result(timeout=args.timeout)
        except concurrent.futures.TimeoutError:
            print_error("任务执行超时, 强制终止线程池")
            executor._threads.clear()
            concurrent.futures.thread._threads_queues.clear()
            fail_count += len([f for f in futures if not f.done()])

    total_elapsed_time = time.time() - total_start_time
    print_success(
        f"批量构建完成, 成功 {success_count} 个, 失败 {fail_count} 个, 跳过 {skip_count} 个"
    )
    print_success(f"总耗时: {total_elapsed_time:.2f} 秒")


def single_build(args, system, architecture, output_file, zip_file):
    """执行单个平台和架构的构建"""
    try:
        # 设置环境变量
        env = os.environ.copy()
        env["GOOS"] = system
        env["GOARCH"] = architecture

        # 获取构建时间
        build_time = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

        # 处理Git信息
        ldflags = args.ldflags
        if args.git:
            if _git_info_cache["version"] is None:
                return False
            ldflags = LD_FLAGS_TEMPLATE.format(
                app_name=BASE_OUTPUT_NAME,
                git_version=_git_info_cache["version"],
                git_commit=_git_info_cache["commit"],
                commit_time=_git_info_cache["commit_time"],
                build_time=build_time,
                tree_state=_git_info_cache["status"],
            )

        # 执行构建
        build_config = BuildConfig(
            go_compiler=args.go_compiler,
            output_file=output_file,
            entry_file=args.entry,
            ldflags=ldflags,
            use_vendor_in_build=args.use_vendor_in_build,
            is_batch=True,
            args=args,
        )

        # 构建
        build_result = build_go_app(build_config)

        # 压缩
        if build_result and args.zip and zip_file:
            zip_executable(output_file, zip_file, True)

        return build_result
    except Exception as e:
        print_error(f"构建 {system}/{architecture} 失败: {str(e)}")
        return False


def get_git_info():
    """获取 Git 版本信息"""
    global _git_info_cache

    # 如果缓存不存在或无效, 则尝试获取git信息
    if _git_info_cache["version"] is None:
        try:
            git_version = subprocess.run(
                ["git", "describe", "--tags", "--always", "--dirty"],
                capture_output=True,
                text=True,
                check=True,
                timeout=10,
                encoding="utf-8",
            ).stdout.strip()
            git_commit = subprocess.run(
                ["git", "rev-parse", "--short", "HEAD"],
                capture_output=True,
                text=True,
                check=True,
                timeout=10,
                encoding="utf-8",
            ).stdout.strip()
            git_commit_time = subprocess.run(
                ["git", "log", "-1", "--format=%cd", "--date=iso"],
                capture_output=True,
                text=True,
                check=True,
                timeout=10,
                encoding="utf-8",
            ).stdout.strip()
            git_status = (
                "dirty"
                if subprocess.run(
                    ["git", "status", "--porcelain"],
                    capture_output=True,
                    text=True,
                    timeout=10,
                    encoding="utf-8",
                ).stdout.strip()
                else "clean"
            )
            format_time = datetime.strptime(
                git_commit_time, "%Y-%m-%d %H:%M:%S %z"
            ).strftime("%Y-%m-%dT%H:%M:%SZ")
            _git_info_cache["version"] = git_version
            _git_info_cache["commit"] = git_commit
            _git_info_cache["commit_time"] = format_time
            _git_info_cache["status"] = git_status
        except subprocess.CalledProcessError:
            print_error("警告: 无法获取 Git 版本信息, 可能是当前目录不是 Git 仓库。")
            return None
        except subprocess.TimeoutExpired:
            print_error("获取 Git 版本信息超时。")
            return None

    return (
        _git_info_cache["version"],
        _git_info_cache["commit"],
        _git_info_cache["commit_time"],
        _git_info_cache["status"],
    )


def generate_output_file_name(base_name, system, git_version=None):
    """根据操作系统生成默认输出文件名, 可选的git版本号"""
    # 创建输出目录, 如果不存在则创建
    os.makedirs(DEFAULT_OUTPUT_DIR, exist_ok=True)

    name = base_name
    if git_version is not None:
        name = f"{name}_{git_version}"

    if system == "windows":
        return os.path.join(DEFAULT_OUTPUT_DIR, f"{name}.exe")
    return os.path.join(DEFAULT_OUTPUT_DIR, name)


def generate_zip_file_name(output_base_name, system, architecture, git_version=None):
    """根据输出文件名、操作系统和架构生成默认的 zip 文件名, 可选的git版本号"""
    os.makedirs(DEFAULT_OUTPUT_DIR, exist_ok=True)
    name = f"{output_base_name}_{system}_{architecture}"
    if git_version is not None:
        name = f"{name}_{git_version}"
    return os.path.join(DEFAULT_OUTPUT_DIR, f"{name}.zip")


def pre_build_checks(go_compiler, entry_file, use_vendor):
    """构建前的检查工作"""
    if not check_go_installed(go_compiler):
        return False
    if not check_go_mod_file():
        return False
    if not check_entry_file(entry_file):
        return False
    if use_vendor:
        try:
            run_go_mod_vendor(go_compiler)
        except Exception as e:
            print_error(str(e))
            return False
    try:
        run_go_mod_tidy(go_compiler, use_vendor)
        run_code_check(go_compiler)
        run_gofmt(go_compiler)
    except Exception as e:
        print_error(str(e))
        return False
    return True


def parse_arguments():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(description="构建 Go 应用程序")
    parser.add_argument(
        "-env",
        "--env",
        action="append",
        help="添加自定义环境变量, 格式为KEY=VALUE, 可多次使用, 例如: --env KEY=VALUE",
    )
    parser.add_argument(
        "-o", "--output", help="指定输出文件名(无需指定后缀)", default=BASE_OUTPUT_NAME
    )
    parser.add_argument(
        "-e", "--entry", help="指定入口文件路径", default=DEFAULT_ENTRY_FILE
    )
    parser.add_argument(
        "-l", "--ldflags", help="指定构建时的链接器标志", default=DEFAULT_LDFLAGS
    )
    parser.add_argument(
        "-g",
        "--go-compiler",
        help="指定 Go 编译器的绝对路径",
        default=DEFAULT_GO_COMPILER,
    )
    parser.add_argument(
        "-v",
        "--use-vendor",
        action="store_true",
        help="是否通过 vendor 克隆依赖到项目目录下",
        default=DEFAULT_USE_VENDOR,
    )
    parser.add_argument(
        "-b",
        "--use-vendor-in-build",
        action="store_true",
        help="是否在构建阶段使用 vendor 目录",
        default=DEFAULT_USE_VENDOR_IN_BUILD,
    )
    parser.add_argument(
        "-z",
        "--zip",
        action="store_true",
        help="是否将构建成功的可执行文件打包到 ZIP 文件中",
        default=DEFAULT_PACKAGE_ZIP,
    )
    parser.add_argument("--zip-file", help="指定打包输出的 ZIP 文件名", default=None)
    parser.add_argument(
        "-git",
        action="store_true",
        help="在构建时注入 Git 仓库的版本信息",
        default=DEFAULT_INJECT_GIT_INFO,
    )
    parser.add_argument(
        "-c",
        "--current-platform-only",
        action="store_true",
        help="仅构建当前运行平台的可执行文件",
        default=DEFAULT_CURRENT_PLATFORM_ONLY,
    )
    parser.add_argument(
        "-s",
        "--simple-name",
        action="store_true",
        help="使用简单文件名格式(不包含系统架构信息)",
        default=DEFAULT_SIMPLE_NAME,
    )
    parser.add_argument(
        "-p",
        "--platform",
        help="指定目标平台(如linux/windows/darwin)",
        default=None,
    )
    parser.add_argument(
        "-a",
        "--arch",
        help="指定目标架构(如amd64/arm64)",
        default=None,
    )
    parser.add_argument(
        "-batch",
        action="store_true",
        help="启用批量构建模式, 构建所有支持的平台和架构组合",
        default=False,
    )
    parser.add_argument(
        "-w",
        "--max-workers",
        type=int,
        help="批量构建模式下的最大并发线程数",
        default=DEFAULT_CONCURRENCY,
    )
    parser.add_argument(
        "-t",
        "--timeout",
        type=int,
        help="批量构建模式下每个任务的超时时间(秒), 默认30分钟(1800秒)",
        default=DEFAULT_TIMEOUT,
    )
    parser.add_argument(
        "-i",
        "--install",
        help="安装指定路径的可执行文件到GOPATH/bin目录",
        default=None,
    )
    parser.add_argument(
        "-ai",
        "--auto-install",
        action="store_true",
        help="构建完成后自动安装可执行文件到GOPATH/bin目录",
        default=False,
    )
    parser.add_argument(
        "-f",
        "--force",
        action="store_true",
        help="在安装模式下强制覆盖已存在的文件",
        default=False,
    )

    args = parser.parse_args()  # 解析命令行参数

    # 处理平台简写
    if args.platform and args.platform in PLATFORM_SHORTCUTS:
        args.platform = PLATFORM_SHORTCUTS[args.platform]

    # 处理架构简写
    if args.arch and args.arch in ARCHITECTURE_SHORTCUTS:
        args.arch = ARCHITECTURE_SHORTCUTS[args.arch]

    # 检查不支持的架构组合
    if args.platform == "darwin" and args.arch in ("386", "arm"):
        print_error(f"不支持的架构组合: darwin/{args.arch}, macOS不支持32位架构")
        return None

    return args


# 主程序入口 #
def install_executable(executable_path, args=None):
    """将可执行文件安装到GOPATH/bin目录

    参数:
        executable_path: 要安装的可执行文件路径
        args: 命令行参数对象, 包含force等标志
    """
    # 检查GOPATH环境变量
    gopath = os.getenv("GOPATH")
    if not gopath:
        print_error("未找到GOPATH环境变量, 请先设置GOPATH")
        return False

    # 检查可执行文件是否存在
    if not os.path.exists(executable_path):
        print_error(f"可执行文件 {executable_path} 不存在")
        return False

    # 创建GOPATH/bin目录(如果不存在)
    bin_path = os.path.join(gopath, "bin")
    if not os.path.exists(bin_path):
        try:
            os.makedirs(bin_path)
            print_success(f"已创建目录 {bin_path}")
        except OSError as e:
            print_error(f"创建目录 {bin_path} 失败: {str(e)}")
            return False

    # 获取目标路径
    target_path = os.path.join(bin_path, os.path.basename(executable_path))

    # 检查目标文件是否已存在
    if os.path.exists(target_path):
        if args and args.force:
            try:
                os.remove(target_path)
                print_success(f"已删除已存在的文件 {target_path}")
            except OSError as e:
                print_error(f"删除文件 {target_path} 失败: {str(e)}")
                return False
        else:
            print_error(
                f"文件 {target_path} 已存在, 请先删除或重命名, 或者使用-f/--force参数强制覆盖"
            )
            return False

    # 移动文件
    try:
        os.rename(executable_path, target_path)
        print_success(f"已安装到 {target_path}")
        return True
    except OSError as e:
        print_error(f"安装失败: {str(e)}")
        return False


def main():
    # 记录开始时间
    start_time = time.time()

    # 解析命令行参数
    args = parse_arguments()

    # 如果指定了单独安装参数
    if args.install:
        if not install_executable(args.install, args):
            sys.exit(1)
        sys.exit(0)

    # 检查批量构建模式下是否启用了简单文件名格式
    if args.batch and args.simple_name:
        print_error("批量构建模式下不能使用简单文件名格式, 请移除-s/--simple-name参数")
        return None

    # 如果启用了git标志, 提前获取git信息
    if args.git:
        print_success("正在获取 Git 信息...")
        if not get_git_info():
            sys.exit(1)

    # 如果是批量构建模式
    if args.batch:
        try:
            batch_build(args)
        except Exception as e:
            print_error(f"批量构建失败: {str(e)}")
            sys.exit(1)
        sys.exit(0)

    entry_file = args.entry  # 指定入口文件路径
    ldflags = args.ldflags  # 指定构建时的链接器标志
    go_compiler = args.go_compiler  # 指定 Go 编译器的绝对路径
    use_vendor = args.use_vendor  # 是否通过 vendor 克隆依赖到项目目录下
    zip_flag = args.zip  # 是否将构建成功的可执行文件打包到 ZIP 文件中
    inject_git = args.git  # 是否注入 Git 信息
    output_base_name = args.output  # 指定输出文件名

    # 获取操作系统和架构信息
    system = args.platform if args.platform else platform.system().lower()
    architecture = args.arch if args.arch else platform.machine().lower()

    # 检查是否仅构建当前平台
    if args.current_platform_only and system != platform.system().lower():
        print_error(
            f"当前平台为 {platform.system().lower()}, 不允许构建 {system} 平台的可执行文件"
        )
        sys.exit(1)

    # 自动转换x86_64为amd64
    if architecture == "x86_64":
        architecture = "amd64"

    # 校验平台和架构是否支持
    if system not in SUPPORTED_PLATFORMS:
        print_error(f"不支持的平台: {system}, 支持的平台: {SUPPORTED_PLATFORMS}")
        print_error(
            "支持的平台简写: "
            + ", ".join([f"{k}({v})" for k, v in PLATFORM_SHORTCUTS.items()])
        )
        sys.exit(1)
    if architecture not in SUPPORTED_ARCHITECTURES:
        print_error(
            f"不支持的架构: {architecture}, 支持的架构: {SUPPORTED_ARCHITECTURES}"
        )
        print_error(
            "支持的架构简写: "
            + ", ".join([f"{k}({v})" for k, v in ARCHITECTURE_SHORTCUTS.items()])
        )
        sys.exit(1)

    # 验证文件路径
    if not os.path.exists(entry_file):
        print_error(f"入口文件 {entry_file} 不存在")
        sys.exit(1)

    # 执行构建前的检查工作
    if not pre_build_checks(go_compiler, entry_file, use_vendor):
        sys.exit(1)

    # 获取构建时间
    build_time = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

    # 根据参数注入 Git 信息
    if inject_git:
        # 直接从缓存获取git信息
        git_version = _git_info_cache["version"]
        git_commit = _git_info_cache["commit"]
        format_time = _git_info_cache["commit_time"]
        git_status = _git_info_cache["status"]
        if git_version is None:
            sys.exit(1)
        # 注入 Git 信息到链接器标志
        ldflags = LD_FLAGS_TEMPLATE.format(
            app_name=output_base_name,
            git_version=git_version,
            git_commit=git_commit,
            commit_time=format_time,
            build_time=build_time,
            tree_state=git_status,
        )
        print_success(f"Git信息已注入: {git_version} ({git_commit})")

    # 生成默认的 zip 文件名
    if args.zip_file is None:
        zip_file = generate_zip_file_name(
            output_base_name, system, architecture, git_version if inject_git else None
        )
    else:
        zip_file = args.zip_file

    # 如果指定了简单文件名格式, 生成则生成不带系统架构信息的输出文件名
    if args.simple_name:
        base_name = f"{BASE_OUTPUT_NAME}"
        # 生成默认的输出文件名
        output_file = generate_output_file_name(base_name, system, None)
    else:
        # 生成带有系统和架构信息的默认输出文件名
        base_name = f"{BASE_OUTPUT_NAME}_{system}_{architecture}"
        # 生成默认的输出文件名
        output_file = generate_output_file_name(
            base_name, system, git_version if inject_git else None
        )

    # 执行构建命令
    print_success("开始构建...")
    build_config = BuildConfig(
        go_compiler=go_compiler,
        output_file=output_file,
        entry_file=entry_file,
        ldflags=ldflags,
        use_vendor_in_build=args.use_vendor_in_build,
        is_batch=args.batch,
        args=args,
    )
    build_result = build_go_app(build_config)

    # 判断构建结果
    if build_result:
        if not args.batch:
            print_success("构建完成。")
        if zip_flag:
            zip_executable(output_file, zip_file, args.batch)
    else:
        print_error("构建失败, 请检查错误信息。")

    # 记录结束时间
    end_time = time.time()
    # 计算构建耗时
    elapsed_time = end_time - start_time
    print_success(f"本次构建耗时: {elapsed_time:.2f} 秒")

    # 单独构建模式下自动安装
    if args.auto_install and not args.batch:
        if not install_executable(output_file, args):
            sys.exit(1)
        sys.exit(0)


if __name__ == "__main__":
    main()
