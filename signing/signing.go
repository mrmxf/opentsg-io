//Package signing provides the methods to save a signature of a given file
package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

//Message sign takes message in string form, generates a sha256 of it
// and then encodes it using a private key. fname is used to identify the file being signed
//, so that the signaute name reflects that of the file.
func MessageSign(mes, fname string) error {
	//get the key and extract the data we can use

	keyDat, err := os.ReadFile("./private.pem")
	if err != nil {
		return err
	}
	block, _ := pem.Decode(keyDat)
	key, err := x509.DecryptPEMBlock(block, []byte("mrmxf"))
	if err != nil {
		return err
	}
	//change the block to key
	priv, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return err
	}

	//hash the data and genereate the signeature with it
	hashed := sha256.Sum256([]byte(mes))
	sig, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA256, hashed[:])

	if err != nil {
		return fmt.Errorf("Error from signing: %s", err)
	}
	//save in a sha256 file as binary file
	l, _ := os.Create(fname + ".sha256")
	defer l.Close()
	l.Write(sig)

	return nil
}
