package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
)

func (k Keeper) GetCommunityPoolAddr() sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf(distribution.ModuleName))))
}
func (k Keeper) GetDepositAddr(dappID uint) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("gridDepositAddr_%d", dappID))))
}
func (k Keeper) GetLuckyPoolAddr(dappID uint) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("gridLuckyPoolAddr_%d", dappID))))
}
func (k Keeper) GetFeeAddr(dappID uint) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("gridFeeAddr%d_", dappID))))
}
func (k Keeper) CreateGrid(ctx sdk.Context, dappID uint, gridItem *types.GridItem) (grid *types.Grid, id uint64, err sdk.Error) {
	dapp := k.GetDapp(ctx, dappID)
	if dapp == nil {
		return grid, id, sdk.ErrInternal("dapp is not found")
	}
	if dapp.OnlyOwnerCanCreateGrid && !dapp.Owner.Equals(gridItem.Owner) {
		return grid, id, sdk.ErrInternal("dapp is only owner can create grid")
	}
	if dapp.EndHeight > 0 && dapp.EndHeight < ctx.BlockHeight() {
		return grid, id, sdk.ErrInternal("dapp is disabled")
	}

	if !dapp.OwnerMinDeposit.IsZero() && gridItem.OwnerDeposit.IsLT(dapp.OwnerMinDeposit) {
		return grid, id, sdk.ErrInternal("owner min deposit is " + dapp.OwnerMinDeposit.String())
	}
	params := k.GetParam(ctx)
	if !params.CreateGridFee.IsZero() {
		if err := k.SendGrid999FeeCoins(ctx, gridItem.Owner, sdk.NewCoins(params.CreateGridFee)); err != nil {
			return grid, id, err
		}
	}
	store := ctx.KVStore(k.storeKey)
	id = k.getLastGridID(store, dappID)

	grid, _ = k.GetGrid(ctx, dappID, id)
	if grid != nil {
		if grid.Height+dapp.MaxBlocksGridCreate < ctx.BlockHeight() {
			id += 1
			k.setNewGridID(store, dappID, id)
			grid = types.DefaultGrid(ctx.BlockHeight(), dappID)
		}
	} else {
		grid = types.DefaultGrid(ctx.BlockHeight(), dappID)
		k.setNewGridID(store, dappID, id)
	}
	if err = k.checkMaxDepositAmount(dapp.MaxDepositAmount, grid.TotalDeposit, gridItem.OwnerDeposit); err != nil {
		return
	}
	grid.TotalDeposit = grid.TotalDeposit.Add(sdk.NewCoins(gridItem.OwnerDeposit))

	if gridItem.ZeroValued {
		gridItem.TotalNumber = -9223372036854775808
	} else {
		gridItem.OwnerNumber = k.GetRand(ctx, gridItem.Owner, dappID, dapp.RandNumberNegativeCriticalValue, types.NumberDigit)
		gridItem.TotalNumber = gridItem.OwnerNumber + gridItem.OwnerDeposit.Amount.Quo(sdk.TokensFromConsensusPower(1)).Int64()
	}

	maxGrid := dapp.Ranks.GetMaxRank()
	if len(grid.Items) == maxGrid {
		if grid.Items[maxGrid-1].OwnerDeposit.IsLT(gridItem.OwnerDeposit) {
			if err := k.SendCoins(ctx, false, k.GetDepositAddr(dappID), grid.Items[maxGrid-1].Owner, sdk.NewCoins(grid.Items[maxGrid-1].OwnerDeposit)); err != nil {
				return grid, id, err
			}
			grid.Items[maxGrid-1] = *gridItem
		} else {
			return grid, id, sdk.ErrInternal("can not create grid now")
		}
	} else {
		grid.Items = append(grid.Items, *gridItem)
	}

	if err := k.SendCoins(ctx, false, gridItem.Owner, k.GetDepositAddr(dappID), sdk.NewCoins(gridItem.OwnerDeposit)); err != nil {
		return grid, id, err
	}

	grid.SortByNumber(dapp.RankType)
	for i, _ := range grid.Items {
		grid.Items[i].Index = uint(i)
	}
	grid.ID = id
	k.SetGrid(ctx, grid)

	return
}
func (k Keeper) checkMaxDepositAmount(maxDepositAmount, totalDeposit sdk.Coins, deposit sdk.Coin) sdk.Error {
	if maxDepositAmount.IsZero() {
		return nil
	}
	if totalDeposit.Add(sdk.NewCoins(deposit)).IsAllLTE(maxDepositAmount) {
		return nil
	}
	return sdk.ErrInternal("can not deposit to the grid any more,max deposit amount is " + maxDepositAmount.String())
}

