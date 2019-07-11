package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"time"

	pdf "github.com/jung-kurt/gofpdf"
	qr "github.com/skip2/go-qrcode"

	"github.com/GoKillers/libsodium-go/cryptobox"

	hdw "github.com/sporkins/hdw-go"
	kms "github.com/sporkins/kms-go"
)

var rawStdEncoding = base64.StdEncoding

func usage() {
	fmt.Println("Creates HD wallet accounts using a random mnemonic, or, when set, using the given '--mnemonic' value")
	flag.PrintDefaults()
	fmt.Println("\n  example")
	fmt.Println(`    % ./hdw --coin 175 --password "secret"`)
	fmt.Println(`    % ./hdw --coin 175 --mnemonic "wire own magic faint cabin ranch palm property tourist riot clarify tomorrow cruise open symptom"	`)
}

func main() {
	mnemonicInput := flag.String("mnemonic", "", "A BIP-39 mnemonic to use, or a randomly generated one when n	ot set")
	coin := flag.Int("coin", 0, "Coin type used in derivation path, default 0")
	account := flag.Int("account", 0, "Account to use for derivation, default 0")

	kmsResourceID := flag.String("kms-resource-id", "", "kms resource used to encrypt key data, if not passed, will print raw data")
	kmsVersionID := flag.Int("kms-key-version", 1, "The version of the key to use, default 1 (used on ly if kms-resource-id is passed)")

	boxPkBase64 := flag.String("box-pk", "", "base64 encoded cryptobox public key to encrypt mnemonic")
	name := flag.String("name", "", "name of key")
	keyName := flag.String("keyname", "", "key name of key")

	printMnemonic := flag.Bool("print-mnemonic", false, "Boolean that controls if the mnemonic will be printed or not")

	flag.Parse()
	if *coin < 0 {
		fmt.Println("Coin must be greater than zero")
		usage()
		os.Exit(1)
	}
	if *account < 0 {
		fmt.Println("Account must be greater than zero")
		usage()
		os.Exit(1)
	}
	var acc hdw.Account
	if *mnemonicInput == "" {
		acc = hdw.Generate("", *account, *coin)
	} else {
		acc = hdw.GenerateFromMnemonic("", *mnemonicInput, *account, *coin)
	}

	kmsEncData := kmsEncrypt(acc, *kmsResourceID, *kmsVersionID)
	boxSealEncData := boxSealEcrypt(acc, *boxPkBase64)

	if kmsEncData != "" {
		fmt.Printf("\nKMS encrypted key data:\n\n%s\n", kmsEncData)
	}
	if boxSealEncData != "" {
		qrData := hdw.JSON(QRData{
			Type:                       "asset",
			CipherTextBase64:           boxSealEncData,
			DistributorPublicKeyBase64: *boxPkBase64,
			Key:                        *keyName,
			Name:                       *name,
		})

		qrFile := writeQR(*name, string(qrData))
		fmt.Printf("\nQR file %s\n\n", qrFile)
	}

	acc.PrintDerived(0, 10)

	if *printMnemonic {
		fmt.Printf("mnemonic:\n\n%s\n", acc.Mnemonic)
	}

}

func base64Encode(b []byte) []byte {
	return []byte(rawStdEncoding.EncodeToString(b))
}

func boxSeal(message []byte, boxPkBase64 string) []byte {
	boxPk, err := rawStdEncoding.DecodeString(boxPkBase64)
	if err != nil {
		println(fmt.Sprintf("Decoding boxSealPkBase64: %s", err))
		os.Exit(-1)
	}

	c, exit := cryptobox.CryptoBoxSeal(message, boxPk)
	if exit != 0 {
		println(fmt.Sprintf("CryptoBoxSeal failed: %d", exit))
		os.Exit(-1)
	}
	return c
}

func kmsEncrypt(account hdw.Account, kmsResourceID string, kmsVersion int) string {
	var kmsEncKeyData string
	if kmsResourceID != "" {
		kmsEncKeyData = string(base64Encode(kmsClient(kmsResourceID, kmsVersion).Encrypt(hdw.JSON(account.AccountJSON()))))
	}
	return kmsEncKeyData
}

func boxSealEcrypt(account hdw.Account, boxPkBase64 string) string {
	var boxSealEncKeyData string
	if boxPkBase64 != "" {
		boxSealEncKeyData = string(base64Encode(boxSeal(hdw.JSON(account.AccountJSON()), boxPkBase64)))
	}
	return boxSealEncKeyData
}

func kmsClient(kmsResourceID string, kmsVersion int) kms.KMSClient {
	return kms.NewKMSClient(fmt.Sprintf("%s/cryptoKeyVersions/%d", kmsResourceID, kmsVersion))
}

func writeQR(title string, content string) string {
	var rdr *bytes.Reader
	t := time.Now().Unix()
	qrFileName := fmt.Sprintf("%d_qr.pdf", t)

	qrBytes, _ := qr.Encode(content, qr.Low, 300)
	rdr = bytes.NewReader(qrBytes)

	qrPdf := pdf.New("P", "pt", "letter", "")
	qrPdf.AddPage()
	qrPdf.SetFont("Times", "", 12)
	qrPdf.SetTextColor(0, 0, 0)
	_, lineH := qrPdf.GetFontSize()
	qrPdf.Write(lineH, title)
	qrPdf.Ln(35)
	qrPdf.RegisterImageOptionsReader("qr", pdf.ImageOptions{ImageType: "png", ReadDpi: false}, rdr)
	qrPdf.ImageOptions("qr", 30, 0, 300, 300, true, pdf.ImageOptions{ImageType: "png", ReadDpi: false}, 0, "")
	qrPdf.OutputFileAndClose(qrFileName)
	return qrFileName
}

//QRData QR Data
type QRData struct {
	Type                       string `json:"type"`
	DistributorPublicKeyBase64 string `json:"distributorPublicKeyBase64"`
	CipherTextBase64           string `json:"cipherTextBase64"`
	Key                        string `json:"key"`
	Name                       string `json:"name"`
}
