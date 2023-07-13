package tiffup

import (
	"bytes"
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
	//file, _ := os.Open("./testimages/standard.tiff")
	//decode to get the colour values
	//baseVals, _ := tiff.Decode(file)
	//assign the colour to the correct type of image NGRBA64 and replace the colour values
	//readImage := image.NewNRGBA64(baseVals.Bounds())

	//draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

	matchFiles := []string{"./testdata/16b.tiff"}
	for i := range matchFiles {
		file, _ := os.Open(matchFiles[i])
		//decode to get the colour values
		baseVals, err := tiff.Decode(file)
		fmt.Println(err)
		//assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		draw.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)
		hnormal := sha256.New()
		htest := sha256.New()

		//open the files to be tested
		fi, _ := file.Stat()
		fnormal := make([]byte, fi.Size())
		file.Read(fnormal) //, _ := os.ReadFile(matchFiles[i])

		out := new(bytes.Buffer)
		Encode(out, readImage)

		//write the hash with the information
		hnormal.Write(fnormal)
		htest.Write(out.Bytes())
		Convey("Checking the image is saved and matches the example file exactly", t, func() {
			Convey("using a basic 16 bit tiff as an input", func() {
				Convey("An identical file is returned to the example 16 bit file", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}
}
