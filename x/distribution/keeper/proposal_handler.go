package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// HandleCommunityPoolSpendProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolSpendProposal(ctx sdk.Context, k Keeper, p types.CommunityPoolSpendProposal) sdk.Error {
	if k.blacklistedAddrs[p.Recipient.String()] {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is blacklisted from receiving external funds", p.Recipient))
	}

	err := k.DistributeFromFeePool(ctx, p.Amount, p.Recipient)
	if err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("transferred %s from the community pool to recipient %s", p.Amount, p.Recipient))
	return nil
}

// HashGard
// HandleStakeIssueLockedSpendProposal is a handler for executing a passed community spend proposal
func HandleStakeIssueLockedSpendProposal(ctx sdk.Context, k Keeper, p types.StakeIssueLockedSpendProposal) sdk.Error {
	if k.blacklistedAddrs[p.Recipient.String()] {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is blacklisted from receiving external funds", p.Recipient))
	}

	if err := k.stakingKeeper.StakeIssueLockedSpend(ctx, p.Denom, p.Recipient, p.Amount); err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("transferred %s from the stake issue locked to recipient %s", p.Amount, p.Recipient))
	return nil
}
