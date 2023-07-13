package signing

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSigningErrors(t *testing.T) {

	fileToSign := "./testfiles/test.dpx.txt"
	//make simple keys for each test
	/*
		priv := exec.Command("openssl", "genrsa", "-aes128", "-passout", "pass:mrmxf", "-out", "private.pem", "512")
		pub := exec.Command("openssl", "rsa", "-in", "private.pem", "-passin", "pass:mrmxf", "-pubout", "-out", "public.pem")
		priv.Output()
		pub.Output()
	*/

	file, _ := os.ReadFile(fileToSign)
	genErr := MessageSign(string(file), fileToSign)
	Convey("Checking errors are caught", t, func() {
		Convey(fmt.Sprintf("using a %s as the file to sign with no key available", fileToSign), func() {
			Convey("An error of \"open ./private.pem: no such file or directory\" is returned", func() {
				So(genErr.Error(), ShouldEqual, "open ./private.pem: no such file or directory")
			})
		})
	})

	priv := exec.Command("openssl", "genrsa", "-aes128", "-passout", "pass:fakepass", "-out", "private.pem", "512")
	priv.Output()
	file2, _ := os.ReadFile(fileToSign)
	genErrBadPass := MessageSign(string(file2), fileToSign)
	Convey("Checking errors are caught when the wrong password is used", t, func() {
		Convey(fmt.Sprintf("using a %s as the file to sign", fileToSign), func() {
			Convey("An error of \"x509: decryption password incorrect\"  is returned", func() {
				So(genErrBadPass.Error(), ShouldEqual, "x509: decryption password incorrect")
			})
		})
	})
	os.Remove("private.pem")
}
func TestSigning(t *testing.T) {

	fileToSign := []string{"./testfiles/test.dpx.txt", "./testfiles/test.tiff.txt"}
	//make simple keys for each test
	priv := exec.Command("openssl", "genrsa", "-aes128", "-passout", "pass:mrmxf", "-out", "private.pem", "512")
	pub := exec.Command("openssl", "rsa", "-in", "private.pem", "-passin", "pass:mrmxf", "-pubout", "-out", "public.pem")
	priv.Output()
	pub.Output()
	for _, sFile := range fileToSign {
		file, _ := os.ReadFile(sFile)
		genErr := MessageSign(string(file), sFile)
		Convey("Checking signature files are generated", t, func() {
			certifCheck := exec.Command("openssl", "dgst", "-sha256", "-verify", "public.pem", "-signature", sFile+".sha256", sFile) //pseudo command line
			opOut, _ := certifCheck.CombinedOutput()
			Convey(fmt.Sprintf("using a %s as the file to sign", sFile), func() {
				Convey("No error is returned amd the openssl verifies the file", func() {
					So(genErr, ShouldBeNil)
					So(string(opOut), ShouldEqual, "Verified OK\n")
				})
			})
		})
	}
	//delete the keys afterwards
	os.Remove("private.pem")
	os.Remove("public.pem")
}
