package tests

import (
	"fmt"
	"os"
	"sy_backend/util"
	"testing"
)

func TestImage(t *testing.T) {
	path := "../resource/menu/304f1b6d5f9b063eb0301fe147c5422.jpg"

	info, _ := os.Stat(path)
	fmt.Printf("info.Size(): %f\n", float64(info.Size())/float64(1024*1024))

	err := util.CompressImage(path)
	if err != nil {
		t.Fatal(err)
	}

	info, _ = os.Stat(path)
	fmt.Printf("info.Size(): %f\n", float64(info.Size())/float64(1024*1024))
}
