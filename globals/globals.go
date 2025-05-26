package globals

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	_ "embed"
	"hash"
)

var (
	// 支持的哈希算法列表
	SupportedAlgorithms = map[string]func() hash.Hash{
		"md5":    md5.New,
		"sha1":   sha1.New,
		"sha256": sha256.New,
		"sha512": sha512.New,
	}

	// 输出哈希值的文件名
	OutputFileName = "checksum.hash"

	// 输出对比结果的文件名
	OutputCheckFileName = "check_dir.check"

	// 时间戳格式
	TimestampFormat = "2006-01-02 15:04:05"
)

// hash 子命令帮助信息
//
//go:embed help/help_hash.txt
var HashHelp string

// size 子命令帮助信息
//
//go:embed help/help_size.txt
var SizeHelp string

// check 子命令帮助信息
//
//go:embed help/help_check.txt
var CheckHelp string

// find 子命令帮助信息
//
//go:embed help/help_find.txt
var FindHelp string

// fck 主命令帮助信息
//
//go:embed help/help.txt
var FckHelp string
