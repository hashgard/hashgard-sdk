package grid999

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"strings"
)

// NewHandler returns a handler for "scene" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgGridParams:
			return handleMsgGridParams(ctx, k, msg)
		case types.MsgDisableDapp:
			return handleMsgDisableDapp(ctx, k, msg)
		case types.MsgDappGenerate:
			return handleMsgDappGenerate(ctx, k, msg)
		case types.MsgDappCreateGrid:
			return handleMsgDappCreateGrid(ctx, k, msg)
		case types.MsgDappDeposit:
			return handleMsgDappDeposit(ctx, k, msg)
		case types.MsgDappWithdraw:
			return handleMsgDappWithdraw(ctx, k, msg)
		case types.MsgDappWithdrawFees:
			return handleMsgDappWithdrawFees(ctx, k, msg)
		case types.MsgDappWithdrawLucky:
			return handleMsgDappWithdrawLucky(ctx, k, msg)
		case types.MsgGridWithdrawFees:
			return handleMsgGridWithdrawFees(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized grid999 message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
func handleMsgGridParams(ctx sdk.Context, k Keeper, msg types.MsgGridParams) sdk.Result {
	params := k.GetParam(ctx)
	if !strings.Contains(params.Owners, msg.Sender.String()) {
		return sdk.ErrInternal("Permission denied").Result()
	}
	msg.Params.Owners = params.Owners
	k.SetParam(ctx, msg.Params)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappCreateGrid(ctx sdk.Context, k Keeper, msg types.MsgDappCreateGrid) sdk.Result {
	gridItem := &types.GridItem{
		Owner:        msg.Sender,
		OwnerDeposit: msg.Deposit,
		GridType:     msg.GridType,
		ZeroValued:   msg.ZeroValued,
		Prepaid:      msg.Prepaid,
	}
	_, id, err := k.CreateGrid(ctx, msg.DappID, gridItem)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			types.EventTypeGridCreate,
			sdk.NewAttribute(types.AttributeDappID, fmt.Sprintf("%d", msg.DappID)),
			sdk.NewAttribute(types.AttributeGridId, fmt.Sprintf("%d", id)),
			sdk.NewAttribute(types.AttributeOwnerNumber, fmt.Sprintf("%d", gridItem.OwnerNumber)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappDeposit(ctx sdk.Context, k Keeper, msg types.MsgDappDeposit) sdk.Result {
	lucky, fee, err := k.DepositGrid(ctx, msg.DappID, msg.Index, msg.GridId, msg.Sender, msg.Deposit)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			types.EventTypeGridDeposit,
			sdk.NewAttribute(types.AttributeDappID, fmt.Sprintf("%d", msg.DappID)),
			sdk.NewAttribute(types.AttributeGridId, fmt.Sprintf("%d", msg.GridId)),
			sdk.NewAttribute(types.AttributeGridLuckyDeposit, lucky.String()),
			sdk.NewAttribute(types.AttributeWithdrawFee, fee.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappWithdraw(ctx sdk.Context, k Keeper, msg types.MsgDappWithdraw) sdk.Result {
	withdrawDeposit, withdrawRewards, withdrawLucky, fee, err := k.WithdrawGrid(ctx, msg.DappID, msg.GridId, msg.Sender)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			types.EventTypeGridWithdraw,
			sdk.NewAttribute(types.AttributeDappID, fmt.Sprintf("%d", msg.DappID)),
			sdk.NewAttribute(types.AttributeGridId, fmt.Sprintf("%d", msg.GridId)),
			sdk.NewAttribute(types.AttributeWithdrawDeposit, withdrawDeposit.String()),
			sdk.NewAttribute(types.AttributeWithdrawRewards, withdrawRewards.String()),
			sdk.NewAttribute(types.AttributeWithdrawLucky, withdrawLucky.String()),
			sdk.NewAttribute(types.AttributeWithdrawFee, fee.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappWithdrawFees(ctx sdk.Context, k Keeper, msg types.MsgDappWithdrawFees) sdk.Result {
	withdraw, fee, err := k.WithdrawFees(ctx, msg.DappID, msg.Sender, msg.To)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeWithdraw, withdraw.String()),
			sdk.NewAttribute(types.AttributeWithdrawFee, fee.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappWithdrawLucky(ctx sdk.Context, k Keeper, msg types.MsgDappWithdrawLucky) sdk.Result {

	fees, err := k.WithdrawLucky(ctx, msg.Sender, msg.DappID, msg.Amount)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeAddress, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeWithdrawFee, fees.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDappGenerate(ctx sdk.Context, k Keeper, msg types.MsgDappGenerate) sdk.Result {
	msg.Dapp.Owner = msg.Sender
	if err := k.GenerateDapp(ctx, &msg.Dapp); err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeDappID, fmt.Sprintf("%d", msg.Dapp.ID)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
func handleMsgDisableDapp(ctx sdk.Context, k Keeper, msg types.MsgDisableDapp) sdk.Result {
	if err := k.DisableDapp(ctx, msg.Sender, msg.DappID, msg.Height); err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeDappID, fmt.Sprintf("%d", msg.DappID)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgGridWithdrawFees(ctx sdk.Context, k Keeper, msg types.MsgGridWithdrawFees) sdk.Result {
	if err := k.CheckOwner(ctx, msg.Sender); err != nil {
		return err.Result()
	}
	withdraw, err := k.WithdrawGrid999Fee(ctx, msg.To)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeWithdraw, withdraw.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
