package util

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func CompressImage(imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	file.Seek(0, 0)

	var img image.Image
	var format string

	switch t := http.DetectContentType(buffer); t {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
		format = "jpeg"
	case "image/png":
		img, err = png.Decode(file)
		format = "png"
	default:
		return errors.New("expected image ContentType")
	}
	if err != nil {
		return err
	}
	file.Close()
	file, err = os.OpenFile(imagePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var opt jpeg.Options
	opt.Quality = 75

	switch format {
	case "jpeg":
		err = jpeg.Encode(file, img, &opt)
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(file, img)
	}
	return err
}
