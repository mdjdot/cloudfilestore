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

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
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
		meta.UpdateFileMeta(fileMeta)

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
	if r.Method == http.MethodGet {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/meta.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()

		// fileHash := r.Form["filehash"][0]
		fileHash := r.PostForm["filehash"][0]
		// fileHash := r.PostFormValue("filehash")
		fileMeta := meta.GetFileMeta(fileHash)
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
