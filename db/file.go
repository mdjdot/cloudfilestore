package db

import (
	filedb "cloudfilestore/db/mysql"
	"log"
)

// OnFileUploadFinished 文件上传完成，保存元信息到数据库
func OnFileUploadFinished(fileHash, fileName string, fileSize int64, fileAddr string) bool {
	stmt, err := filedb.FileDB.Prepare("insert ignore into tbl_file (file_sha1,file_name,file_size,file_addr,status) values(?,?,?,?,1)")
	if err != nil {
		log.Println("failed to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		log.Println("failed to insert data, err: ", err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); rf <= 0 && err == nil {
		log.Println("文件已存在，文件hash：", fileHash)
		return true
	}
	return true
}

// TableFile 文件结构体
type TableFile struct {
	FileHash string
	FileName string
	// FileSize sql.NullInt64
	FileSize int64
	FileAddr string
}

// GetFileMeta 从数据库获取文件元信息
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := filedb.FileDB.Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tableFile := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&tableFile.FileHash, &tableFile.FileAddr, &tableFile.FileName, &tableFile.FileSize)
	if err != nil {
		return nil, err
	}
	return &tableFile, nil
}
