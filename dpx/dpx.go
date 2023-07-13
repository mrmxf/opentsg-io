// Package dpx contains the encoding methods for dpx files
package dpx

import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"math"
	"reflect"
)

type Options struct {
	Bitdepth int
}

// Encode writes the image m to the file based on the bitdepth, if it is nil a 16 bit image is saved
func Encode(w io.Writer, m *image.NRGBA64, options *Options) error {

	/*
			switch m.(type) {
		case *image.NRGBA64:
			im := m.(*image.NRGBA64)
			data = saveNRGBA64(i, pixSize, alpha)

		default:
			return fmt.Errorf("Image of unkown type only images of type *image.NRGBA64 can be saved")
		}*/
	//generate the header and body
	btotal := headerGen(m, options.Bitdepth)
	//conver to io.reader and write as a stream loop through as a buffer etc
	_, err := w.Write(btotal)
	return err
}

func headerGen(canvas *image.NRGBA64, bitdepth int) []byte {
	//Magic Numb]er
	//offset
	//version number
	//total file size

	//image orientation
	//numer of image elemanets
	//pixels per line
	//lines per image

	//image 1
	//data sign
	//descriptot
	//transfer charterisitc
	//colormetric specification
	//bit depth
	//packing
	//encoding
	//offset to data

	//assign all these values
	var header dpxHeaders

	header.MagicNumber = "SDPX"
	header.Offset = 8192
	header.HeaderVersion = "V2.0"
	header.Creator = "MR MXF's golang dpx writer"
	header.EncryptionKey = 4294967295

	header.GenericSectionHeader = 1664
	header.IndustrySpecificHeader = 384
	header.UserDefinedHeader = 6144

	header.ImageOrientation = 0
	header.NumberOFImages = 1
	header.PixelsPerLine = uint32(canvas.Bounds().Max.X)
	header.LinesPerElement = uint32(canvas.Bounds().Max.Y)

	header.DataStructure1.DataSign = 0
	header.DataStructure1.Descriptor = 50 //change to 51 for rgba
	header.DataStructure1.TransferChar = 0
	header.DataStructure1.ColorimetricSpec = 0

	header.DataStructure1.Encoding = 0
	header.DataStructure1.DataOffset = 8192
	header.DataStructure1.HighData = 65335

	var ImageBuf []byte
	switch bitdepth {
	case 8:
		ImageBuf = encode8(canvas.Pix, canvas.Bounds().Max.X, canvas.Bounds().Max.Y)
		header.DataStructure1.Packing = 0
	case 10:
		ImageBuf = encode10(canvas.Pix, canvas.Bounds().Max.X, canvas.Bounds().Max.Y)
		header.DataStructure1.Packing = 1
	case 12:
		ImageBuf = encode12(canvas.Pix, canvas.Bounds().Max.X, canvas.Bounds().Max.Y)
		header.DataStructure1.Packing = 0
	case 16:
		ImageBuf = encode16(canvas.Pix, canvas.Bounds().Max.X, canvas.Bounds().Max.Y)
		header.DataStructure1.Packing = 0
	}

	bd := uint8(bitdepth)
	header.DataStructure1.BitDepth = bd
	highdata := math.Pow(2, float64(bd)) - 1

	//fmt.Println(highdata)
	header.DataStructure1.HighData = uint32(highdata)
	header.FileSize = uint32(8192 + len(ImageBuf))
	//fmt.Println(header.FileSize)
	//these are all the compulsory bits of data

	//assign extra data if required like date time etc

	b, _ := header.byteConvert()
	//fmt.Println(len(b))
	b = append(b, ImageBuf...)
	return b
}

func (dpxhead *dpxHeaders) byteConvert() ([]byte, error) {

	v := reflect.ValueOf(dpxhead).Elem()
	vtype := v.Type()

	var data bytes.Buffer

	//loop through the struct values assigning values from the bytes based
	//on the offset and length of each header
	for i := 0; i < v.NumField(); i++ {
		fieldname := vtype.Field(i).Name
		//write to a length of bytes as specified
		dint := dpxhead.fieldTobyte(fieldname, length[i])
		//update the final user defined length to access the data
		data.Write(dint)
	}
	return data.Bytes(), nil
}