func (k Keeper) DepositGrid(ctx sdk.Context, dappID, indexSeq uint, id uint64, sender sdk.AccAddress, deposit sdk.Coin) (luckyCoin, feeCoin sdk.Coin, err sdk.Error) {
	dapp := k.GetDapp(ctx, dappID)
	if dapp == nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("dapp is not found")
	}
	if deposit.IsLT(dapp.MemberMinDeposit) {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("member min deposit is " + dapp.MemberMinDeposit.String())
	}
	grid, err := k.GetGrid(ctx, dappID, id)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}
	if !dapp.GridCreateCanDeposit && grid.Height+dapp.MaxBlocksGridCreate > ctx.BlockHeight() {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("can not deposit to the grid now")
	}
	if err = k.checkMaxDepositAmount(dapp.MaxDepositAmount, grid.TotalDeposit, deposit); err != nil {
		return
	}
	index := -1
	for i, v := range grid.Items {
		if v.Index == indexSeq {
			index = i
			break
		}
	}
	if index < 0 {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("grid with no designated owner")
	}
	ranks := dapp.Ranks.GetRank(len(grid.Items))
	if len(ranks) == 0 {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("can not deposit to the grid now")
	}
	endBlocks := grid.Height + dapp.MaxBlocksGridCreate + dapp.MaxBlocksGridDeposit
	if endBlocks < ctx.BlockHeight() {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("can not deposit to the grid any more,end blocks is " + fmt.Sprintf("%d", endBlocks))
	}
	if len(grid.Items[index].Deposits) >= dapp.PerGridMaxDeposits {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("can not deposit to the grid any more,max deposits is " + fmt.Sprintf("%d", dapp.PerGridMaxDeposits))
	}
	if grid.Items[index].Locked {
		return sdk.Coin{}, sdk.Coin{}, sdk.ErrInternal("can not deposit to the locked grid")
	}

	feeCoin = sdk.NewCoin(deposit.Denom, sdk.ZeroInt())
	params := k.GetParam(ctx)
	if !dapp.Voted && !params.DepositFee.IsZero() {
		feeCoin = params.DepositFee
		if err = k.SendGrid999FeeCoins(ctx, sender, sdk.NewCoins(feeCoin)); err != nil {
			return sdk.Coin{}, sdk.Coin{}, err
		}
	}

	luckyCoin = k.depositToLuckyPool(ctx, dappID, deposit, dapp, sender)
	if !luckyCoin.IsZero() {
		if luckyCoin.IsLT(deposit) {
			if k.SendCoins(ctx, false, sender, k.GetLuckyPoolAddr(dappID), sdk.NewCoins(luckyCoin)); err != nil {
				return sdk.Coin{}, sdk.Coin{}, err
			}
			deposit = deposit.Sub(luckyCoin)
		} else {
			luckyCoin = sdk.NewCoin(deposit.Denom, sdk.ZeroInt())
		}
	}

	if err = k.SendCoins(ctx, false, sender, k.GetDepositAddr(dappID), sdk.NewCoins(deposit)); err != nil {
		return
	}

	if grid.Items[index].MemberDeposit.IsZero() {
		grid.Items[index].MemberDeposit = deposit
	} else {
		grid.Items[index].MemberDeposit = grid.Items[index].MemberDeposit.Add(deposit)
	}

	grid.TotalDeposit = grid.TotalDeposit.Add(sdk.NewCoins(deposit))
	var number int64 = 0
	if !grid.Items[index].ZeroValued {
		number = k.GetRand(ctx, sender, dappID, dapp.RandNumberNegativeCriticalValue, types.NumberDigit)
		grid.Items[index].TotalNumber += number
		grid.Items[index].TotalNumber += deposit.Amount.Quo(sdk.TokensFromConsensusPower(1)).Int64()
	}
	grid.Items[index].Deposits = append(grid.Items[index].Deposits, fmt.Sprintf("%s_%d_%s_%s_%s", sender.String(), number, deposit.String(), feeCoin.String(), luckyCoin.String()))

	grid.SortByNumber(dapp.RankType)

	store := ctx.KVStore(k.storeKey)
	store.Set(KeyGrid(dappID, id), k.cdc.MustMarshalBinaryLengthPrefixed(grid))
	return
}

