package dpx

type dpxHeaders struct {
	//File Information Header
	MagicNumber            string `json:"Magic Number"`
	Offset                 uint32 `json:"Offset"`
	HeaderVersion          string `json:"Header Version"`
	FileSize               uint32 `json:"File Size"`
	DittoKey               uint32 `json:"Ditto Key"`
	GenericSectionHeader   uint32 `json:"Generic Section Header"`
	IndustrySpecificHeader uint32 `json:"Industry Specific Header"`
	UserDefinedHeader      uint32 `json:"User Defined Header"`
	ImageFileName          string `json:"Image File Name"`
	CreationDate           string `json:"Creation Date"`
	Creator                string `json:"Creator"`
	ProjectName            string `json:"Project Name"`
	Copyright              string `json:"Copyright"`
	EncryptionKey          uint32 `json:"EncryptionKey"`
	TBD                    string `json:"-"` //hide this but keep it in the mixer
	//Image Information Header
	ImageOrientation uint16 `json:"Image orientation"`
	NumberOFImages   uint16 `json:"Number of image elements"`
	PixelsPerLine    uint32 `json:"Pixels Per Line"`
	LinesPerElement  uint32 `json:"Lines per image element"`
	//Image Information Header for the 8 image layers
	DataStructure1 imageInfo `json:" Data structure for image element 1"`
	DataStructure2 imageInfo `json:" Data structure for image element 2"`
	DataStructure3 imageInfo `json:" Data structure for image element 3"`
	DataStructure4 imageInfo `json:" Data structure for image element 4"`
	DataStructure5 imageInfo `json:" Data structure for image element 5"`
	DataStructure6 imageInfo `json:" Data structure for image element 6"`
	DataStructure7 imageInfo `json:" Data structure for image element 7"`
	DataStructure8 imageInfo `json:" Data structure for image element 8"`
	TBDImageInfo   string    `json:"-"` //hide this but keep it in the mixer
	//Image Source Information Header
	Xoffset           uint32   `json:"X offset"`
	Yoffset           uint32   `json:"Y offset"`
	Xcenter           uint32   `json:"X center"`
	YCenter           uint32   `json:"Y center"`
	XogSize           uint32   `json:"X original size"`
	YogSize           uint32   `json:"Y original size"`
	SourceFile        string   `json:"Source image filename"`
	SourceImageDate   string   `json:"Source image date/time"`
	InputDeviceName   string   `json:"Input device name"`
	InputDeviceSerial string   `json:"Input device serial number"`
	BorderVal         []uint16 `json:"Border validity (XL,XR,YT,YB)border"`
	PixelRatio        []uint32 `json:"Pixel aspect ratio (horizontal:vertical)"`

	XScanSize      uint32 `json:"X scanned size"`
	YScanSize      uint32 `json:"Y scanned size"`
	TBDImageSource string `json:"-"` //hide this but keep it in the mixer

	//motion picture film information header
	FilmMFG          string `json:"Film mfg"`
	FilType          string `json:"Film Type"`
	OffsetinPerf     string `json:"Offset in Perfs"`
	Prefix           string `json:"Prefix"`
	Count            string `json:"Count"`
	Format           string `json:"Format"`
	FramePosition    uint32 `json:"Frame position in sequence"`
	SequenceLength   uint32 `json:"Sequence length (frames)"`
	HeldCount        uint32 `json:"Held count (1 = default)"`
	OgFrameRate      uint32 `json:"Frame rate of original (frames/s)"`
	ShutterAngle     uint32 `json:"Shutter angle of camera in degrees"`
	FrameID          string `json:"Frame identification"`
	SlateInformation string `json:"Slate information"`
	TBDMotionPicture string `json:"-"` //hide this but keep it in the mixer

	//Television Information Header
	SMPTETime     uint32 `json:"SMPTE time code"`
	SMPTEUser     uint32 `json:"SMPTE user bits"`
	Interlace     uint8  `json:"Interlace "`
	FielNum       uint8  `json:"Field number"`
	VideoSignal   uint8  `json:"Video signal standard"`
	Zero          uint8  `json:"Zero"`
	HorizSamp     uint32 `json:"Horizontal sampling rate "`
	VertSamp      uint32 `json:"Vertical sampling rate"`
	TemporalSamp  uint32 `json:"Temporal sampling rate or frame rate"`
	TimeOffset    uint32 `json:"Time offset from sync to first pixel "`
	Gamma         uint32 `json:"Gamma"`
	BlackLevel    uint32 `json:"Black level code value"`
	BlackGain     uint32 `json:"Black gain"`
	BreakPoint    uint32 `json:"Breakpoint"`
	Reference     uint32 `json:"Reference white level code value"`
	Integration   uint32 `json:"Integration time (s)"`
	TBDTelevision string `json:"-"` //hide this but keep it in the mixer
	//User Defined
	UserID      string `json:"User Identification"`
	UserDefined string `json:"User Defined Content"`
}

//image info is the structure for the 8 image layers
type imageInfo struct {
	DataSign         uint32 `json:"Data Sign"`
	LowData          uint32 `json:"Reference Low Data"`
	LowQuanity       uint32 `json:"Reference Low Quantity"`
	HighData         uint32 `json:"Reference High Data"`
	HighQuanity      uint32 `json:"Reference High Quantity"`
	Descriptor       uint8  `json:"Descriptor"`
	TransferChar     uint8  `json:"Transfer characteristic"`
	ColorimetricSpec uint8  `json:"Colorimetric specification"`
	BitDepth         uint8  `json:"Bit Depth"`
	Packing          uint16 `json:"Packing"`
	Encoding         uint16 `json:"Encoding"`
	DataOffset       uint32 `json:"Offset to Data"`
	EOLPadding       uint32 `json:"End-of-line Padding"`
	EOIPadding       uint32 `json:"End-of-image padding"`
	ImageDescrip     string `json:"Description of Image Element"`
}

//length holds the byte length for each header, in the order they are declared in the struct
var length []int = []int{
	//File Information Header
	4, 4, 8, 4, 4, 4, 4, 4, 100, 24, 100, 200, 200, 4, 104,
	//Image Information Header
	2, 2, 4, 4,
	//Data Structure for image elements
	72, 72, 72, 72, 72, 72, 72, 72, 52,
	//Image Source Information
	4, 4, 4, 4, 4, 4, 100, 24, 32, 32, 8, 8,
	//Data structure for additional information
	4, 4, 20,
	//motion picture film information header
	2, 2, 2, 6, 4, 32, 4, 4, 4, 4, 4, 32, 100, 56,
	//Television Information
	4, 4, 1, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 76,
	//User Defined
	32, 6112} //last digit isn't used and is 0 for the time being

//imageLength holds the byte length for each image header, in the order they are declared in the struct
var imageLength []int = []int{4, 4, 4, 4, 4, 1, 1, 1, 1, 2, 2, 4, 4, 4, 32}

//Offset is the header offset in length of bytes and is the cumulative value of length
// at each point in the array

var offset []int
var imageOffset []int
