package sporkins

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

//Mnemonic bip39 mnemonic
type Mnemonic struct {
	mnemonic string
	entropy  []byte
}

//Seed seed from mnemonic
type Seed struct {
	mnemonic string
	entropy  []byte
	seed     []byte
}

//MasterKey master key
type MasterKey struct {
	mnemonic  string
	entropy   []byte
	seed      []byte
	masterKey *hdkeychain.ExtendedKey
	net       *chaincfg.Params
}

type ChildKey struct {
	childKey *hdkeychain.ExtendedKey
}

type Account struct {
	masterKey  MasterKey
	accountKey *hdkeychain.ExtendedKey
	coin       int
	account    int
}

type DerivedKey struct {
	account Account
	address string
	derived *hdkeychain.ExtendedKey
}

//AccountJSON json of data
type AccountJSON struct {
	Mnemonic  string `json:"mnemonic"`
	Seed      string `json:"seed"`
	RootSK    string `json:"rootsk"`
	AccountSK string `json:"accountsk"`
	AccountPK string `json:"accountpk"`
}
