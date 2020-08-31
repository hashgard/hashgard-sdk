package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// HashGard
// GetCmdStakeIssueToken implements the create validator command handler.
func GetCmdStakeIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token [path/stake_issue_token.json]",
		Args:  cobra.ExactArgs(1),
		Short: "issue a new token if validator self delegation more than required stake gard and lock ggt",
		Example: fmt.Sprintf(`$ %s tx %s issue-token stake_issue_token.json

Notes for node owners to issue TOKEN:
1. The issuance of TOKEN must be at least self-staking GARD and locked GGT according to the rules. Please check the amount and GGT lock period by querying the staking issuance TOKEN configuration command;
2. To issue TOKEN, you can configure whether to pre-mine the specified amount to the specified address;
3. You can mine the tokens issued by your node only when your node propose blocks. You can specify the amount of your node and all voting nodes that can be mined in each block;
4. After your node unbonded, the minimum staking part of the self-staking GARD according to the issuance rules will be transferred to the lock-up address corresponding to the issued TOKEN, and the rest will be returned to the owner's address;
5. After your node unbonded, the locked GGT will be transferred to the locked address corresponding to the issued TOKEN according to the rules within the locked height, otherwise it will be returned to the owner's address;
6. After your node unbonded, anyone can submit a proposal to withdraw GARD and GGT from the locked address corresponding to the issued TOKEN, and all nodes will vote to agree whether they can be withdrawn;

json file：
The default decimals is 6 digits,denom start with 'u',you can choose to pre mint or not.
{
	"denom": "uexp",
	"total_supply": "100000000000000",
	"pre_mint_address": "",
	"pre_mint_amount": "0",
	"description": {
		"whole_name": "Experience",
		"website": "https://www.exp.top/",
		"icon": "https://www.exp.top/img/icon.png",
		"details": "exp"
	},
	"genesis_height": 100000,
	"per_block_mint": [{
			"start_height": 100000,
			"self_node_amount": "10000000",
			"others_node_amount": "5000000"
		},
		{
			"start_height": 8200000,
			"self_node_amount": "6000000",
			"others_node_amount": "600000"
		}
	]
}`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			issueToken := &types.StakeIssueToken{}
			err = json.Unmarshal(contents, issueToken)
			if err != nil {
				return err
			}

			msg := types.MsgStakeIssueToken{Sender: cliCtx.GetFromAddress(), StakeIssueToken: *issueToken}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdStakeIssueTokenEdit implements the create validator command handler.
func GetCmdStakeIssueTokenEdit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token-edit [path/stake_issue_token_edit.json]",
		Args:  cobra.ExactArgs(1),
		Short: "edit the description of the issued token",
		Example: fmt.Sprintf(`$ %s tx %s issue-token-edit stake_issue_token_edit.json
json file：
{
	"denom": "uexp",
	"description": {
		"whole_name": "Experience",
		"website": "https://www.exp.top/",
		"icon": "https://www.exp.top/img/icon.png",
		"details": "exp"
	}
}`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgStakeIssueTokenEdit{}
			err = json.Unmarshal(contents, msg)
			if err != nil {
				return err
			}
			msg.Sender = cliCtx.GetFromAddress()
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdStakeIssueTokenConfig implements the create validator command handler.
func GetCmdStakeIssueTokenConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token-config [path/stake_issue_token_config.json]",
		Args:  cobra.ExactArgs(1),
		Short: "Configure token issuance rule parameters",
		Example: fmt.Sprintf(`$ %s tx %s issue-token-config stake_issue_token_config.json
json file：
{
	"min_self_delegation": "1000000",
    "lock_period_height":10,
	"lock_coins": [{
			"denom": "uggt",
			"amount": "10000000"
		}]
}`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			config := &types.StakeIssueTokenConfig{}
			err = json.Unmarshal(contents, config)
			if err != nil {
				return err
			}

			msg := types.MsgStakeIssueTokenConfig{Sender: cliCtx.GetFromAddress(), Config: *config}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
