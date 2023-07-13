package exr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"math"

	"github.com/x448/float16"
	"gitlab.com/mmTristan/tpg-core/aces"
)

const (
	version = uint32(2)
	//used to differentiate between most things within the file
	terminator = 0x00
)

var (
	magic = []byte{0x76, 0x2F, 0x31, 0x01}
)

func Encode(w io.Writer, in image.Image) error {
	b := in.Bounds().Max
	//open exr are little endian
	dir := binary.LittleEndian
	//uint32 for 3 bytes as no extra information for a basic exr file
	head := dir.AppendUint32(magic, version)

	// Generate the information to be used in the headers

	//ensure it's all in alphabetical order
	channels := []rune{'B', 'G', 'R'} //
	//check to see if there's any transparent spaces
	var alpha bool
	for y := 0; y < b.Y; y++ {
		for x := 0; x < b.X; x++ {
			_, _, _, a := in.At(x, y).RGBA()
			if a != 0xffff {
				alpha = true
				break
			}
		}
	} //if there is alpha add to the list of channels
	if alpha {
		alpha := []rune{'A'}
		channels = append(alpha, channels...)
	}

	chans := make([]byte, len(channels)*18+1) //generate the channel metadata
	for i, c := range channels {
		pos := i * 18
		chans[pos] = byte(c)
		chans[pos+1] = byte(terminator)
		//0 for uint32 bit and unsigned char and resrved
		//is 1 as we are using linear data
		dir.PutUint64(chans[pos+2:pos+10:pos+10], uint64(1))
		dir.PutUint32(chans[pos+10:pos+14:pos+14], 1)
		dir.PutUint32(chans[pos+14:pos+18:pos+18], 1)
	}
	chans[len(chans)-1] = byte(terminator)
	//the window is the canvas size
	window := make([]byte, 16)
	dir.PutUint32(window[8:12], uint32(b.X-1))  //the first 8 values are
	dir.PutUint32(window[12:16], uint32(b.Y-1)) //coordiantes of 0,0

	pixRatio := make([]byte, 4)
	dir.PutUint32(pixRatio, math.Float32bits(1)) //the ratio is always one

	//assign the chromacities in sRGB colour space
	var chromacities []byte
	//r,g,b,whitepoint x then y
	points := []float32{0.64, 0.33, 0.3, 0.6, 0.15, 0.06, 0.3127, 0.3290}
	for _, p := range points {
		bits := dir.AppendUint32([]byte{}, math.Float32bits(p))
		chromacities = append(chromacities, bits...)
	}
	//list of headers to be used for the exr file
	heads := []headers{
		{"channels", "chlist", chans},
		{"compression", "compression", []byte{0x00}},
		{"chromaticities", "chromaticities", chromacities},
		{"dataWindow", "box2i", window},
		{"displayWindow", "box2i", window},
		{"lineOrder", "lineOrder", []byte{0x00}},
		{"pixelAspectRatio", "float", pixRatio},
		{"screenWindowCenter", "v2f", []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"screenWindowWidth", "float", pixRatio}, //used pixel ratio as both are float 1s
	}
	var buf bytes.Buffer
	//loop through heads and write
	for _, h := range heads {
		//data as bytes
		size := dir.AppendUint32([]byte{}, uint32(len(h.data)))

		buf.Write([]byte(h.name))
		buf.Write([]byte{terminator})
		buf.Write([]byte(h.dataTag))
		buf.Write([]byte{terminator})
		buf.Write(size)
		buf.Write(h.data)
	}

	head = append(head, buf.Bytes()...)
	head = append(head, byte(terminator)) //add the terminator to signal end of the headers

	offset := uint64(len(head) + b.Y*(8))

	dataSize := len(channels) * 2 // each channel has 2 bytes for 16 bit floats

	lineLength := uint64(b.X * dataSize)

	//first 8 is num then next
	step := 8 + lineLength //bytes + len per y
	//add the offset for each y line
	for y := uint64(0); int(y) < b.Y; y++ {
		head = append(head, dir.AppendUint64([]byte{}, step*y+offset)...)
	}

	//write the header
	w.Write(head)

	pixSize := uint32(dataSize) * uint32(b.X)

	var data []byte
	switch in.(type) {
	case *image.NRGBA64:
		i := in.(*image.NRGBA64)
		data = saveNRGBA64(i, pixSize, alpha)
	case *aces.ARGBA:
		i := in.(*aces.ARGBA)
		data = saveARGBA128(i, pixSize, alpha)
	default:
		return fmt.Errorf("Image of unkown type only images of type *image.NRGBA64 and *aces.ARGBA can be saved")
	}

	//write the image bytes
	w.Write(data)

	return nil

}

