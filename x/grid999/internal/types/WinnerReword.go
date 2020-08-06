package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type WinnerRewards struct {
	Winner1 []sdk.Dec `json:"winner1" yaml:"winner1"`
	Winner2 []sdk.Dec `json:"winner2" yaml:"winner2"`
	Winner3 []sdk.Dec `json:"winner3" yaml:"winner3"`
}

func DefaultWinnerRewards() WinnerRewards {
	return WinnerRewards{
		Winner1: []sdk.Dec{sdk.NewDec(1)},
		Winner2: []sdk.Dec{sdk.NewDecWithPrec(6, 2), sdk.NewDecWithPrec(4, 2)},
		Winner3: []sdk.Dec{sdk.NewDecWithPrec(5, 2), sdk.NewDecWithPrec(3, 2), sdk.NewDecWithPrec(2, 2)},
	}
}
func (w WinnerRewards) IsNoWinnerRewards(rank int) bool {
	rewards, _ := w.GetWinnerRewards(rank)
	total := sdk.ZeroDec()
	for _, v1 := range rewards {
		total = total.Add(v1)
	}
	return total.TruncateInt().Int64() == 0
}
func (w WinnerRewards) GetWinnerRewards(rank int) ([]sdk.Dec, bool) {
	switch rank {
	case 1:
		return w.Winner1, true
	case 2:
		return w.Winner2, true
	case 3:
		return w.Winner3, true
	default:
		return []sdk.Dec{}, false
	}
}
func (w WinnerRewards) ValidateBasic() sdk.Error {
	if len(w.Winner1) == 0 && len(w.Winner2) == 0 && len(w.Winner3) == 0 {
		return sdk.ErrInternal("wrong winner rewards")
	}
	for i := 1; i <= 3; i++ {
		rewards, _ := w.GetWinnerRewards(i)
		total := sdk.ZeroDec()
		for _, v1 := range rewards {
			if v1.IsNegative() {
				return sdk.ErrInternal("wrong winner rewards")
			}
			if checkDec(v1) {
				return sdk.ErrInternal("wrong winner rewards")
			}
			total = total.Add(v1)
		}
		if total.TruncateInt().Int64() != 1 && total.TruncateInt().Int64() != 0 {
			return sdk.ErrInternal("wrong winner rewards")
		}

	}
	return nil
}

func checkDec(dec sdk.Dec) bool {
	return len(strings.Split(fmt.Sprintf("%v", strings.Trim(dec.String(), "0")), ".")[1]) > 4
}
