package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	hdw "github.com/sporkins/hdw-go"
	kms "github.com/sporkins/kms-go"
)

var rawStdEncoding = base64.StdEncoding.WithPadding(base64.NoPadding)

func usage() {
	fmt.Println("Creates HD wallet accounts using a random mnemonic, or, when set, using the given '--mnemonic' value")
	flag.PrintDefaults()
	fmt.Println("\n  example")
	fmt.Println(`    % ./hdgenerator --coin 175 --password "secret"`)
	fmt.Println(`    % ./hdgenerator --coin 175 --password "secret" --change 0 --index 0`)
}

func main() {
	password := flag.String("password", "", "The password used for the BIP-39 seed")
	mnemonicInput := flag.String("mnemonic", "", "A BIP-39 mnemonic to use, or a randomly generated one when not set")
	coin := flag.Int("coin", 0, "Coin type used in derivation path")
	account := flag.Int("account", 0, "Account to use for derivation")
	network := flag.String("network", "mainnet", "chain network to use, either mainnet or testnet")
	kmsResourceID := flag.String("kms-resource-id", "", "kms resource used to encrypt key data")
	kmsVersionID := flag.Int("kms-key-version", 1, "The version of the key to use")
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
	if !hdw.Networks[*network] {
		fmt.Println("network must be either mainnet or testnet")
		usage()
		os.Exit(1)
	}

	var mnemonic hdw.Mnemonic
	if *mnemonicInput == "" {
		mnemonic = hdw.GenerateMnemonic()
	} else {
		mnemonic = hdw.FromMnemonic(*mnemonicInput)
	}

	acc := mnemonic.GenerateSeed(*password).
		GenerateMasterKey(hdw.NetworkParams(*coin, *network)).
		Account(uint32(*coin), uint32(*account))
	print(acc, *kmsResourceID, *kmsVersionID)
}

func print(account hdw.Account, kmsResourceID string, kmsVersion int) {
	accountJSON := account.AccountJSON()

	if kmsResourceID == "" {
		println(string(accountJSON))
	} else {
		kmsClient := kms.NewKMSClient(fmt.Sprintf("%s/cryptoKeyVersions/%d", kmsResourceID, kmsVersion))
		encryptdAccountJSON := kmsClient.Encrypt(accountJSON)
		println(fmt.Sprintf("encrypted JSON base64:\t%s", rawStdEncoding.EncodeToString(encryptdAccountJSON)))
	}

	account.PrintDerived(0, 10)

}
