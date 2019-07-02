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

//GenerateMnemonic generates mnemonic
func GenerateMnemonic() Mnemonic {
	var mnemonic string
	entropy, err := bip39.NewEntropy(256)
	checkError(err, "Error generating entropy")
	mnemonic, err = bip39.NewMnemonic(entropy)
	checkError(err, "Error generating mnemonic")
	return Mnemonic{mnemonic, entropy}
}

//FromMnemonic generates mnemonic
func FromMnemonic(mnemonic string) Mnemonic {
	entropy, err := bip39.EntropyFromMnemonic(mnemonic)
	checkError(err, "Error generating entropy")
	return Mnemonic{mnemonic, entropy}
}

//GenerateSeed generate sedd from password
func (m Mnemonic) GenerateSeed(password string) Seed {
	seed := bip39.NewSeed(m.mnemonic, password)
	return Seed{m.mnemonic, m.entropy, seed}
}

//GenerateMasterKey for network
func (s Seed) GenerateMasterKey(network *chaincfg.Params) MasterKey {
	masterKey, err := hdkeychain.NewMaster(s.seed, network)
	checkError(err, "Error getting master key")
	return MasterKey{s.mnemonic, s.entropy, s.seed, masterKey, network}
}

//ChildH childH
func (m MasterKey) ChildH(i uint32) ChildKey {
	acc, err := m.masterKey.Child(hdkeychain.HardenedKeyStart + i)
	checkError(err, fmt.Sprintf("error i %d", i))
	return ChildKey{acc}
}

//ChildH childH
func (m ChildKey) ChildH(i uint32) ChildKey {
	acc, err := m.childKey.Child(hdkeychain.HardenedKeyStart + i)
	checkError(err, fmt.Sprintf("error i %d", i))
	return ChildKey{acc}
}

//Child child
func (a Account) Child(i uint32) DerivedKey {
	key, err := a.accountKey.Child(i)
	checkError(err, fmt.Sprintf("error i %d", i))
	address, err := key.Address(a.masterKey.net)
	checkError(err, "error generating address from derived key")
	return DerivedKey{a, address.String(), key}
}

//Child child
func (d DerivedKey) Child(i uint32) DerivedKey {
	key, err := d.derived.Child(i)
	checkError(err, fmt.Sprintf("error i %d", i))
	address, err := key.Address(d.account.masterKey.net)
	checkError(err, "error generating address from derived key")
	return DerivedKey{d.account, address.String(), key}
}

//Account account
func (m MasterKey) Account(coin uint32, account uint32) Account {
	acc := m.ChildH(44).
		ChildH(coin).
		ChildH(account)
	return Account{m, acc.childKey, int(coin), int(account)}
}

//Derive derive change/index path
func (a Account) Derive(change uint32, index uint32) DerivedKey {
	dk := a.Child(change).
		Child(index)
	return dk
}

//AccountPK return the account extended public key
func (a Account) AccountPK() string {
	acctPk, err := a.accountKey.Neuter()
	checkError(err, "error neutering account")
	return acctPk.String()
}

//AccountSK return the account extended private key
func (a Account) AccountSK() string {
	return a.accountKey.String()
}

//Mnemonic return the mnemonic for this account
func (a Account) Mnemonic() string {
	return a.masterKey.mnemonic
}

//AccountJSON convert to JSON
func (a Account) AccountJSON() []byte {
	acctPk, err := a.accountKey.Neuter()
	checkError(err, "error neutering account")
	account := AccountJSON{
		Mnemonic:  a.masterKey.mnemonic,
		Seed:      hex.EncodeToString(a.masterKey.seed),
		RootSK:    a.masterKey.masterKey.String(),
		AccountPK: acctPk.String(),
		AccountSK: a.AccountSK(),
	}
	acctJSON, err := json.MarshalIndent(account, "", " ")
	checkError(err, "error marshalling json")
	return acctJSON
}

//PrintDerived derive accounts
func (a Account) PrintDerived(change uint32, count int) {
	for i := 0; i < count; i++ {
		var d = a.Derive(uint32(change), uint32(i))
		println(fmt.Sprintf("m/44'/%d'/%d'/%d/%d\t%s", a.coin, a.account, int(change), i, d.address))
	}
}

func checkError(err error, msg string) {
	if err != nil {
		println(fmt.Sprintf("%s, error = %s", msg, err))
		os.Exit(-1)
	}
}