func (k Keeper) depositToLuckyPool(ctx sdk.Context, dappID uint, deposit sdk.Coin, dapp *types.Dapp, sender sdk.AccAddress) (luckyCoin sdk.Coin) {

	if !dapp.LuckyPoolRatio.IsZero() {
		return sdk.NewCoin(deposit.Denom, dapp.LuckyPoolRatio.MulInt(deposit.Amount).TruncateInt())
	}

	if dapp.DepositToLuckyPoolDigit > 0 {
		return sdk.NewCoin(deposit.Denom, sdk.TokensFromConsensusPower(k.GetLuckyRand(ctx, sender, dappID, dapp.DepositToLuckyPoolDigit)))
	}

	return sdk.NewCoin(deposit.Denom, sdk.ZeroInt())
}
func (k Keeper) WithdrawGrid(ctx sdk.Context, dappID uint, id uint64, sender sdk.AccAddress) (withdrawDeposit, withdrawRewards, withdrawLucky, fee sdk.Coins, err sdk.Error) {
	grid, err := k.GetGrid(ctx, dappID, id)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, err
	}
	dapp := k.GetDapp(ctx, dappID)
	if dapp == nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.ErrInternal("dapp is not found")
	}
	rewardsCanWithdraw := false
	depositCanWithdraw := false

	if grid.Height+dapp.MaxBlocksGridCreate+dapp.MaxBlocksGridDeposit+dapp.MaxBlocksGridRewardsWithdraw <= ctx.BlockHeight() {
		rewardsCanWithdraw = true
	}
	if grid.Height+dapp.MaxBlocksGridCreate+dapp.MaxBlocksGridDeposit+dapp.MaxBlocksGridDepositWithdraw <= ctx.BlockHeight() {
		depositCanWithdraw = true
	}

	if !rewardsCanWithdraw && !depositCanWithdraw {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.ErrInternal("can not withdraw from the grid now")
	}
	ranks := dapp.Ranks.GetRank(len(grid.Items))

	if dapp.Voted || len(ranks) == 0 || (!dapp.MinDepositAmount.IsZero() && dapp.MinDepositAmount.IsAnyGT(grid.TotalDeposit)) {
		withdrawDeposit, err = k.backDeposit(ctx, grid, dapp, sender)
		if err != nil {
			return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, err
		}
		withdrawRewards = k.getDefaultCoin(dapp)
		withdrawLucky = withdrawRewards
		fee = withdrawRewards
		return
	}

	var totalLastDeposit = k.getDefaultCoin(dapp)
	for _, v := range dapp.Ranks.GetLoser(ranks) {
		totalLastDeposit = k.Add(totalLastDeposit, grid.Items[v-1].OwnerDeposit)
		totalLastDeposit = k.Add(totalLastDeposit, grid.Items[v-1].MemberDeposit)
	}
	rewardsFee := k.getDefaultCoin(dapp)
	withdrawDeposit, withdrawRewards, rewardsFee, err = k.withdrawRewards(ctx, grid, dapp, rewardsCanWithdraw, sender, totalLastDeposit, dappID, ranks)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, err
	}

	deposit := k.getDefaultCoin(dapp)
	deposit, fee, err = k.deductFee(ctx, grid, dapp, depositCanWithdraw, sender, ranks)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, err
	}

	if depositCanWithdraw {
		withdrawDeposit = withdrawDeposit.Add(deposit)
		fee = rewardsFee.Add(fee)
		if !withdrawDeposit.IsZero() {
			if !strings.Contains(strings.Join(grid.Deposits, ","), sender.String()) {
				if err = k.SendCoins(ctx, true, k.GetDepositAddr(grid.DappID), sender, withdrawDeposit); err != nil {
					return
				}
				grid.Deposits = append(grid.Deposits, fmt.Sprintf("%s_%s", sender.String(), withdrawDeposit.String()))
			}
		}
	} else {
		fee = k.getDefaultCoin(dapp)
		withdrawDeposit = fee
	}

	withdrawLucky, err = k.withdrawLucky(ctx, grid, dapp, sender, dappID, ranks)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, sdk.Coins{}, err
	}

	k.SetGrid(ctx, grid)
	return
}
func (k Keeper) getDefaultCoin(dapp *types.Dapp) sdk.Coins {
	if dapp.OwnerMinDeposit.Denom == dapp.MemberMinDeposit.Denom {
		return sdk.NewCoins(sdk.NewCoin(dapp.OwnerMinDeposit.Denom, sdk.ZeroInt()))
	}
	return sdk.NewCoins(
		sdk.NewCoin(dapp.OwnerMinDeposit.Denom, sdk.ZeroInt()),
		sdk.NewCoin(dapp.MemberMinDeposit.Denom, sdk.ZeroInt()))
}
func (k Keeper) MulRate(dapp *types.Dapp, coin sdk.Coins, rate sdk.Dec) sdk.Coins {
	if dapp.OwnerMinDeposit.Denom == dapp.MemberMinDeposit.Denom {
		return sdk.NewCoins(
			sdk.NewCoin(dapp.OwnerMinDeposit.Denom, rate.MulInt(coin.AmountOf(dapp.OwnerMinDeposit.Denom)).TruncateInt()))
	}
	ownerMinDeposit := sdk.NewCoin(dapp.OwnerMinDeposit.Denom, rate.MulInt(coin.AmountOf(dapp.OwnerMinDeposit.Denom)).TruncateInt())
	memberMinDeposit := sdk.NewCoin(dapp.MemberMinDeposit.Denom, rate.MulInt(coin.AmountOf(dapp.MemberMinDeposit.Denom)).TruncateInt())
	return sdk.NewCoins(ownerMinDeposit, memberMinDeposit)
}
func (k Keeper) Add(coins sdk.Coins, coinB sdk.Coin) sdk.Coins {
	return coins.Add(sdk.NewCoins(coinB))
}
func (k Keeper) deductWinnerFee(ctx sdk.Context, dapp *types.Dapp, withdrawDeposit, withdrawRewards sdk.Coins, sender sdk.AccAddress, dappID uint) (feeCoin, retWithdrawDeposit, retWithdrawRewards sdk.Coins, err sdk.Error) {

	feeCoin = k.getDefaultCoin(dapp)
	retWithdrawRewards = withdrawRewards
	retWithdrawDeposit = withdrawDeposit

	params := k.GetParam(ctx)
	if params.WithdrawRewardsFee.IsZero() {
		return
	}
	feeCoin = k.Add(feeCoin, params.WithdrawRewardsFee)
	if withdrawRewards.IsAllGTE(feeCoin) {
		if err = k.SendGrid999FeeCoins(ctx, k.GetDepositAddr(dappID), feeCoin); err != nil {
			return
		}
		retWithdrawRewards = retWithdrawRewards.Sub(feeCoin)
		return
	}
	if withdrawDeposit.IsAllGTE(feeCoin) {
		if err = k.SendGrid999FeeCoins(ctx, k.GetDepositAddr(dappID), feeCoin); err != nil {
			return
		}
		retWithdrawDeposit = retWithdrawDeposit.Sub(feeCoin)
		return
	}
	if err = k.SendGrid999FeeCoins(ctx, sender, feeCoin); err != nil {
		return
	}
	return
}
func (k Keeper) withdrawRewards(ctx sdk.Context, grid *types.Grid, dapp *types.Dapp, rewardsCanWithdraw bool, sender sdk.AccAddress, totalLastDeposit sdk.Coins, dappID uint, ranks [][]int) (withdrawDeposit, withdrawRewards, fee sdk.Coins, err sdk.Error) {
	withdrawDeposit = k.getDefaultCoin(dapp)
	withdrawRewards = k.getDefaultCoin(dapp)
	fee = k.getDefaultCoin(dapp)

	if strings.Contains(strings.Join(grid.Rewards, ","), sender.String()) {
		return
		//return withdrawDeposit, withdrawRewards, fee, sdk.ErrInternal("you have withdraw")
	}
	for _, v := range dapp.Ranks.GetWinner(ranks) {
		if sender.Equals(grid.Items[v-1].Owner) {
			var rewards sdk.Coins
			rewards, err = k.ownerWithdraw(dapp, totalLastDeposit, len(ranks[0]), v-1)
			if err != nil {
				return
			}
			withdrawDeposit = k.Add(withdrawDeposit, grid.Items[v-1].OwnerDeposit)
			withdrawRewards = withdrawRewards.Add(rewards)
		}
		for _, deposits := range grid.Items[v-1].Deposits {
			if strings.HasPrefix(deposits, sender.String()) {
				rewards, deposit, err := k.memberWithdraw(dapp, deposits, totalLastDeposit, grid.Items[v-1].MemberDeposit.Amount, len(ranks[0]), v-1)
				if err != nil {
					return withdrawDeposit, withdrawRewards, fee, sdk.ErrInternal(err.Error())
				}
				withdrawDeposit = k.Add(withdrawDeposit, deposit)
				withdrawRewards = withdrawRewards.Add(rewards)
			}
		}
	}
	if !rewardsCanWithdraw || withdrawRewards.IsZero() {
		//if err = k.SendCoins(ctx, true, k.GetDepositAddr(dappID), sender, withdrawDeposit); err != nil {
		//	return
		//}
		return
	}
	feeCoin := k.getDefaultCoin(dapp)
	feeCoin, withdrawDeposit, withdrawRewards, err = k.deductWinnerFee(ctx, dapp, withdrawDeposit, withdrawRewards, sender, dappID)
	if err != nil {
		return
	}
	fee = fee.Add(feeCoin)
	if !dapp.FeeRatio.IsZero() && len(ranks) == 2 {
		feeCoin = k.MulRate(dapp, withdrawRewards, dapp.FeeRatio)
		if err = k.SendCoins(ctx, true, k.GetDepositAddr(dappID), k.GetFeeAddr(grid.DappID), feeCoin); err != nil {
			return
		}
		withdrawRewards = withdrawRewards.Sub(feeCoin)
		fee = fee.Add(feeCoin)
		grid.Fees = append(grid.Fees, fmt.Sprintf("%s_%s", sender.String(), fee.String()))
	}

	if err = k.SendCoins(ctx, true, k.GetDepositAddr(dappID), sender, withdrawRewards); err != nil {
		return
	}

	grid.Rewards = append(grid.Rewards, fmt.Sprintf("%s_%s", sender.String(), withdrawRewards.String()))
	return
}
func (k Keeper) withdrawLucky(ctx sdk.Context, grid *types.Grid, dapp *types.Dapp, sender sdk.AccAddress, dappID uint, ranks [][]int) (withdrawLucky sdk.Coins, err sdk.Error) {
	withdrawLucky = k.getDefaultCoin(dapp)
	if dapp.LuckyPoolRewardsDigit == 0 {
		return
	}
	if strings.Contains(strings.Join(grid.Lucky, ","), sender.String()) {
		return
		//return withdrawLucky, sdk.ErrInternal("you have withdraw")
	}
	for _, v := range dapp.Ranks.GetLoser(ranks) {
		if sender.Equals(grid.Items[v-1].Owner) {
			lucky := sdk.NewCoin(grid.Items[v-1].OwnerDeposit.Denom, sdk.TokensFromConsensusPower(k.GetLuckyRand(ctx, sender, dappID, dapp.LuckyPoolRewardsDigit)))
			withdrawLucky = k.Add(withdrawLucky, lucky)
		}
		for _, deposits := range grid.Items[v-1].Deposits {
			if strings.HasPrefix(deposits, sender.String()) {
				values := strings.Split(deposits, "_")
				deposit, err := sdk.ParseCoin(values[2])
				if err != nil {
					return withdrawLucky, sdk.ErrInternal(err.Error())
				}
				lucky := sdk.NewCoin(deposit.Denom, sdk.TokensFromConsensusPower(k.GetLuckyRand(ctx, sender, dappID, dapp.LuckyPoolRewardsDigit)))
				withdrawLucky = k.Add(withdrawLucky, lucky)
			}
		}
	}
	if withdrawLucky.IsZero() {
		return
	}
	coins := k.bankKeeper.GetCoins(ctx, k.GetLuckyPoolAddr(dappID))
	lucky := k.getDefaultCoin(dapp)
	for _, coin := range coins {
		amount := withdrawLucky.AmountOf(coin.Denom)
		if !amount.IsZero() && coin.Amount.GTE(amount) {
			lucky = k.Add(lucky, sdk.NewCoin(coin.Denom, amount))
		}
	}

	withdrawLucky = lucky
	if withdrawLucky.IsZero() {
		return
	}

	if err := k.SendCoins(ctx, true, k.GetLuckyPoolAddr(dappID), sender, withdrawLucky); err == nil {
		grid.Lucky = append(grid.Lucky, fmt.Sprintf("%s_%s", sender.String(), withdrawLucky.String()))
		return withdrawLucky, nil
	}

	return k.getDefaultCoin(dapp), nil
}
func (k Keeper) deductFee(ctx sdk.Context, grid *types.Grid, dapp *types.Dapp, depositCanWithdraw bool, sender sdk.AccAddress, ranks [][]int) (totalDeposit, totalFee sdk.Coins, err sdk.Error) {
	totalDeposit = k.getDefaultCoin(dapp)
	totalFee = k.getDefaultCoin(dapp)

	if len(ranks) == 2 {
		return
	}

	if strings.Contains(strings.Join(grid.Fees, ","), sender.String()) {
		//return totalDeposit, totalFee, sdk.ErrInternal("you have withdraw")
		return
	}
	for _, v := range dapp.Ranks.GetFeePayer(ranks) {
		if sender.Equals(grid.Items[v-1].Owner) {
			totalDeposit = k.Add(totalDeposit, grid.Items[v-1].OwnerDeposit)
			totalFee = k.Add(totalFee, sdk.NewCoin(grid.Items[v-1].OwnerDeposit.Denom, dapp.FeeRatio.MulInt(grid.Items[v-1].OwnerDeposit.Amount).TruncateInt()))
		}
		for _, deposits := range grid.Items[v-1].Deposits {
			if strings.HasPrefix(deposits, sender.String()) {
				values := strings.Split(deposits, "_")
				deposit, err := sdk.ParseCoin(values[2])
				if err != nil {
					return totalDeposit, totalFee, sdk.ErrInternal(err.Error())
				}
				totalDeposit = k.Add(totalDeposit, deposit)
				totalFee = k.Add(totalFee, sdk.NewCoin(deposit.Denom, dapp.FeeRatio.MulInt(deposit.Amount).TruncateInt()))
			}
		}
	}
	totalDeposit = totalDeposit.Sub(totalFee)
	if !depositCanWithdraw {
		return
	}
	if totalFee.IsZero() {
		return
	}
	if err = k.SendCoins(ctx, true, k.GetDepositAddr(grid.DappID), k.GetFeeAddr(grid.DappID), totalFee); err != nil {
		return
	}
	grid.Fees = append(grid.Fees, fmt.Sprintf("%s_%s", sender.String(), totalFee.String()))
	return
}

