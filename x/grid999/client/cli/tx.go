package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Grid999 engine transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdGenerate(cdc),
		GetCmdDisable(cdc),
		GetCmdCreateGrid(cdc),
		GetCmdDepositGrid(cdc),
		GetCmdWithdrawGrid(cdc),
		GetCmdUpdateGridParams(cdc),
		GetCmdWithdrawFeeGrid(cdc),
		GetCmdWithdrawDappFeeGrid(cdc),
		GetCmdWithdrawLuckyDeposit(cdc),
	)...)

	return txCmd
}

func GetCmdCreateGrid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [dapp_id] [amount] [zero_valued] [locked]",
		Short:   "Create a grid from a dapp",
		Args:    cobra.ExactArgs(4),
		Example: fmt.Sprintf(`$ %s tx %s create 100000000ugard false false`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}
			zeroValued, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			locked, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			msg := types.MsgDappCreateGrid{Sender: cliCtx.GetFromAddress(),
				DappID:     uint(dappID),
				Deposit:    amount,
				ZeroValued: zeroValued,
				Locked:     locked}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdDepositGrid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deposit [dapp_id] [grid_id] [index] [amount]",
		Short:   "Deposit to the dapp grid",
		Args:    cobra.ExactArgs(4),
		Example: fmt.Sprintf(`$ %s tx %s deposit 100000 1 1 1000000000ugard`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			gridId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			index, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoin(args[3])
			if err != nil {
				return err
			}
			msg := types.MsgDappDeposit{DappID: uint(dappID), Sender: cliCtx.GetFromAddress(), Deposit: amount, GridId: gridId, Index: uint(index)}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdWithdrawGrid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw [dapp_id] [grid-id]",
		Short:   "Depositor withdraw rewards & deposit & lucky from the dapp grid",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf(`$ %s tx %s withdraw 100000 1`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			gridId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			msg := types.MsgDappWithdraw{DappID: uint(dappID), Sender: cliCtx.GetFromAddress(), GridId: gridId}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdWithdrawDappFeeGrid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fee [dapp_id] [to]",
		Short:   "Dapp owner withdraw fees from the dapp",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf(`$ %s tx %s withdraw-fee 100000 <address>`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgDappWithdrawFees{DappID: uint(dappID), Sender: cliCtx.GetFromAddress(), To: addr}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdWithdrawFeeGrid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-grid-fee [to]",
		Short:   "Grid999 engine owner withdraw fees",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`$ %s tx %s withdraw-grid-fee <address>`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgGridWithdrawFees{Sender: cliCtx.GetFromAddress(), To: addr}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdUpdateGridParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [path/grid999_params.json]",
		Short: "Update grid999 engine params",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`$ %s tx %s update-params grid_params.json
json file：
{
	"per_grid_max_deposits": 500,
	"generate_dapp_fee": {
		"denom": "ugard",
		"amount": "1000000000"
	},
	"create_grid_fee": {
		"denom": "ugard",
		"amount": "1000000000"
	},
	"deposit_fee": {
		"denom": "ugard",
		"amount": "0"
	},
	"withdraw_rewards_fee": {
		"denom": "ugard",
		"amount": "10000000"
	},
	"fee_withdraw_fee": "0.2"
}
`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			params := types.Params{}
			err = json.Unmarshal(contents, &params)
			if err != nil {
				return err
			}
			msg := types.MsgGridParams{Sender: cliCtx.GetFromAddress(), Params: params}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdWithdrawLuckyDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-lucky [dapp_id] [amount]",
		Short:   "Dapp owner withdraw lucky deposit to community pool addr",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf(`$ %s tx %s withdraw-lucky 100000 1000ugard`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgDappWithdrawLucky{DappID: uint(dappID), Sender: cliCtx.GetFromAddress(), Amount: amount}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdGenerate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-dapp [path/generate_params.json]",
		Short: "Generate a dapp from grid999 engine",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`$ %s tx %s generate-dapp generate_params.json
json file：
{
	"max_blocks_grid_create": 9,
	"max_blocks_grid_deposit": 9,
	"grid_create_can_deposit": true,
	"voted": false,
	"only_owner_can_create_grid": false,
	"rand_number_negative_critical_value": 50,
	"per_grid_max_deposits": 9,
	"max_blocks_grid_rewards_withdraw": 0,
	"max_blocks_grid_deposit_withdraw": 100,
	"rank_type": 1,
	"min_deposit_amount": [{
		"denom": "ugard",
		"amount": "0"
	}],
	"max_deposit_amount": [{
		"denom": "ugard",
		"amount": "0"
	}],
	"owner_min_deposit": {
		"denom": "ugard",
		"amount": "10000000000"
	},
	"member_min_deposit": {
		"denom": "ugard",
		"amount": "1000000000"
	},
	"max_per_deposit": {
		"denom": "ugard",
		"amount": "100000000000"
	},
	"owner_rewards_ratio": "0.2",
	"ree_ratio": "0.2",
	"lucky_pool_ratio": "0.2",
	"lucky_pool_rewards_digit": 2,
	"deposit_to_lucky_pool_digit": 2,
	"ranks": {
		"rank2": [
			[1],
			[2]
		],
		"rank3": [
			[1],
			[2],
			[3]
		],
		"rank4": [
			[1],
			[2, 3],
			[4]
		],
		"rank5": [
			[1],
			[2, 3, 4],
			[5]
		],
		"rank6": [
			[1, 2],
			[3, 4],
			[5, 6]
		],
		"rank7": [
			[1, 2],
			[3, 4, 5],
			[6, 7]
		],
		"rank8": [
			[1, 2, 3],
			[4, 5],
			[6, 7, 8]
		],
		"rank9": [
			[1, 2, 3],
			[4, 5, 6],
			[7, 8, 9]
		]
	},
	"winner_rewards": {
		"winner1": ["1"],
		"winner2": ["0.6", "0.4"],
		"winner3": ["0.5", "0.3", "0.2"]
	}
}`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			params := types.Dapp{}
			err = json.Unmarshal(contents, &params)
			if err != nil {
				return err
			}

			msg := types.MsgDappGenerate{Sender: cliCtx.GetFromAddress(), Dapp: params}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
func GetCmdDisable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable-dapp [dapp_id] [height]",
		Short:   "Dapp owner disable at the specified height",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf(`$ %s tx %s gdisable-dapp 100000 88888`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			dappID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			height, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			msg := types.MsgDisableDapp{Sender: cliCtx.GetFromAddress(), DappID: uint(dappID), Height: height}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			return nil
		},
	}
	return cmd
}
