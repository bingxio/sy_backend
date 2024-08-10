package api

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sy_backend/config"
	"sy_backend/util"
)

const (
	menu = "/menu"
)

func NewResource(
	folder,
	nanoId string,
	file *multipart.FileHeader,
) (string, error) {
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
	path = fmt.Sprintf(
		"%s/%s%s",
		path,
		nanoId,
		filepath.Ext(file.Filename),
	)
	dst, _ := os.Create(path)
	defer dst.Close()

	_, _ = io.Copy(dst, src)

	if err := util.CompressImage(path); err != nil {
		return "", err
	}
	return path, nil
}

func DeleteResource(path string) error {
	return os.Remove(path)
}
