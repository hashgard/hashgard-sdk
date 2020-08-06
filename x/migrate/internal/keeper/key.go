package keeper

import (
	"fmt"
)

var (
	KeyErc20MigrateExchange = []byte{0x00}
	TxMigrateKeyIDPrefix    = "tsi"
	TxMigrateKeyPrefix      = "tsk"
	Erc20TxMigrateKeyPrefix = "erc20:"
)

func KeyERC20TxMigrateKey(erc20TxHash string) []byte {
	return []byte(fmt.Sprintf("%s%s", Erc20TxMigrateKeyPrefix, erc20TxHash))
}
func KeyERC20TxMigratePrefix() []byte {
	return []byte(fmt.Sprintf("%s", Erc20TxMigrateKeyPrefix))
}
func KeyNextTxMigrateID(migrateKey string) []byte {
	return []byte(fmt.Sprintf("%s:%s", TxMigrateKeyIDPrefix, migrateKey))
}
func KeyTxMigrate(migrateKey string, id uint64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%d", TxMigrateKeyPrefix, migrateKey, id))
}
func KeyTxMigratePrefix() []byte {
	return []byte(fmt.Sprintf("%s:", TxMigrateKeyPrefix))
}
