// Package savefunc contains the different methods for saving files
package savefunc

import (
	"image"
	"image/png"
	"os"

	"github.com/mrmxf/opentsg-io/csvsave"
	"github.com/mrmxf/opentsg-io/dpx"
	"github.com/mrmxf/opentsg-io/exr"
	"github.com/mrmxf/opentsg-io/tiffup"
	"golang.org/x/image/tiff"
)

// writeTiffFile saves the file as a tiff
func WriteTiffFile(f *os.File, img *image.NRGBA64) error {

	//save the file depending on if it's transparent
	if img.Opaque() {
		return tiffup.Encode(f, img)
	} else {
		return tiff.Encode(f, img, nil)
	}
}

// writePngFile saves the file as a png
func WritePngFile(f *os.File, image *image.NRGBA64) error {
	return png.Encode(f, image)
}

func WriteExrFile(f *os.File, image image.Image) error {
	return exr.Encode(f, image)
}

func WriteDPXFile(f *os.File, image *image.NRGBA64, bit int) error {
	//default all files to 16 bit
	if bit == 0 {
		bit = 16
	}

	return dpx.Encode(f, image, &dpx.Options{Bitdepth: bit})
}

func WriteCSVFile(filename string, image *image.NRGBA64) error {
	return csvsave.Encode(filename, image)
}
