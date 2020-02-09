package meta

import "cloudfilestore/db"

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// UpdateFileMetaDB 更新文件元信息到数据库
func UpdateFileMetaDB(fileMeta FileMeta) bool {
	return db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
}

// GetFileMeta 获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB 从数据库获取元文件信息
func GetFileMetaDB(fileSha1 string) FileMeta {
	fileMeta, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}
	}
	return FileMeta{
		FileSha1: fileMeta.FileHash,
		FileName: fileMeta.FileName,
		FileSize: fileMeta.FileSize,
		Location: fileMeta.FileAddr,
	}
}

// RemoveFileMeta 删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
