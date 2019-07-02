package sporkins

import (
	"github.com/btcsuite/btcd/chaincfg"
)

var Networks = map[string]bool{
	"mainnet": true,
	"testnet": true,
}

//RvnMainNetParams rvn mainnet params
var RvnMainNetParams = chaincfg.Params{
	PubKeyHashAddrID: 0x3c,                            // R
	HDPrivateKeyID:   [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:    [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
	HDCoinType:       175,
}

//RvnTestNetParams rvn mainnet params
var RvnTestNetParams = chaincfg.Params{
	PubKeyHashAddrID: 0x6f,                            // R
	HDPrivateKeyID:   [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:    [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	HDCoinType:       175,
}

//NetworkParams get network params
func NetworkParams(coin int, network string) *chaincfg.Params {
	switch coin {
	case 1, 2:
		switch network {
		case "mainnet":
			return &chaincfg.MainNetParams
		case "testnet":
			return &chaincfg.TestNet3Params
		}
	case 175:
		switch network {
		case "mainnet":
			return &RvnMainNetParams
		case "testnet":
			return &RvnTestNetParams
		}
	}
	return nil
}
