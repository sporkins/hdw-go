package sporkins

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

const derivationPath = "m/44'/%d'/0'/0/%d"

//GenerateFromMnemonic generate account from mnenomicInput for account and coin
func GenerateFromMnemonic(password string, mnemonicInput string, account int, coin int) Account {
	mnemonic := fromMnemonic(mnemonicInput)
	seed := bip39.NewSeed(mnemonic, "")
	return generateAccount(mnemonic, seed, NetworkParams(coin), uint32(account))
}

//Generate generate account with random mnemonic for account and coin
func Generate(password string, account int, coin int) Account {
	mnemonic := generateMnemonic()
	seed := bip39.NewSeed(mnemonic, password)
	return generateAccount(mnemonic, seed, NetworkParams(coin), uint32(account))
}

func generateMnemonic() string {
	var mnemonic string
	entropy, err := bip39.NewEntropy(256)
	checkError(err, "Error generating entropy")
	mnemonic, err = bip39.NewMnemonic(entropy)
	checkError(err, "Error generating mnemonic")
	return mnemonic
}

func fromMnemonic(mnemonic string) string {
	_, err := bip39.EntropyFromMnemonic(mnemonic)
	checkError(err, "Error generating entropy")
	return mnemonic
}

//GenerateAccount generate account
func generateAccount(mnemonic string, seed []byte, network *chaincfg.Params, account uint32) Account {
	var key *hdkeychain.ExtendedKey
	masterKey, _ := hdkeychain.NewMaster(seed, network)
	key, _ = masterKey.Child(hdkeychain.HardenedKeyStart + 44)
	key, _ = key.Child(hdkeychain.HardenedKeyStart + network.HDCoinType)
	key, _ = key.Child(hdkeychain.HardenedKeyStart + account)

	return Account{
		Mnemonic:   mnemonic,
		Seed:       seed,
		MasterKey:  masterKey,
		AccountKey: key,
		Coin:       network.HDCoinType,
		Net:        network,
	}
}

//Derive derive change/index path
func (a Account) Derive(change uint32, index uint32) DerivedKey {
	changeKey, _ := a.AccountKey.Child(change)
	indexKey, _ := changeKey.Child(index)
	address, _ := indexKey.Address(a.Net)
	return DerivedKey{
		Account: a,
		Change:  change,
		Index:   index,
		Address: address.String(),
		Derived: indexKey,
	}
}

//AccountPK return the account extended public key
func (a Account) AccountPK() string {
	acctPk, err := a.AccountKey.Neuter()
	checkError(err, "error neutering account")
	return acctPk.String()
}

//AccountSK return the account extended private key
func (a Account) AccountSK() string {
	return a.AccountKey.String()
}

//AccountJSON convert to JSON
func (a Account) AccountJSON() AccountJSON {
	acctPk, err := a.AccountKey.Neuter()
	checkError(err, "error neutering account")
	account := AccountJSON{
		Mnemonic:  a.Mnemonic,
		Seed:      hex.EncodeToString(a.Seed),
		RootSK:    a.MasterKey.String(),
		AccountPK: acctPk.String(),
		AccountSK: a.AccountSK(),
	}
	return account
}

//JSON convert to JSON
func JSON(k interface{}) []byte {
	JSON, err := json.MarshalIndent(k, "", " ")
	checkError(err, "error marshalling json")
	return JSON
}

//PrintDerived derive accounts
func (a Account) PrintDerived(change uint32, count int) {
	println("\nAddresses:")
	for i := 0; i < count; i++ {
		var d = a.Derive(uint32(change), uint32(i))
		println(fmt.Sprintf("\tm/44'/%d'/%d'/%d/%d\t%s", a.Coin, a.Account, change, i, d.Address))
	}
}

func checkError(err error, msg string) {
	if err != nil {
		println(fmt.Sprintf("%s, error = %s", msg, err))
		os.Exit(-1)
	}
}
