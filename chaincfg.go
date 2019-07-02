package sporkins

import (
	"github.com/btcsuite/btcd/chaincfg"
)

//RvnMainNetParams rvn mainnet params
var RvnMainNetParams = chaincfg.Params{
	PubKeyHashAddrID: 0x3c,                            // R
	HDPrivateKeyID:   [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:    [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
	HDCoinType:       175,
}

//NetworkParams get network params
func NetworkParams(coin int) *chaincfg.Params {
	switch coin {
	case 0:
		return &chaincfg.MainNetParams
	case 1:
		return &chaincfg.TestNet3Params
	case 175:
		return &RvnMainNetParams
	}
	return nil
}
