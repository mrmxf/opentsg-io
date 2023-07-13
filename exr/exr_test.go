package exr

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"
	"gitlab.com/mmTristan/tpg-core/aces"
	. "github.com/smartystreets/goconvey/convey"
)

// identify -verbose test.exr
func TestNRGBA64Write(t *testing.T) {

	vals := []uint16{0xffff, 0x8000, 0}
	target := []string{"./testdata/full.exr", "./testdata/half.exr", "./testdata/black.exr"}

	for i, v := range vals {

		//generate a box of certain colours
		box := image.NewNRGBA64(image.Rect(0, 0, 100, 100))
		colors := make(map[int]color.NRGBA64)
		colors[0] = color.NRGBA64{R: v, A: 0xffff}
		colors[1] = color.NRGBA64{G: v, A: 0xffff}
		colors[2] = color.NRGBA64{B: v, A: 0xffff}
		colors[3] = color.NRGBA64{R: v, G: v, B: v, A: 0xffff}

		for y := 0; y < box.Bounds().Max.Y; y++ {
			c := colors[y%4]
			for x := 0; x < box.Bounds().Max.X; x++ {
				box.SetNRGBA64(x, y, c)

			}
		}
		var mock bytes.Buffer
		Encode(&mock, box)

		file, _ := os.ReadFile(target[i])

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(file)
		htest.Write(mock.Bytes())

		Convey("Checking the exr files are saved with nrgba64", t, func() {
			Convey(fmt.Sprintf("Compared the generated file to %v", target[i]), func() {
				Convey("No error is returned and the file matches", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}

}

// TestACESWRITE checks that files written in aces will match the nrgba64 for integers
func TestACESWrite(t *testing.T) {

	vals := []uint16{0xffff, 0x8000, 0}
	target := []string{"./testdata/full.exr", "./testdata/half.exr", "./testdata/black.exr"}

	for i, v := range vals {

		//generate a box of certain colours
		box := aces.NewARGBA(image.Rect(0, 0, 100, 100))
		colors := make(map[int]color.NRGBA64)
		colors[0] = color.NRGBA64{R: v, A: 0xffff}
		colors[1] = color.NRGBA64{G: v, A: 0xffff}
		colors[2] = color.NRGBA64{B: v, A: 0xffff}
		colors[3] = color.NRGBA64{R: v, G: v, B: v, A: 0xffff}

		for y := 0; y < box.Bounds().Max.Y; y++ {
			c := colors[y%4]
			for x := 0; x < box.Bounds().Max.X; x++ {
				box.Set(x, y, c)

			}
		}
		var mock bytes.Buffer
		Encode(&mock, box)

		file, _ := os.ReadFile(target[i])

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(file)
		htest.Write(mock.Bytes())

		//png treats it as a rgb48 which aces is technically representing here
		// so do tiff files
		//f, _ := os.Create(fmt.Sprintf("p%v.exr", i))
		//fmt.Println(Encode(f, box))
		Convey("Checking the exr files are saved with aces", t, func() {
			Convey(fmt.Sprintf("Compared the generated file to %v", target[i]), func() {
				Convey("No error is returned and the file matches", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}

}
