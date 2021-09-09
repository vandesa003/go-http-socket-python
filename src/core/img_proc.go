package main

import (
	"errors"
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// func main() {
// 	img := gocv.IMRead("data/test.png", gocv.IMReadColor)
// 	fmt.Println(img.Size())
// 	blur, err := preprocess(img)
// 	if err != nil {
// 		fmt.Println("error in preprocess")
// 		fmt.Println(err)
// 	}
// 	gocv.IMWrite("data/res.png", blur)
// 	flag, r, err := keyPointsColor(blur, 50)
// 	fmt.Println(flag, r)
// 	a, b, err := mainColor(blur)
// 	fmt.Println(a, b)

// }

func maxInSlice(arr []float64) float64 {
	var max float64 = -1.0
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func minInSlice(arr []float64) float64 {
	var min float64 = 256.0
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}

func sumInSlice(arr []float64) float64 {
	var sum float64 = 0.0
	for _, v := range arr {
		sum += v
	}
	return sum
}

// from https://github.com/hybridgroup/gocv/issues/833
func GetVecbAt(m gocv.Mat, row int, col int) []uint8 {
	ch := m.Channels()
	v := make([]uint8, ch)

	for c := 0; c < ch; c++ {
		v[c] = m.GetUCharAt(row, col*ch+c)
	}

	return v
}

func preprocess(img gocv.Mat) (gocv.Mat, error) {
	gray := gocv.NewMat()
	blur := gocv.NewMat()
	shape := img.Size()
	fmt.Println(shape)
	switch {
	case len(shape) == 3 && img.Size()[2] == 3:
		gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	case len(shape) == 3 && img.Size()[2] == 4:
		gocv.CvtColor(img, &gray, gocv.ColorBGRAToGray)
	case len(img.Size()) == 2:
		gray = img
	default:
		return blur, errors.New("wrong input.")
	}
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &blur, image.Point{3, 3}, 0.0, 0.0, gocv.BorderDefault)
	return blur, nil
}

func keyPointsColor(grayImg gocv.Mat, drift int8) (bool, float64, error) {
	h, w := grayImg.Size()[0], grayImg.Size()[1]
	var hBound int
	var wBound int
	if w <= 10 || h <= 10 {
		return false, 0, errors.New("image size ({}, {}) too small!")
	}
	if h < 100 {
		hBound = 3
	} else {
		hBound = 5
	}
	if w < 100 {
		wBound = 3
	} else {
		wBound = 5
	}
	topLeft := grayImg.Region(image.Rect(0, 0, wBound, hBound)) // left, top, right, bottom
	topRight := grayImg.Region(image.Rect(w-wBound, 0, w, hBound))
	bottomLeft := grayImg.Region(image.Rect(0, h-hBound, wBound, h))
	bottomRight := grayImg.Region(image.Rect(w-wBound, h-hBound, w, h))
	var keyPoints = [4]float64{
		topLeft.Mean().Val1,
		topRight.Mean().Val1,
		bottomLeft.Mean().Val1,
		bottomRight.Mean().Val1,
	}
	if maxInSlice(keyPoints[:])-minInSlice(keyPoints[:]) < float64(drift) {
		r := sumInSlice(keyPoints[:]) / 4
		return true, r, nil
	} else {
		return true, 0, nil
	}
}

// calculate main color at corners.
func mainColor(grayImg gocv.Mat) (uint8, float32, error) {
	const bins int = 256
	var colorHist [bins]float32
	var pixelNum float32 = 0
	var maxColor uint8 = 0
	var maxRatio float32 = 0

	mask := gocv.NewMat()
	ch := []int{0}
	src := []gocv.Mat{grayImg}
	hist := gocv.NewMat()
	size := []int{bins}
	ranges := []float64{0, 256}
	gocv.CalcHist(src, ch, mask, &hist, size, ranges, false)

	for i := 0; i < hist.Size()[1]; i++ { // row
		for j := 0; j < hist.Size()[0]; j++ { // col
			v := hist.GetFloatAt(i, j)
			pixelNum += v
			colorHist[j] = v
		}
	}
	for i := 0; i < bins; i++ {
		r := float32(colorHist[i]) / float32(pixelNum)
		if r > maxRatio {
			maxColor = uint8(i)
			maxRatio = r
		}
	}
	fmt.Println(pixelNum)

	return maxColor, maxRatio, nil
}

func generateMask(img gocv.Mat, th float32, drift int8) (gocv.Mat, error) {
	mask := gocv.NewMat()
	gray, err := preprocess(img)
	if err != nil {
		return mask, err
	}
	isSame, keyColor, err := keyPointsColor(gray, drift)
	maxColor, maxRatio, err := mainColor(gray)
	if isSame == false {
		fmt.Println("Warning: key points color not same!")
		return mask, errors.New("key points color not same!")
	}
	if uint8(keyColor) != maxColor {
		fmt.Println("Warning: key color != max color!")
	}
	if maxRatio > th {
		// TODO: finish this(marquee_todo)

	} else {

	}
	return mask, nil
}

func blend(img gocv.Mat, mask gocv.Mat) (gocv.Mat, error) {
	rgba := gocv.NewMat()
	img.CopyToWithMask(&rgba, mask)
	return rgba, nil
}