func (dpxhead *dpxHeaders) fieldTobyte(fieldname string, length int) []byte {
	structVal := reflect.ValueOf(dpxhead).Elem()
	structField := structVal.FieldByName(fieldname)

	b := valueGet(structField.Interface())
	if len(b) != length {
		filler := make([]byte, length-len(b))
		b = append(b, filler...)
	}
	return b
}

func (ii *imageInfo) byteConvert() []byte {

	v := reflect.ValueOf(ii).Elem()
	vtype := v.Type()

	var data bytes.Buffer

	//loop through the struct values assigning values from the bytes based
	//on the offset and length of each header
	for i := 0; i < v.NumField(); i++ {
		fieldname := vtype.Field(i).Name
		//write to a length of bytes as specified
		dataInt := ii.iifieldToByte(fieldname, imageLength[i])
		//update the final user defined length to access the data
		data.Write(dataInt)
	}

	return data.Bytes()
}

func (ii *imageInfo) iifieldToByte(fieldname string, length int) []byte {
	structVal := reflect.ValueOf(ii).Elem()
	structField := structVal.FieldByName(fieldname)
	b := valueGet(structField.Interface())
	if len(b) != length {
		filler := make([]byte, length-len(b))
		b = append(b, filler...)
	}
	return b
}

func valueGet(vinter interface{}) []byte {
	var b []byte
	switch vinter.(type) {
	case string:
		b = []byte(vinter.(string))
	case uint32:
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, vinter.(uint32))
	case uint16:
		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, vinter.(uint16))
	case uint8:
		bint := byte(vinter.(uint8))
		b = append(b, bint)
	case []uint32:
		v := vinter.([]uint32)

		for i := 0; i < len(v); i++ {
			bint := make([]byte, 4)
			binary.BigEndian.PutUint32(bint, v[i])
			//series of integer arrays
			b = append(b, bint...)
		}
	case []uint16:
		v := vinter.([]uint16)

		for i := 0; i < len(v); i++ {
			bint := make([]byte, 2)
			binary.BigEndian.PutUint16(bint, v[i])
			//series of integer arrays
			b = append(b, bint...)
		}
	case imageInfo:
		tempIm := vinter.(imageInfo)
		//generate a struct to assign to the element
		b = tempIm.byteConvert()
	default:
		//err = fmt.Errorf("Invalid type for field ")
	}

	return b
}

// 32 bit words of 16 bit
func encode16(pix []uint8, dx, dy int) []byte {
	//stride is the image between vertically adjacent pixels
	off := 0
	//assign to 3/4 as we are ignoring the alpha channel
	buf := make([]byte, (len(pix)*3)/4)
	for y := 0; y < dy; y++ {
		min := y * dx * 8
		max := min + dx*8

		for i := min; i < max; i += 8 {
			/*
				// An image.RGBA64's Pix is in big-endian order.
				r := uint16(pix[i])<<8 | uint16(pix[i+1])
				g := uint16(pix[i+2])<<8 | uint16(pix[i+3])
				b := uint16(pix[i+4])<<8 | uint16(pix[i+5])
				//ignore alpha channel at the mo
				//a1 := uint16(pix[i+6])<<8 | uint16(pix[i+7])
				// big endian dpx files
				//change to append when getting rid of the alpha channel
				buf[off+0] = byte(b >> 8)
				buf[off+1] = byte(b)
				buf[off+2] = byte(g >> 8)
				buf[off+3] = byte(g)
				buf[off+4] = byte(r >> 8)
				buf[off+5] = byte(r)
			*/
			//bgr order direct from the pixels
			buf[off+0] = byte(pix[i])
			buf[off+1] = byte(pix[i+1])
			buf[off+2] = byte(pix[i+2])
			buf[off+3] = byte(pix[i+3])
			buf[off+4] = byte(pix[i+4])
			buf[off+5] = byte(pix[i+5])
			off += 6
		}
	}
	return buf
}

