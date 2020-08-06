package keeper

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"strconv"
	"strings"
)

var (
	numberReg, _ = regexp.Compile("[^0-9]+")
)

func (k Keeper) GetRandRound(ctx sdk.Context) (count int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyRand(ctx.BlockHeight()))
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	}
	store.Set(KeyRand(ctx.BlockHeight()), k.cdc.MustMarshalBinaryLengthPrefixed(count+1))
	return
}

func (k Keeper) GetRand(ctx sdk.Context, sender sdk.AccAddress, dappId, randNumberNegativeCriticalValue uint, digit int) int64 {
	randKey := k.getRandNumber(ctx, sender.String(), dappId, k.GetRandRound(ctx))
	v, _ := strconv.ParseInt(randKey[0:digit], 10, 64)
	v1, _ := strconv.Atoi(randKey[0:2])
	if uint(v1) <= randNumberNegativeCriticalValue {
		return -v
	}
	return v
}
func (k Keeper) GetLuckyRand(ctx sdk.Context, sender sdk.AccAddress, dappId uint, digit int) int64 {
	randKey := k.getRandNumber(ctx, sender.String(), dappId, k.GetRandRound(ctx))
	if digit > len(randKey) {
		digit = len(randKey)
	}
	v, _ := strconv.ParseInt(randKey[0:digit], 10, 64)
	return v
}
func (k Keeper) getRandNumber(ctx sdk.Context, sender string, dappId uint, randRound int) string {
	randValue := string(numberReg.ReplaceAll([]byte(k.getRandStr(ctx, sender, dappId, randRound)), []byte("")))
	return strings.TrimLeft(randValue, "0")
}
func (k Keeper) getRandStr(ctx sdk.Context, sender string, dappId uint, randRound int) string {
	randKey := fmt.Sprintf("%s%s%s%d%d", hex.EncodeToString(ctx.TxBytes()), hex.EncodeToString(ctx.BlockHeader().LastCommitHash), strings.ReplaceAll(sender, sdk.Bech32MainPrefix, ""), dappId, randRound)

	sha := sha1.New()
	sha.Write([]byte(randKey))
	shaValues1 := sha.Sum(nil)

	sha.Reset()
	sha.Write(shaValues1)
	shaValues2 := sha.Sum(nil)

	randValue := fmt.Sprintf("%s%s%s", hex.EncodeToString(shaValues1), hex.EncodeToString(shaValues2), randKey)
	return randValue
}
