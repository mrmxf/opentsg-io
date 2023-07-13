package dpx

/*

type dpxDecode struct {
	r         io.ReaderAt
	byteOrder binary.ByteOrder
	config    image.Config

	orien   uint16
	elemNum uint16

	desc     uint8
	bitdepth uint8

	buf   []byte
	off   int    // Current offset in the bytes.
	v     uint32 // Buffer value for reading with arbitrary bit depths.
	nbits uint   // Remaining number of bits in v.
}

func Decode(r io.Reader) (image.Image, error) {
	d, err := newDpxDecode(r)
	//get the decoded info from the file
	if err != nil {
		return nil, err
	}
	//imgRect := image.Rect(0, 0, d.config.Width, d.config.Height)
	//then extract the info in someway and slap it in here
	if d.desc != 50 {
		return nil, fmt.Errorf("only rgb files are currently opened")
	}

	if (d.bitdepth != 8) && (d.bitdepth != 16) {
		return nil, fmt.Errorf("Bitdepths of %b are not supported", d.bitdepth)
	}
	return nil, nil
}

func newDpxDecode(r io.Reader) (*dpxDecode, error) {
	dpxF := &dpxDecode{
		r: r.(io.ReaderAt),
	}
	//get the FIH
	h := make([]byte, 2080)
	if _, err := dpxF.r.ReadAt(h, 0); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	switch string(h[0:4]) {
	case "SDPX":
		dpxF.byteOrder = binary.BigEndian
	case "XPDS":
		dpxF.byteOrder = binary.LittleEndian
	default:
		return nil, fmt.Errorf("Not a valid dpx file")
	}

	imageOffset := dpxF.byteOrder.Uint32(h[4:8])
	fileSize := dpxF.byteOrder.Uint32(h[16:20])

	dpxF.orien = dpxF.byteOrder.Uint16(h[768:770])
	dpxF.elemNum = dpxF.byteOrder.Uint16(h[770:772])

	dpxF.config.Width = int(dpxF.byteOrder.Uint32(h[772:776]))
	dpxF.config.Height = int(dpxF.byteOrder.Uint32(h[776:780]))

	//convert between []byte and a single uint8
	desc := h[800:801]
	dpxF.desc = desc[0]
	bitdpth := h[803:804]
	dpxF.bitdepth = bitdpth[0]
	fmt.Println(h[780:806])

	dataStructOff := dpxF.byteOrder.Uint32(h[808:812])

	// check it's got the numbe
	// assign the offset to the decoder
	// get the bytes

	h = make([]byte, fileSize-imageOffset)
	if _, err := dpxF.r.ReadAt(h, int64(dataStructOff)); err != nil {
		return nil, err
	}
	dpxF.buf = h

	return dpxF, nil
}
*/
