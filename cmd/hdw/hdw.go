package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	hdw "github.com/sporkins/hdw-go"
)

var rawStdEncoding = base64.StdEncoding.WithPadding(base64.NoPadding)

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
		acc = hdw.Generate(*account, *coin)
	} else {
		acc = hdw.GenerateFromMnemonic(*mnemonicInput, *account, *coin)
	}

	println(string(hdw.JSON(acc.AccountJSON())))

	acc.PrintDerived(0, 10)

	fmt.Printf("mnemonic:\n\n%s\n", acc.Mnemonic)

}

func base64Encode(b []byte) []byte {
	return []byte(rawStdEncoding.EncodeToString(b))
}
