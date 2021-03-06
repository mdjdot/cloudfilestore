package handler

import (
	"cloudfilestore/meta"
	"cloudfilestore/util"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// HomeHandler 处理主页
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/view/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("StatusInternalServerError"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// 接受文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			log.Println("failed to receive the file")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./temp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			log.Println("failed to create file")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("StatusServiceUnavailable"))
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			log.Println("failed to save file")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("StatusServiceUnavailable"))
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetaDB(fileMeta)

		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("upload finished"))
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler 处理上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}

// GetFileMetaHandler 处理查询文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		fileHash := r.Form["filehash"][0]
		// fileHash := r.PostForm["filehash"][0]
		// fileHash := r.PostFormValue("filehash")
		// fileMeta := meta.GetFileMeta(fileHash)
		fileMeta := meta.GetFileMetaDB(fileHash)
		if fileMeta.FileSize <= 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("文件未找到"))
			return
		}
		data, err := json.Marshal(fileMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// DownloadHandler 处理下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		fileSha1 := r.Form.Get("filehash")
		// fileMeta := meta.GetFileMeta(fileSha1)
		fileMeta := meta.GetFileMetaDB(fileSha1)

		file, err := os.Open(fileMeta.Location)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		w.Header().Set("Content-Type", "application/octect-stream")
		w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// FileMetaUpdateHandler 更新文件元信息（重命名）
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		opType := r.Form.Get("op")
		fileSha1 := r.Form.Get("filehash")
		newFileNam := r.Form.Get("filename")

		if opType != "0" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("StatusForbidden"))
			return
		}
		// fileMeta := meta.GetFileMeta(fileSha1)
		fileMeta := meta.GetFileMetaDB(fileSha1)
		fileMeta.FileName = newFileNam
		meta.UpdateFileMeta(fileMeta)

		data, err := json.Marshal(fileMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// FileDeleteHandler 处理文件删除
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		fileSha1 := r.FormValue("filehash")
		if len(fileSha1) > 0 {
			// fileMeta := meta.GetFileMeta(fileSha1)
			fileMeta := meta.GetFileMetaDB(fileSha1)
			if fileMeta.FileSize > 0 {
				meta.RemoveFileMeta(fileSha1)
				os.Remove(fileMeta.Location)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("删除成功"))
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("StatusForbidden"))
	}
}
