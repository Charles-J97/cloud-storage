package middleware

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	// MergeFileCMD : 通过shell合并分块文件
	MergeFileCMD = `
	#!/bin/bash
	# 需要进行合并的分片所在的目录
	chunkDir=$1
	# 合并后的文件的完成路径(目录＋文件名)
	mergePath=$2
	
	echo "分块合并，输入目录: " $chunkDir
	
	if [ ! -f $mergePath ]; then
			echo "$mergePath not exist"
	else
			rm -f $mergePath
	fi
	
	for chunk in $(ls $chunkDir | sort -n)
	do
			cat $chunkDir/${chunk} >> ${mergePath}
	done
	
	echo "合并完成，输出：" mergePath
	`

	// FileSha1CMD : 计算文件sha1值
	FileSha1CMD = `
	#!/bin/bash
	sha1sum $1 | awk '{print $1}'
	`
)

// ComputeSha1ByShell : 通过调用shell来计算文件sha1
// @return  (string, error): (文件hash, 错误信息)
func ComputeSha1ByShell(destPath string) (string, error) {
	cmdStr := strings.Replace(FileSha1CMD, "$1", destPath, 1)
	hashCmd := exec.Command("bash", "-c", cmdStr)
	if filehash, err := hashCmd.Output(); err != nil {
		fmt.Println(err)
		return "", err
	} else {
		reg := regexp.MustCompile("\\s+")
		return reg.ReplaceAllString(string(filehash), ""), nil
	}
}

// MergeChuncksByShell : 通过调用shell来合并文件分块，分块文件名须有序 (如分块名分别为: 1, 2, 3, ...)
// @return bool: 合并成功将返回true, 否则返回false
func MergeChuncksByShell(chunkDir string, destPath string, fileSha1 string) bool {
	// 合并分块
	cmdStr := strings.Replace(MergeFileCMD, "$1", chunkDir, 1)
	cmdStr = strings.Replace(cmdStr, "$2", destPath, 1)
	mergeCmd := exec.Command("bash", "-c", cmdStr)
	if _, err := mergeCmd.Output(); err != nil {
		fmt.Println(err)
		return false
	}

	//// 计算合并后的文件hash
	//if filehash, err := ComputeSha1ByShell(destPath); err != nil {
	//	fmt.Println(err)
	//	return false
	//} else if string(filehash) != fileSha1 { // 判断文件hash是否符合给定值
	//	fmt.Println(filehash + " " + fileSha1)
	//	return false
	//} else {
	//	fmt.Println("check sha1: " + destPath + " " + filehash + " " + fileSha1)
	//}

	return true
}
