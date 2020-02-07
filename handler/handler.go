package handler

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
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
		newFile, err := os.Create("./temp/" + head.Filename)
		if err != nil {
			log.Println("failed to create file")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("StatusServiceUnavailable"))
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Println("failed to save file")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("StatusServiceUnavailable"))
			return
		}
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("upload finished"))
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler 处理上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}
