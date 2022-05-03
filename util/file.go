package util

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	FileSha1CMD = `
	#!/bin/bash
	sha1sum $1 | awk '{print $1}'
	`

	// FileSizeCMD : 计算文件大小
	FileSizeCMD = `
	#!/bin/bash
	ls -l $1 | awk '{print $5}'
	`

	// FileChunksDelCMD : 删除文件分块
	FileChunksDelCMD = `
	#!/bin/bash
	chunkDir="/data/chunks/"
	targetDir=$1
	# 增加条件判断，避免误删  (指定的路径包含且不等于chunkDir)
	if [[ $targetDir =~ $chunkDir ]] && [[ $targetDir != $chunkDir ]]; then 
	  rm -rf $targetDir
	fi
	`
)

//ComputeSha1ByShell 通过调用shell计算sha1
func ComputeSha1ByShell(destPath string) (string, error) {

	// Replace 替换字符 包含 $1 替换成新的字符串
	cmdStr := strings.Replace(FileSha1CMD, "$1", destPath, 1)
	hashCmd := exec.Command("bash", "-c", cmdStr)
	if fileHash, err := hashCmd.Output(); err != nil {
		return "", err
	} else {
		reg := regexp.MustCompile("\\s+")
		s := reg.ReplaceAllString(string(fileHash), "")
		return s, nil
	}
}

// RemovePathByShell : 通过调用shell来删除制定目录
// @return bool: 合并成功将返回true, 否则返回false
func RemovePathByShell(destPath string) bool {
	cmdStr := strings.Replace(FileChunksDelCMD, "$1", destPath, 1)
	delCmd := exec.Command("bash", "-c", cmdStr)
	if _, err := delCmd.Output(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