// 32 words of 8 bit
func encode8(pix []uint8, dx, dy int) []byte {
	off := 0
	rSize := float64((len(pix)*3)/4) * 0.5 //3/4 for alpha 0.5 for 8 bit
	buf := make([]byte, int(rSize))
	for y := 0; y < dy; y++ {
		min := y * dx * 8
		max := min + dx*8

		for i := min; i < max; i += 8 {
			pline := pix[i : i+7 : i+7]
			// An image.RGBA64's Pix is in big-endian order.
			r := uint16(pline[0])<<8 | uint16(pline[1])
			g := uint16(pline[2])<<8 | uint16(pline[3])
			b := uint16(pline[4])<<8 | uint16(pline[5])
			/*
					r := uint16(pix[i])<<8 | uint16(pix[i+1])
				g := uint16(pix[i+2])<<8 | uint16(pix[i+3])
				b := uint16(pix[i+4])<<8 | uint16(pix[i+5])
			*/
			//ignore alpha channel at the mo
			//a1 := uint16(pix[i+6])<<8 | uint16(pix[i+7])
			// big endian dpx files
			//change to append when getting rid of the alpha channel
			buf[off+0] = byte(b >> 8)
			buf[off+1] = byte(g >> 8)
			buf[off+2] = byte(r >> 8)
			off += 3
		}
	}
	return buf
}

func encode10(pix []uint8, dx, dy int) []byte {

	ppl := math.Ceil(float64(dx) * 4)
	off := 0
	buf := make([]byte, int(ppl)*dy)
	// go through each line and assign pixels for each one
	for y := 0; y < dy; y++ {
		min := y * dx * 8
		max := min + dx*8
		for i := min; i < max; i += 8 {

			r := uint16(pix[i])<<8 | uint16(pix[i+1])
			g := uint16(pix[i+2])<<8 | uint16(pix[i+3])
			b := uint16(pix[i+4])<<8 | uint16(pix[i+5])
			//shift to 10 bit equivalents etc
			r = r >> 6
			g = g >> 6
			b = b >> 6

			//full 32 word of 00bgr with packing
			buf[off] = byte(r >> 2)        //8 r
			buf[off+1] = byte(r<<6 | g>>4) //2 r //6g
			buf[off+2] = byte(g<<4 | b>>6) //4g //4b
			buf[off+3] = byte(b << 2 & 0b11111100)

			off = off + 4
		}
	}
	return buf
}

