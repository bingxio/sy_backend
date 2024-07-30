package api

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"sy_backend/config"
)

const (
	menu = "/menu"
)

func NewResource(folder string, file *multipart.FileHeader) string {
	src, _ := file.Open()
	defer src.Close()

	path := fmt.Sprintf("%s%s", config.Conf.ResourcePath, folder)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		log.Println(err)
	}
	path = fmt.Sprintf("%s/%s", path, file.Filename)
	dst, _ := os.Create(path)
	defer dst.Close()

	_, _ = io.Copy(dst, src)
	return path
}

func DeleteResource(path string) {
	_ = os.Remove(path)
}