type headers struct {
	name    string
	dataTag string
	data    []byte
}

func saveNRGBA64(base *image.NRGBA64, pixSize uint32, alpha bool) []byte {
	dir := binary.LittleEndian
	b := base.Bounds().Max
	var body bytes.Buffer

	s := base.Stride
	for y := uint32(0); int(y) < b.Y; y++ {
		//write the y pos and data length
		body.Write(dir.AppendUint32([]byte{}, y))
		body.Write(dir.AppendUint32([]byte{}, pixSize))
		pixLine := base.Pix[int(y)*s : int(y)*s+s : int(y)*s+s]
		r := make([]byte, b.X*2)
		g := make([]byte, b.X*2)
		bc := make([]byte, b.X*2)
		a := make([]byte, b.X*2)
		for i := 0; i < len(pixLine); i += 8 {
			//make all uint16 from floats etc

			rF := float16.Fromfloat32(float32(uint16(pixLine[i])<<8|uint16(pixLine[i+1])) / 65535)
			gF := float16.Fromfloat32(float32(uint16(pixLine[i+2])<<8|uint16(pixLine[i+3])) / 65535)
			bF := float16.Fromfloat32(float32(uint16(pixLine[i+4])<<8|uint16(pixLine[i+5])) / 65535)
			aF := float16.Fromfloat32(float32(uint16(pixLine[i+6])<<8|uint16(pixLine[i+7])) / 65535)
			// add the colour to their respective channel
			pos := (i / 8)
			dir.PutUint16(r[pos*2:pos*2+2:pos*2+2], rF.Bits())
			dir.PutUint16(bc[pos*2:pos*2+2:pos*2+2], bF.Bits())
			dir.PutUint16(g[pos*2:pos*2+2:pos*2+2], gF.Bits())
			dir.PutUint16(a[pos*2:pos*2+2:pos*2+2], aF.Bits())

		}
		if alpha {
			body.Write(a)
		}
		body.Write(bc)
		body.Write(g)
		body.Write(r)
	}

	return body.Bytes()
}

func saveARGBA128(base *aces.ARGBA, pixSize uint32, alpha bool) []byte {
	dir := binary.LittleEndian
	b := base.Bounds().Max
	var body bytes.Buffer

	s := base.Stride
	for y := uint32(0); int(y) < b.Y; y++ {
		//write the y pos and data length
		body.Write(dir.AppendUint32([]byte{}, y))
		body.Write(dir.AppendUint32([]byte{}, pixSize))
		pixLine := base.Pix[int(y)*s : int(y)*s+s : int(y)*s+s]
		r := make([]byte, b.X*2)
		g := make([]byte, b.X*2)
		bc := make([]byte, b.X*2)
		a := make([]byte, b.X*2)
		for i := 0; i < len(pixLine); i += 16 {
			//make all half floats from the bits
			rF := float16.Fromfloat32(math.Float32frombits((uint32(pixLine[i])<<24 | uint32(pixLine[i+1])<<16 | uint32(pixLine[i+2])<<8 | uint32(pixLine[i+3]))) / 65535)
			gF := float16.Fromfloat32(math.Float32frombits((uint32(pixLine[i+4])<<24 | uint32(pixLine[i+5])<<16 | uint32(pixLine[i+6])<<8 | uint32(pixLine[i+7]))) / 65535)
			bF := float16.Fromfloat32(math.Float32frombits((uint32(pixLine[i+8])<<24 | uint32(pixLine[i+9])<<16 | uint32(pixLine[i+10])<<8 | uint32(pixLine[i+11]))) / 65535)
			aF := float16.Fromfloat32(math.Float32frombits((uint32(pixLine[i+12])<<24 | uint32(pixLine[i+13])<<16 | uint32(pixLine[i+14])<<8 | uint32(pixLine[i+15]))) / 65535)
			//add them to their respective channels
			pos := (i / 16)
			dir.PutUint16(r[pos*2:pos*2+2:pos*2+2], rF.Bits())
			dir.PutUint16(bc[pos*2:pos*2+2:pos*2+2], bF.Bits())
			dir.PutUint16(g[pos*2:pos*2+2:pos*2+2], gF.Bits())
			dir.PutUint16(a[pos*2:pos*2+2:pos*2+2], aF.Bits())

		}
		if alpha {
			body.Write(a)
		}
		body.Write(bc)
		body.Write(g)
		body.Write(r)
	}

	return body.Bytes()
}
