package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group migrate queries under a subcommand
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the grid999 engine module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(client.PostCommands(
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryDappList(queryRoute, cdc),
		GetCmdQueryDapp(queryRoute, cdc),
		GetCmdQueryGridList(queryRoute, cdc),
		GetCmdQueryGrid(queryRoute, cdc))...)

	return queryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current grid999 engine parameters",
		Args:  cobra.NoArgs,
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s params`, version.ClientName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.OutputFormat = "json"

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParams)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}
			return cliCtx.PrintOutput(params)
		},
	}
}
func GetCmdQueryDappList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dapp-list [start_id] [limit]",
		Short: "Query dapp list",
		Args:  cobra.ExactArgs(2),
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s dapp-list 0 10`, version.ClientName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.OutputFormat = "json"
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryDappList, args[0], args[1]), nil,
			)
			if err != nil {
				return err
			}
			var grid types.Dapps
			if err := cdc.UnmarshalJSON(res, &grid); err != nil {
				return err
			}

			return cliCtx.PrintOutput(grid)
		},
	}
}
func GetCmdQueryDapp(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dapp [id]",
		Short: "Query a dapp",
		Args:  cobra.ExactArgs(1),
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s dapp 100000`, version.ClientName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.OutputFormat = "json"
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryDapp, args[0]), nil,
			)
			if err != nil {
				return err
			}
			var grid types.Dapp
			if err := cdc.UnmarshalJSON(res, &grid); err != nil {
				return err
			}

			return cliCtx.PrintOutput(grid)
		},
	}
}
func GetCmdQueryGrid(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "grid [dapp_id] [grid_id]",
		Short: "Query a grid",
		Args:  cobra.ExactArgs(2),
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s grid 1000000 1`, version.ClientName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.OutputFormat = "json"
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryGrid999, args[0], args[1]), nil,
			)
			if err != nil {
				return err
			}
			var grid types.Grid
			if err := cdc.UnmarshalJSON(res, &grid); err != nil {
				return err
			}

			return cliCtx.PrintOutput(grid)
		},
	}
}
func GetCmdQueryGridList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "grid-list [dapp_id] [start_id] [limit]",
		Short: "Query grid list",
		Args:  cobra.ExactArgs(3),
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s grid-list 1000000 0 10`, version.ClientName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.OutputFormat = "json"
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s/%s/%s", queryRoute, types.QueryGrid999List, args[0], args[1], args[2]), nil,
			)
			if err != nil {
				return err
			}
			var grid types.Grids
			if err := cdc.UnmarshalJSON(res, &grid); err != nil {
				return err
			}

			return cliCtx.PrintOutput(grid)
		},
	}
}