func (k Keeper) backDeposit(ctx sdk.Context, grid *types.Grid, dapp *types.Dapp, sender sdk.AccAddress) (totalDeposit sdk.Coins, err sdk.Error) {
	totalDeposit = k.getDefaultCoin(dapp)
	if strings.Contains(strings.Join(grid.Backs, ","), sender.String()) {
		return totalDeposit, sdk.ErrInternal("you have withdraw")
	}
	for _, item := range grid.Items {
		if sender.Equals(item.Owner) {
			totalDeposit = k.Add(totalDeposit, item.OwnerDeposit)
		}
		for _, deposits := range item.Deposits {
			values := strings.Split(deposits, "_")
			coin, err1 := sdk.ParseCoin(values[2])
			if err1 != nil {
				return totalDeposit, sdk.ErrInternal(err1.Error())
			}
			if sender.String() == values[0] {
				totalDeposit = k.Add(totalDeposit, coin)
			}
		}
	}
	if totalDeposit.IsZero() {
		return totalDeposit, sdk.ErrInternal("you have not deposit")
	}
	if err = k.SendCoins(ctx, true, k.GetDepositAddr(grid.DappID), sender, totalDeposit); err != nil {
		return
	}
	grid.Backs = append(grid.Backs, fmt.Sprintf("%s_%s", sender.String(), totalDeposit.String()))
	k.SetGrid(ctx, grid)
	return
}
func (k Keeper) ownerWithdraw(dapp *types.Dapp, totalLastDeposit sdk.Coins, seqs, seq int) (rewardsCoin sdk.Coins, err sdk.Error) {
	rewardsCoin = k.getDefaultCoin(dapp)
	if winners, ok := dapp.WinnerRewards.GetWinnerRewards(seqs); ok {
		rewardsCoin = k.MulRate(dapp, totalLastDeposit, winners[seq])
		if rewardsCoin.IsZero() {
			return
		}
		rewardsCoin = k.MulRate(dapp, rewardsCoin, dapp.OwnerRewardsRatio)
	}
	return
}
func (k Keeper) memberWithdraw(dapp *types.Dapp, deposits string, totalLastDeposit sdk.Coins, memberDeposit sdk.Int, seqs, seq int) (rewardsCoin sdk.Coins, deposit sdk.Coin, err error) {
	rewardsCoin = k.getDefaultCoin(dapp)
	values := strings.Split(deposits, "_")
	deposit, err = sdk.ParseCoin(values[2])
	if err != nil {
		return
	}
	if deposit.IsZero() || memberDeposit.IsZero() {
		return
	}
	if winners, ok := dapp.WinnerRewards.GetWinnerRewards(seqs); ok {
		rewardsCoin = k.MulRate(dapp, totalLastDeposit, winners[seq])
		if rewardsCoin.IsZero() {
			return
		}
		rewardsCoin = k.MulRate(dapp, rewardsCoin, sdk.OneDec().Sub(dapp.OwnerRewardsRatio))
		rewardsCoin = k.MulRate(dapp, rewardsCoin, sdk.NewDecFromInt(deposit.Amount).QuoInt(memberDeposit))
	}
	return
}

