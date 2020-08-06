package migrate

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
	"strings"
)

// NewHandler returns a handler for "migrate" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgERC20MigrateExchange:
			return handleMsgERC20MigrateExchange(ctx, k, msg)
		case MsgERC20Migrate:
			return handleMsgERC20Migrate(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized org message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
func handleMsgERC20MigrateExchange(ctx sdk.Context, k Keeper, msg MsgERC20MigrateExchange) sdk.Result {

	exchange := types.ERC20MigrateExchange{
		ExchangeFrom: msg.ExchangeFrom,
		Allows:       msg.Allows}

	if err := k.SetErc20MigrateExchange(ctx, exchange); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			types.EventTypeERC20MigrateExchange,
			sdk.NewAttribute(types.AttributeExchange, msg.ExchangeFrom.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}

}
func checkAllows(ctx sdk.Context, k Keeper, msg MsgERC20Migrate) (err sdk.Error, migrate types.ERC20MigrateExchange) {
	err, migrate = k.GetErc20MigrateExchange(ctx)
	if err != nil {
		return
	}
	senderStr := msg.Sender.String()
	allows := strings.Split(migrate.Allows, ",")
	for _, v := range allows {
		if v == senderStr {
			return
		}
	}
	return sdk.ErrInternal("not allowed:" + senderStr), migrate
}
func handleMsgERC20Migrate(ctx sdk.Context, k Keeper, msg MsgERC20Migrate) sdk.Result {
	params := k.GetParams(ctx)
	if params.ERC20MigrateDisabled {
		return sdk.ErrInternal("ERC20 coin migrate to mainnet coin is disabled").Result()
	}
	err, migrate := checkAllows(ctx, k, msg)
	if err != nil {
		return err.Result()
	}
	// The migrate service has duplicated control of erc 20 txhash
	//chain case 1
	if ctx.BlockHeight() < 24000 {
		if err := k.AddErc20Migrate(ctx, msg.EthTxHash); err != nil {
			return err.Result()
		}
	} else if ctx.BlockHeight() < 27000 {
		//chain case 2
		if err := k.AddErc20Migrate(ctx, msg.Erc20Address); err != nil {
			return err.Result()
		}
	} else {
		//chain case 3,txhash control again
		if err := k.AddErc20Migrate(ctx, msg.Erc20TxHash); err != nil {
			return err.Result()
		}
	}

	event := sdk.NewEvent(
		types.EventTypeERC20Migrate,
		sdk.NewAttribute(types.AttributeSigner, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeErc20Address, msg.Erc20Address),
		sdk.NewAttribute(types.AttributeErc20TxHash, msg.Erc20TxHash),
		sdk.NewAttribute(types.AttributeEthTxHash, msg.EthTxHash),
	)
	if migrate.ExchangeFrom != nil {
		if err := k.GetBankKeeper().SendCoins(ctx, migrate.ExchangeFrom, msg.To, msg.Amount); err != nil {
			return err.Result()
		}
		k.GetSceneKeeper().AddSenderTxScene(ctx, migrate.ExchangeFrom)
	} else {
		_, err := k.GetBankKeeper().AddCoins(ctx, msg.To, msg.Amount)

		if err != nil {
			return err.Result()
		}
		if k.GetSupplyKeeper().MintCoins(ctx, types.ModuleName, msg.Amount); err != nil {
			return err.Result()
		}
		event.Attributes = append(event.Attributes, sdk.NewAttribute(types.AttributeTo, msg.To.String()).ToKVPair())
	}
	if migrate.RewardsFrom != nil {
		amount := sdk.NewDecFromBigInt(msg.Amount.AmountOf(sdk.DefaultBondDenom).BigInt())
		p, _ := sdk.NewDecFromStr("0.1")
		rewards := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount.Mul(p).TruncateInt()))
		if err := k.GetBankKeeper().SendCoins(ctx, migrate.RewardsFrom, msg.To, rewards); err == nil {
			k.GetSceneKeeper().AddSenderTxScene(ctx, migrate.RewardsFrom)
		}
	}

	if msg.Channel != nil {
		event.Attributes = append(event.Attributes, sdk.NewAttribute(types.AttributeChannel, msg.Channel.String()).ToKVPair())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		event,
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
