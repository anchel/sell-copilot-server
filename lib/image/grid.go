package image

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"sort"
	"sync"

	"github.com/disintegration/imaging"
)

type List struct {
	images []image.Image
}

type chanImgObj struct {
	img   image.Image
	index int
}

type chanImgObjSlice []*chanImgObj

func (ci chanImgObjSlice) Len() int {
	return len(ci)
}
func (ci chanImgObjSlice) Less(i, j int) bool {
	return ci[i].index < ci[j].index
}
func (ci chanImgObjSlice) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}
func (ci chanImgObjSlice) Sort() {
	sort.Sort(ci)
}

func NewList(list []string) (*List, error) {
	var wg sync.WaitGroup

	imageChan := make(chan *chanImgObj, len(list))

	for index, filePath := range list {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()
			src, err := imaging.Open(filePath)
			if err != nil {
				log.Printf("failed to open image: %v", err)
				// imageChan <- nil
			} else {
				imageChan <- &chanImgObj{src, index}
			}
		}(index, filePath)
	}

	wg.Wait()

	// fmt.Println("len(imageChan)", len(imageChan))

	imageObjList := make(chanImgObjSlice, 0, len(imageChan))
	l := len(imageChan)
	for i := 0; i < l; i++ {
		// fmt.Println("dddd", i)
		imgObj, ok := <-imageChan
		// fmt.Println("ddd", ok, img.Bounds())
		if !ok {
			// fmt.Println("dddd not ok", i, imgObj.index)
			break
		}
		imageObjList = append(imageObjList, imgObj)
	}

	imageObjList.Sort() // 排序

	imageList := make([]image.Image, 0, imageObjList.Len())
	for _, imgObj := range imageObjList {
		imageList = append(imageList, imgObj.img)
	}

	imageSet := &List{imageList}
	return imageSet, nil
}

func (is *List) ApplyGridLayout(dx, gutter int) (image.Image, error) {
	images := is.images
	l := len(images)

	// fmt.Println("l=", l)

	rows := l / dx

	y := l % dx
	if y != 0 {
		rows += 1
	}

	rowGroupList := make([]*RowGroup, rows)
	maxWidth := 0

	sumHeight := 0

	for i := 0; i < rows; i++ {
		offset := i * dx
		end := offset + dx
		if end > l {
			end = l
		}
		// fmt.Println("offset, end", offset, end)
		rg := newRowGroup(images[offset:end])
		rowGroupList[i] = rg

		if rg.Width > maxWidth {
			maxWidth = rg.Width
		}

		sumHeight += rg.Height
	}

	sumWidth := dx*maxWidth + (dx+1)*gutter
	sumHeight = sumHeight + (rows+1)*gutter

	if rows <= 1 { // 只有一行时，宽度为所有图片的宽度之和
		sumWidth = l*maxWidth + (l+1)*gutter
	}

	canvasMaxPoint := image.Point{X: sumWidth, Y: sumHeight}
	canvasRect := image.Rectangle{Min: image.Point{}, Max: canvasMaxPoint}
	canvas := image.NewRGBA(canvasRect)

	var lastRg *RowGroup
	var leiHeight int
	for i, rg := range rowGroupList {
		if lastRg != nil {
			leiHeight += lastRg.Height
		}
		//fmt.Println("preHeight", i, preHeight, rg.Height)
		for y, img := range rg.Images {
			minPoint := image.Point{X: y*maxWidth + (y+1)*gutter, Y: leiHeight + (i+1)*gutter}
			maxPoint := minPoint.Add(image.Point{X: maxWidth, Y: rg.Height})
			rect := image.Rectangle{Min: minPoint, Max: maxPoint}

			draw.Draw(canvas, rect, img, image.Point{}, draw.Src)
		}
		fmt.Println()
		lastRg = rg
	}

	return canvas, nil
}

func Save(img image.Image, dstPath string) error {
	err := imaging.Save(img, dstPath)
	return err
}

func newRowGroup(images []image.Image) *RowGroup {
	rg := &RowGroup{
		Images: images,
	}
	rg.initMax()
	return rg
}

type RowGroup struct {
	Images   []image.Image
	MinWidth int
	MaxWidth int

	MinHeight int
	MaxHeight int

	Width  int
	Height int

	WidthSum int
}

func (rg *RowGroup) initMax() {
	X := rg.Images[0].Bounds().Max.X
	minWidth := X
	maxWidth := X

	for i := 1; i < len(rg.Images); i++ {
		if rg.Images[i].Bounds().Max.X > maxWidth {
			maxWidth = rg.Images[i].Bounds().Max.X
		}
		if rg.Images[i].Bounds().Max.X < minWidth {
			minWidth = rg.Images[i].Bounds().Max.X
		}
	}
	rg.MinWidth = minWidth
	rg.MaxWidth = maxWidth
	rg.Width = (rg.MinWidth + rg.MaxWidth) / 2
	rg.WidthSum = rg.Width * len(rg.Images)
	// fmt.Println("ccc", rg.MinWidth, rg.MaxWidth, rg.Width)

	minHeight := 0
	maxHeight := 0

	// 根据平均宽度来进行缩放
	for i := 0; i < len(rg.Images); i++ {
		img := imaging.Resize(rg.Images[i], rg.Width, 0, imaging.Lanczos)
		rg.Images[i] = img
		if img.Bounds().Max.Y > maxHeight {
			maxHeight = img.Bounds().Max.Y
		}
		if img.Bounds().Max.Y < minHeight {
			minHeight = img.Bounds().Max.Y
		}
	}

	rg.MinHeight = minHeight
	rg.MaxHeight = maxHeight
	rg.Height = rg.MaxHeight
}