func (k Keeper) GetGrids(ctx sdk.Context, dappID uint, numLatest uint64, limit int) []*types.Grid {
	store := ctx.KVStore(k.storeKey)
	list := make([]*types.Grid, 0, limit)
	maxId := k.getLastGridID(store, dappID)
	if numLatest == 0 {
		numLatest = maxId
	}
	for txID := numLatest; txID > 0; txID-- {
		grid, err := k.GetGrid(ctx, dappID, txID)
		if err != nil {
			continue
		}
		list = append(list, grid)
		if len(list) >= limit {
			break
		}
	}
	return list
}
func (k Keeper) ExportGrids(ctx sdk.Context) []types.Grid {
	store := ctx.KVStore(k.storeKey)
	list := make([]types.Grid, 0)
	iterator := sdk.KVStorePrefixIterator(store, KeyGridPrefix())
	defer iterator.Close()
	var bz []byte
	for ; iterator.Valid(); iterator.Next() {
		bz = iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var grid types.Grid
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &grid)
		list = append(list, grid)
	}
	return list
}
func (k Keeper) ExportGridIDs(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)
	list := make([]string, 0)
	iterator := sdk.KVStorePrefixIterator(store, KeyGridIDPrefix())
	defer iterator.Close()
	var bz []byte
	for ; iterator.Valid(); iterator.Next() {
		bz = iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var id uint64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
		list = append(list, fmt.Sprintf("%s_%d", string(iterator.Key()), id))
	}
	return list
}
func (k Keeper) ImportGridIDs(ctx sdk.Context, kv []string) {
	store := ctx.KVStore(k.storeKey)
	for _, v := range kv {
		values := strings.Split(v, "_")
		store.Set([]byte(values[0]), k.cdc.MustMarshalBinaryLengthPrefixed(values[1]))
	}
}
func (k Keeper) GetGrid(ctx sdk.Context, dappID uint, id uint64) (grid *types.Grid, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyGrid(dappID, id))
	if len(bz) == 0 {
		err = sdk.ErrInternal("grid is not found")
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &grid)
	return
}
func (k Keeper) SetGrid(ctx sdk.Context, grid *types.Grid) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyGrid(grid.DappID, grid.ID), k.cdc.MustMarshalBinaryLengthPrefixed(grid))
	return
}

func (k Keeper) SetNewGridID(ctx sdk.Context, dappID uint, id uint64) {
	store := ctx.KVStore(k.storeKey)
	k.setNewGridID(store, dappID, id)
}
func (k Keeper) GetLastGridID(ctx sdk.Context, dappID uint) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	return k.getLastGridID(store, dappID)
}
func (k Keeper) getLastGridID(store sdk.KVStore, dappID uint) (id uint64) {
	bz := store.Get(KeyGridID(dappID))
	id = 1
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	}
	return
}
func (k Keeper) setNewGridID(store sdk.KVStore, dappID uint, id uint64) {
	store.Set(KeyGridID(dappID), k.cdc.MustMarshalBinaryLengthPrefixed(id))
	return
}
