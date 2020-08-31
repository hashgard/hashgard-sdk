package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// HashGard
const (
	// ProposalTypeStakeIssueLockedSpend defines the type for a StakeIssueLockedSpendProposal
	ProposalTypeStakeIssueLockedSpend = "StakeIssueLockedSpend"
)

// Assert StakeIssueLockedSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = StakeIssueLockedSpendProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeStakeIssueLockedSpend)
	govtypes.RegisterProposalTypeCodec(StakeIssueLockedSpendProposal{}, "cosmos-sdk/StakeIssueLockedSpendProposal")
}

// StakeIssueLockedSpendProposal spends from the community pool
type StakeIssueLockedSpendProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Denom       string         `json:"denom" yaml:"denom"`
	Recipient   sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount      sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewStakeIssueLockedSpendProposal creates a new community pool spned proposal.
func NewStakeIssueLockedSpendProposal(title, description, denom string, recipient sdk.AccAddress, amount sdk.Coins) StakeIssueLockedSpendProposal {
	return StakeIssueLockedSpendProposal{title, description, denom, recipient, amount}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp StakeIssueLockedSpendProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp StakeIssueLockedSpendProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp StakeIssueLockedSpendProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp StakeIssueLockedSpendProposal) ProposalType() string {
	return ProposalTypeStakeIssueLockedSpend
}

// ValidateBasic runs basic stateless validity checks
func (csp StakeIssueLockedSpendProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, csp)
	if err != nil {
		return err
	}
	if !csp.Amount.IsValid() {
		return sdk.ErrInternal("amount is invalid")
	}
	if csp.Recipient.Empty() {
		return sdk.ErrInternal("recipient is invalid")
	}
	return nil
}

// String implements the Stringer interface.
func (csp StakeIssueLockedSpendProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`StakeIssueLocked  Spend Proposal:
  Title:       %s
  Description: %s
  Recipient:   %s
  Amount:      %s
`, csp.Title, csp.Description, csp.Recipient, csp.Amount))
	return b.String()
}
