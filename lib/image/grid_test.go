package image

import (
	"fmt"
	"os"
	"testing"
)

func TestGrid(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println("TestGrid", wd)
	is, err := NewList([]string{
		"../../dist/image/upload-1721446571650028.png",
		"../../dist/image/upload-1721470718025484.png",
		"../../dist/image/upload-1721492772406986.jpg",
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	img, err := is.ApplyGridLayout(3, 100)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	err = Save(img, "../../dist/image/merge.jpg")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
}
