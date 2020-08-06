package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Migrate erc20 coin to mainnet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(GetCmdSetERC20MigrateExchange(cdc))...)

	return txCmd
}

func GetCmdSetERC20MigrateExchange(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-exchange [exchange_from] [sender_allows]",
		Short:   "Set up a erc20 coin to mainnet migrate exchange gateway,Can only be set once",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf(`$ %s tx %s set-exchange "" gard1p9kz2z5yll7vp6je5sczjpwfsjym9qrqx9zqh8,gard1rns8r0rzs629avtajcttkydjhcfy3n7na0cjge --from user`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			exchangeMsg := types.MsgERC20MigrateExchange{Sender: cliCtx.GetFromAddress(), Allows: args[1]}
			if len(args[0]) > 0 {
				exchangeFrom, err := sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
				exchangeMsg.ExchangeFrom = exchangeFrom
			}

			// broadcast to a Tendermint node
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{exchangeMsg})
		},
	}
	return cmd
}
