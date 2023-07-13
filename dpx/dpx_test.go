package dpx

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/draw"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"golang.org/x/image/tiff"
)

func TestEncode(t *testing.T) {
	//"testing encoding of 8,10,12 and 16 bit files against known versions, this covers encoding with correct input values
	//this is the standard 16 bit image as a tiff
	file, _ := os.Open("./testimages/standard.tiff")
	//decode to get the colour values
	baseVals, _ := tiff.Decode(file)
	//assign the colour to the correct type of image NGRBA64 and replace the colour values
	readImage := image.NewNRGBA64(baseVals.Bounds())

	draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

	matchFiles := []string{"./testimages/8b.dpx", "./testimages/10b.dpx", "./testimages/12b.dpx", "./testimages/16b.dpx"}
	matchFileDepth := []int{8, 10, 12, 16}
	for i, mf := range matchFiles {
		hnormal := sha256.New()
		htest := sha256.New()

		//open the files to be tested
		fnormal, _ := os.ReadFile(mf)
		ftest := headerGen(readImage, matchFileDepth[i])

		//write the hash with the information
		hnormal.Write(fnormal)
		htest.Write(ftest)
		Convey("Checking the image is saved and matches the example file exactly", t, func() {
			Convey(fmt.Sprintf("using a a tag of %v bits", matchFileDepth[i]), func() {
				Convey("An identical file is returned to the example 8 bit file", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}
}

// go test ./dpx -bench=. -benchtime=40s
// 260         161384890 ns/op
func BenchmarkEncode12(b *testing.B) {
	file, _ := os.Open("./testimages/standard.tiff")
	//decode to get the colour values
	baseVals, _ := tiff.Decode(file)
	//assign the colour to the correct type of image NGRBA64 and replace the colour values
	readImage := image.NewNRGBA64(baseVals.Bounds())

	draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		headerGen(readImage, 12)
	}
}

// 222         194536411 ns/op
func BenchmarkEncode16(b *testing.B) {
	file, _ := os.Open("./testimages/standard.tiff")
	//decode to get the colour values
	baseVals, _ := tiff.Decode(file)
	//assign the colour to the correct type of image NGRBA64 and replace the colour values
	readImage := image.NewNRGBA64(baseVals.Bounds())

	draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		headerGen(readImage, 16)
	}
}

// 374         108842332 ns/op
func BenchmarkEncode8(b *testing.B) {
	file, _ := os.Open("./testimages/standard.tiff")
	//decode to get the colour values
	baseVals, _ := tiff.Decode(file)
	//assign the colour to the correct type of image NGRBA64 and replace the colour values
	readImage := image.NewNRGBA64(baseVals.Bounds())

	draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		headerGen(readImage, 8)
	}
}

//add tests here to check the byte form saving
/*
goos: linux
goarch: amd64
pkg: github.com/mmTristan/tpg-io-private/dpx
cpu: AMD EPYC 7B13
BenchmarkEncode12-16                 820          50008853 ns/op
BenchmarkEncode16-16                 968          44181842 ns/op
BenchmarkEncode8-16                 1339          39128481 ns/op

cpu: AMD EPYC 7B13
BenchmarkEncode12-16                 896          48595325 ns/op
BenchmarkEncode16-16                 992          45209761 ns/op
BenchmarkEncode8-16                 1465          30837210 ns/op


pkg: github.com/mmTristan/tpg-io-private/dpx
cpu: AMD EPYC 7B13
BenchmarkEncode12-16                1051          41338020 ns/op
BenchmarkEncode16-16                1101          42051351 ns/op
BenchmarkEncode8-16                 2132          20902725 ns/op

BenchmarkEncode12-16                1046          45226782 ns/op
BenchmarkEncode16-16                1098          47361543 ns/op
BenchmarkEncode8-16                 1395          29190676 ns/op
*/
