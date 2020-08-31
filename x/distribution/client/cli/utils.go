package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// CommunityPoolSpendProposalJSON defines a CommunityPoolSpendProposal with a deposit
	CommunityPoolSpendProposalJSON struct {
		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Recipient   sdk.AccAddress `json:"recipient" yaml:"recipient"`
		Amount      sdk.Coins      `json:"amount" yaml:"amount"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

// ParseCommunityPoolSpendProposalJSON reads and parses a CommunityPoolSpendProposalJSON from a file.
func ParseCommunityPoolSpendProposalJSON(cdc *codec.Codec, proposalFile string) (CommunityPoolSpendProposalJSON, error) {
	proposal := CommunityPoolSpendProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// HashGard
type (
	// StakeIssueLockedSpendProposalJSON defines a StakeIssueLockedSpendProposal with a deposit
	StakeIssueLockedSpendProposalJSON struct {
		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Denom       string         `json:"denom" yaml:"denom"`
		Recipient   sdk.AccAddress `json:"recipient" yaml:"recipient"`
		Amount      sdk.Coins      `json:"amount" yaml:"amount"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

// StakeIssueLockedSpendProposalJSON reads and parses a StakeIssueLockedSpendProposalJSON from a file.
func ParseStakeIssueLockedSpendJSON(cdc *codec.Codec, proposalFile string) (StakeIssueLockedSpendProposalJSON, error) {
	proposal := StakeIssueLockedSpendProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