func encode12(pix []uint8, dx, dy int) []byte {

	ppl := math.Ceil(float64(dx) * 4.5)
	if int(ppl)%8 != 0 {
		ppl = float64(int(ppl) + 8 - int(ppl)%8)
	}
	//fmt.Println(ppl)
	rSize := float64(ppl * float64(dy))
	buf := make([]byte, int(rSize))
	off := 0

	for y := 0; y < dy; y++ {
		min := y * dx * 8
		max := min + dx*8
		//64 as 8 lots of rgb for a full cycle of packing the 32 bit word
		for i := min; i < max; i += 64 {
			// An image.RGBA64's Pix is in big-endian order.
			var endpoint int
			r, g, b := rgbGet(pix[i : i+6 : i+6])
			r1, g1, b1 := rgbGet(pix[i+8 : i+14 : i+14])
			r2, g2, b2 := rgbGet(pix[i+16 : i+22 : i+22])
			r3, g3, b3 := rgbGet(pix[i+24 : i+30 : i+30])
			r4, g4, b4 := rgbGet(pix[i+32 : i+38 : i+38])
			r5, g5, b5 := rgbGet(pix[i+40 : i+46 : i+46])
			r6, g6, b6 := rgbGet(pix[i+48 : i+54 : i+54])
			r7, g7, b7 := rgbGet(pix[i+56 : i+62 : i+62])

			//break the bytes into their respective pixels and start a new one on a new row
			if (i + 64) > max { //cut out the last 32 bit word
				//set an integer break point
				//this is to fill the rest of the words with 0 for this tyle of packing
				endpoint = (max - i) / 8 //from 1-7 if it's over 64

			}

			buf[off] = byte(b)             //last 8 bits of datum 2
			buf[off+1] = byte(g >> 4)      //first 8 bits of datum 1
			buf[off+2] = byte(g<<4 | r>>8) //last 4 bits of datum 1 //first 4 of datum 0
			buf[off+3] = byte(r)           //last 8 bits of datum 0

			if endpoint == 1 {
				buf[off+4] = byte(0b00000000) //last 4 of 5 first 4 of 4
				buf[off+5] = byte(0b00000000) //last 8 of 4
				buf[off+6] = byte(0b00000000) //first 8 of 3
				buf[off+7] = byte(0b0000 | b>>8)
				off = off + 8
				break
			}
			buf[off+4] = byte(b1<<4 | g1>>8) //last 4 of 5 first 4 of 4
			buf[off+5] = byte(g1)            //last 8 of 4
			buf[off+6] = byte(r1 >> 4)       //first 8 of 3
			buf[off+7] = byte(r1<<4 | b>>8)  ////last 4 of 3first 4 bits of datum 2

			buf[off+8] = byte(g2 >> 4)       // last 8 of 7
			buf[off+9] = byte(g2<<4 | r2>>8) //last 4 of 7 first 4 of 6
			buf[off+10] = byte(r2)           //last 8 of 6
			buf[off+11] = byte(b1 >> 4)      // first 8 of 5

			if endpoint == 2 {
				buf[off+8] = byte(0b00000000) // last 8 of 7
				buf[off+9] = byte(0b00000000) //last 4 of 7 first 4 of 6
				buf[off+10] = byte(0b00000000)
				off = off + 12
				break
			}
			buf[off+12] = byte(g3)            //last 8 of 10
			buf[off+13] = byte(r3 >> 4)       //first 8 of 9
			buf[off+14] = byte(r3<<4 | b2>>8) //last 4 of 9 first 4 of 8
			buf[off+15] = byte(b2)            //last 8 of 8

			if endpoint == 3 {
				buf[off+12] = byte(0b00000000)
				buf[off+13] = byte(0b00000000)
				buf[off+14] = byte(0b00000000 | b2>>8)

				off = off + 16
				break
			}

			buf[off+16] = byte(g4<<4 | r4>>8) //last 4 of 13 first 4 of 12
			buf[off+17] = byte(r4)            //last 8 of 12
			buf[off+18] = byte(b3 >> 4)       //first 8 of 11
			buf[off+19] = byte(b3<<4 | g3>>8) ////last 4 of 11first 4 bits of 10

			if endpoint == 4 {
				buf[off+16] = byte(0b00000000)
				buf[off+17] = byte(0b00000000)
				off = off + 20
				break
			}
			buf[off+20] = byte(r5 >> 4)       // last 8 of 15
			buf[off+21] = byte(r5<<4 | b4>>8) //last 4 of 15 first 4 of 14
			buf[off+22] = byte(b4)            //last 8 of 14
			buf[off+23] = byte(g4 >> 4)       // first 8 of 13

			if endpoint == 5 {
				buf[off+20] = byte(0b00000000) // last 8 of 15
				buf[off+21] = byte(0b0000 | b4>>8)
				off = off + 24
				break
			}

			buf[off+24] = byte(r6)            //last 8 bits of datum 2
			buf[off+25] = byte(b5 >> 4)       //first 8 bits of datum 1
			buf[off+26] = byte(b5<<4 | g5>>8) //last 4 bits of datum 1 //first 4 of datum 0
			buf[off+27] = byte(g5)            //last 8 bits of datum 0

			if endpoint == 6 {
				buf[off+24] = byte(0b00000000)
				off = off + 28
				break
			}

			buf[off+28] = byte(r7<<4 | b6>>8) //last 4 of 5 first 4 of 4
			buf[off+29] = byte(b6)            //last 8 of 4
			buf[off+30] = byte(g6 >> 4)       //first 8 of 3
			buf[off+31] = byte(g6<<4 | r6>>8) ////last 4 of 3first 4 bits of datum 2

			if endpoint == 7 {
				buf[off+28] = byte(0b0000 | b6>>8)
				off = off + 32
			} else {
				buf[off+32] = byte(b7 >> 4)       // last 8 of 7
				buf[off+33] = byte(b7<<4 | g7>>8) //last 4 of 7 first 4 of 6
				buf[off+34] = byte(g7)            //last 8 of 6
				//fmt.Println(i, y, off)
				buf[off+35] = byte(r7 >> 4) // first 8 of 5
				off = off + 36
			}
		}
	}
	//fmt.Println(len(buf))
	//fmt.Println(pix[0:200])
	//fmt.Println(buf[252:300], "yes")
	return buf //x[8192:]
}

// only for generating 12bit images
func rgbGet(pix []uint8) (r, g, b uint16) {
	r = uint16(pix[0])<<8 | uint16(pix[1])
	g = uint16(pix[2])<<8 | uint16(pix[3])
	b = uint16(pix[4])<<8 | uint16(pix[5])

	r = r >> 4
	g = g >> 4
	b = b >> 4
	return
}
