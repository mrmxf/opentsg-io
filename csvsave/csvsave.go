// Package csvsave is used to save independent r,g and b csv files of an image
package csvsave

import (
	"encoding/csv"
	"fmt"
	"image"
	"os"
	"strconv"
)

// Encode generates an .r.csv, .b.csv and a .g.csv file with each csv value representing a pixel
func Encode(filename string, content *image.NRGBA64) error {

	fcutOff := len(filename) - 4
	files := [3]string{filename[:fcutOff] + ".r.csv", filename[:fcutOff] + ".g.csv", filename[:fcutOff] + ".b.csv"}
	colourStrings := imageToString(content)

	for i := range files {
		f, fErr := os.OpenFile(files[i], os.O_RDWR|os.O_CREATE, 0777)
		defer f.Close()
		if fErr != nil {
			//do not run if there's no file to save to
			return fErr
		} else {
			//assign a way to write to the file
			w := csv.NewWriter(f)
			w.WriteAll(colourStrings[i])
		}
	}
	//array of rgb
	fmt.Printf("%s generated \n", files)
	return nil
}

func imageToString(canvas *image.NRGBA64) [][][]string {
	//make an array of rgb
	imgString := make([][][]string, 3)
	for i := range imgString {
		imgString[i] = make([][]string, canvas.Bounds().Max.Y)
		for j := range imgString[i] {
			imgString[i][j] = make([]string, canvas.Bounds().Max.X)
		}
	}
	//loop through each position and assign the rgb section of the array with the values
	for i := 0; i < canvas.Bounds().Max.Y; i++ {
		for j := 0; j < canvas.Bounds().Max.X; j++ {
			cVal := canvas.At(j, i)
			r, g, b, _ := cVal.RGBA()
			colours := []uint32{r, g, b}
			for k := range colours {
				stringColour := strconv.FormatUint(uint64(colours[k]), 10)
				imgString[k][i][j] = stringColour
			}
		}
	}
	return imgString
}
