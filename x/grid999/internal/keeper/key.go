package keeper

import (
	"fmt"
)

var (
	KeyParams     = []byte{0x10}
	KeyNextDappID = []byte{0x11}

	DappKeyPrefix = "dapp"
	GridKeyPrefix = "grid"

	GridIDPrefix  = "gid"
	RandKeyPrefix = "rand"

	ParamsKeyPrefix = "params"
)

func KeyDapp(id uint) []byte {
	return []byte(fmt.Sprintf("%s:%d", DappKeyPrefix, id))
}
func KeyGrid(dappID uint, id uint64) []byte {
	return []byte(fmt.Sprintf("%s:%d:%d", GridKeyPrefix, dappID, id))
}
func KeyGridID(id uint) []byte {
	return []byte(fmt.Sprintf("%s:%d", GridIDPrefix, id))
}
func KeyRand(height int64) []byte {
	return []byte(fmt.Sprintf("%s:%d", RandKeyPrefix, height))
}
func KeyGridIDPrefix() []byte {
	return []byte(fmt.Sprintf("%s:", GridIDPrefix))
}
func KeyGridPrefix() []byte {
	return []byte(fmt.Sprintf("%s:", GridKeyPrefix))
}
