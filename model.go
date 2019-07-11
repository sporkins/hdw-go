package sporkins

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

//Account account
type Account struct {
	Mnemonic   string
	Seed       []byte
	MasterKey  *hdkeychain.ExtendedKey
	AccountKey *hdkeychain.ExtendedKey
	Coin       uint32
	Account    uint32
	Net        *chaincfg.Params
}

//DerivedKey derived key
type DerivedKey struct {
	Account Account
	Change  uint32
	Index   uint32
	Address string
	Derived *hdkeychain.ExtendedKey
}

//AccountJSON json of data
type AccountJSON struct {
	Mnemonic  string `json:"mnemonic"`
	Seed      string `json:"seed"`
	RootSK    string `json:"rootsk"`
	AccountSK string `json:"accountsk"`
	AccountPK string `json:"accountpk"`
}

//QRData QR Data
type QRData struct {
	Type                       string `json:"type"`
	DistributorPublicKeyBase64 string `json:"distributorPublicKeyBase64"`
	CipherTextBase64           string `json:"cipherTextBase64"`
	Key                        string `json:"key"`
	Name                       string `json:"name"`
}
