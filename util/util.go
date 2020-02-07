package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// Sha1Stream 结构体，有一个 hash.Hash 字段 _sha1
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update 方法
func (ss *Sha1Stream) Update(data []byte) {
	if ss._sha1 == nil {
		ss._sha1 = sha1.New()
	}
	ss._sha1.Write(data)
}

// Sum 方法
func (ss *Sha1Stream) Sum() string {
	return hex.EncodeToString(ss._sha1.Sum([]byte("")))
}

// Sha1 函数
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

// FileSha1 函数
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

// MD5 函数
func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

// FileMD5 函数
func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

// PathExists 函数
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetFileSize 函数
func GetFileSize(fileName string) int64 {
	var result int64
	filepath.Walk(fileName, func(path string, fileInfo os.FileInfo, err error) error {
		result = fileInfo.Size()
		return nil
	})
	return result
}
