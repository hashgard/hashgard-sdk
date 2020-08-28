package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
	"strings"
)

type GridItem struct {
	Index         uint           `json:"index" yaml:"index"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
	OwnerDeposit  sdk.Coin       `json:"owner_deposit" yaml:"owner_deposit"`
	MemberDeposit sdk.Coin       `json:"member_deposit" yaml:"member_deposit"`
	OwnerNumber   int64          `json:"owner_number" yaml:"owner_number"`
	TotalNumber   int64          `json:"total_number" yaml:"total_number"`
	Deposits      []string       `json:"deposits" yaml:"deposits"`
	GridType      string         `json:"grid_type" yaml:"grid_type"`
	ZeroValued    bool           `json:"zero_valued" yaml:"zero_valued"`
	Prepaid       string         `json:"prepaid" yaml:"prepaid"`
}

type Grid struct {
	DappID       uint       `json:"dapp_id" yaml:"dapp_id"`
	ID           uint64     `json:"id" yaml:"id"`
	Height       int64      `json:"height" yaml:"height"`
	TotalDeposit sdk.Coins  `json:"total_deposit" yaml:"total_deposit"`
	Items        []GridItem `json:"items" yaml:"items"`
	Fees         []string   `json:"fees" yaml:"fees"`
	Lucky        []string   `json:"lucky" yaml:"lucky"`
	Rewards      []string   `json:"rewards" yaml:"rewards"`
	Deposits     []string   `json:"deposits" yaml:"deposits"`
	Backs        []string   `json:"backs" yaml:"backs"`
	Prepaid      []string   `json:"prepaid" yaml:"prepaid"`
}
type Grids []Grid

func DefaultGrid(height int64, dappID uint) *Grid {
	return &Grid{
		Height:       height,
		DappID:       dappID,
		TotalDeposit: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())),
		Items:        make([]GridItem, 0),
		Fees:         make([]string, 0),
		Lucky:        make([]string, 0),
		Rewards:      make([]string, 0),
		Backs:        make([]string, 0),
	}
}
func (c GridItem) String() string {
	return fmt.Sprintf(`Grid:
		Owner:    %s
		OwnerDeposit:    %s`, c.Owner.String(), c.OwnerDeposit.String())
}
func (c Grids) String() string {
	out := ""
	for _, v := range c {
		out += v.String()
	}
	return strings.TrimSpace(out)
}
func (c Grid) String() string {
	return fmt.Sprintf(`Grid:
		Height:    %d
		Items:     %s`, c.Height, c.Items)
}

func (c Grid) Sort() {
	sort.Slice(c.Items, func(i, j int) bool {
		return c.Items[i].OwnerDeposit.Amount.GT(c.Items[j].OwnerDeposit.Amount)
	})
}
func (c Grid) SortByNumber(rankType uint) {
	sort.Slice(c.Items, func(i, j int) bool {
		if rankType == 1 {
			return c.Items[i].TotalNumber > c.Items[j].TotalNumber
		}
		coinA := c.add(c.Items[j].OwnerDeposit, c.Items[j].MemberDeposit)
		coinB := c.add(c.Items[i].OwnerDeposit, c.Items[i].MemberDeposit)
		return coinA.IsLT(coinB)
	})
}
func (c Grid) add(coinA sdk.Coin, coinB sdk.Coin) sdk.Coin {
	if coinA.Denom != coinB.Denom {
		return coinA
	}
	return coinA.Add(coinB)
}
