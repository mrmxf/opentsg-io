//Package tiffup saves images as tiff files without the alpha channel
// This is based off of the golang tiff image library
package tiffup

import (
	"encoding/binary"
	"image"
	"io"
)

//Encode Saves the image to the specified writer as an 16 bit image with no alpha channel
func Encode(w io.Writer, img *image.NRGBA64) error {
	//write the input for all this
	// ifd tags
	ifd(w, img)
	return nil
}

func encodeRGB(w io.Writer, pix []uint8, dx, dy int) error {
	buf := make([]byte, dy*dx*6)
	off := 0
	for y := 0; y < dy; y++ {
		min := y * dx * 8
		max := min + dx*8
		//increase by 8 each time to compensate for the alpha channel that is skipped
		for i := min; i < max; i += 8 {

			//keep the big endian format of pix
			//just remove the alpha channel
			buf[off+0] = byte(pix[i+0])
			buf[off+1] = byte(pix[i+1])
			buf[off+2] = byte(pix[i+2])
			buf[off+3] = byte(pix[i+3])
			buf[off+4] = byte(pix[i+4])
			buf[off+5] = byte(pix[i+5])

			off += 6

		}

	}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

//Big endian files are written
var beMagic = "MM\x00\x2A"
var enc = binary.BigEndian

func ifd(w io.Writer, img *image.NRGBA64) error {
	d := img.Bounds().Size()

	if _, err := io.WriteString(w, beMagic); err != nil {
		return err
	}

	// imageLen is total number of image bytes
	imageLen := d.X * d.Y * 6

	if err := binary.Write(w, enc, uint32(8+imageLen)); err != nil {
		return err
	}
	//add the image after the ifd
	encodeRGB(w, img.Pix[:], d.X, d.Y)

	photometricInterpretation := uint32(2) //prgb is 2
	samplesPerPixel := uint32(3)
	bitsPerSample := []uint32{16, 16, 16}

	ifd := []ifdEntry{
		{tImageWidth, dtShort, []uint32{uint32(d.X)}},
		{tImageLength, dtShort, []uint32{uint32(d.Y)}},
		{tBitsPerSample, dtShort, bitsPerSample},
		{tPhotometricInterpretation, dtShort, []uint32{photometricInterpretation}},
		{tStripOffsets, dtLong, []uint32{8}},
		{tSamplesPerPixel, dtShort, []uint32{samplesPerPixel}},
		{tRowsPerStrip, dtShort, []uint32{uint32(d.Y)}},
		{tStripByteCounts, dtLong, []uint32{uint32(d.X * d.Y * 6)}},
		// There is currently no support for storing the image
		// resolution, so give a bogus value of 72x72 dpi.
		{tXResolution, dtRational, []uint32{72, 1}}, //denominator and numerator
		{tYResolution, dtRational, []uint32{72, 1}},
		{tResolutionUnit, dtShort, []uint32{2}},
	}
	//12 bytes per ifd
	var bufifd [12]byte
	//length after tags
	pointers := len(ifd) * 12
	var pointArea []byte

	cBuf := make([]byte, 2)
	enc.PutUint16(cBuf[:], uint16(len(ifd)))
	w.Write(cBuf)
	var off uint32
	//8 for initial tag
	//2 for total tags
	//4 for the position of next ifd at the end of the sequence
	totalOff := uint32(imageLen + 8 + 2 + 4 + pointers)

	for _, ent := range ifd {
		//assign tag then the data type values
		enc.PutUint16(bufifd[0:2], uint16(ent.tag))
		enc.PutUint16(bufifd[2:4], uint16(ent.datatype))

		//check if the data fits in the 8:12 bytes
		length := len(ent.data)
		datLen := uint32(length) * lengths[ent.datatype]

		if datLen <= 4 {
			enc.PutUint32(bufifd[4:8], uint32(length)) //count of values
			ent.putData(bufifd[8:12])
		} else {
			if ent.datatype == dtRational {
				datLen /= 2 //change rational to length of a long
				length /= 2
			}
			enc.PutUint32(bufifd[4:8], uint32(length)) //count of values

			enc.PutUint32(bufifd[8:12], totalOff+off)

			inter := make([]byte, datLen)
			ent.putData(inter)
			pointArea = append(pointArea, inter...)
			off += datLen
		}
		w.Write(bufifd[:])
	}
	//loop through and assign values
	//assign 0 as the next group of ifd values as we only add 1
	if err := binary.Write(w, enc, uint32(0)); err != nil {
		return err
	}
	//add the points for ifd over runs
	if _, err := w.Write(pointArea[:]); err != nil {
		return err
	}
	return nil
}

func (e ifdEntry) putData(p []byte) {
	for _, d := range e.data {
		switch e.datatype {
		case dtByte, dtASCII:
			p[0] = byte(d)
			p = p[1:]
		case dtShort:
			enc.PutUint16(p, uint16(d))
			p = p[2:]
		case dtLong, dtRational:
			enc.PutUint32(p, uint32(d))
			p = p[4:] //shift to the next four place in the byte
		}
	}
}

type ifdEntry struct {
	tag      int
	datatype int
	data     []uint32
}

var lengths = [...]uint32{0, 1, 1, 2, 4, 8}

// Data types (p. 14-16 of the spec).
const (
	dtByte     = 1
	dtASCII    = 2
	dtShort    = 3
	dtLong     = 4
	dtRational = 5
)

// Tags (see p. 28-41 of the spec).
const (
	tImageWidth                = 256
	tImageLength               = 257
	tBitsPerSample             = 258
	tPhotometricInterpretation = 262
	tStripOffsets              = 273
	tSamplesPerPixel           = 277
	tRowsPerStrip              = 278
	tStripByteCounts           = 279
	tXResolution               = 282
	tYResolution               = 283
	tResolutionUnit            = 296
)
