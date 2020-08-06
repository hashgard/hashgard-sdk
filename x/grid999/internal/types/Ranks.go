package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Ranks struct {
	Rank2 [][]int `json:"rank2" yaml:"rank2"`
	Rank3 [][]int `json:"rank3" yaml:"rank3"`
	Rank4 [][]int `json:"rank4" yaml:"rank4"`
	Rank5 [][]int `json:"rank5" yaml:"rank5"`
	Rank6 [][]int `json:"rank6" yaml:"rank6"`
	Rank7 [][]int `json:"rank7" yaml:"rank7"`
	Rank8 [][]int `json:"rank8" yaml:"rank8"`
	Rank9 [][]int `json:"rank9" yaml:"rank9"`
}

func DefaultRanks() Ranks {
	return Ranks{
		Rank2: [][]int{{1}, {2}},
		Rank3: [][]int{{1}, {2}, {3}},
		Rank4: [][]int{{1}, {2, 3}, {4}},
		Rank5: [][]int{{1}, {2, 3, 4}, {5}},
		Rank6: [][]int{{1, 2}, {3, 4}, {5, 6}},
		Rank7: [][]int{{1, 2}, {3, 4, 5}, {6, 7}},
		Rank8: [][]int{{1, 2, 3}, {4, 5}, {6, 7, 8}},
		Rank9: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
	}
}
func (r Ranks) GetMaxRank() int {
	for i := 3; i <= MaxGrid; i++ {
		ranks := r.GetRank(i)
		if len(ranks) == 0 {
			return i - 1
		}
	}
	return MaxGrid
}

func (r Ranks) GetRank(rank int) [][]int {
	switch rank {
	case 2:
		return r.Rank2
	case 3:
		return r.Rank3
	case 4:
		return r.Rank4
	case 5:
		return r.Rank5
	case 6:
		return r.Rank6
	case 7:
		return r.Rank7
	case 8:
		return r.Rank8
	case 9:
		return r.Rank9
	default:
		return [][]int{}
	}
}
func (r Ranks) GetWinner(ranks [][]int) []int {
	if len(ranks) == 0 {
		return make([]int, 0)
	}
	return ranks[0]
}
func (r Ranks) GetLoser(ranks [][]int) []int {
	rankSize := len(ranks)
	switch rankSize {
	case 3:
		return ranks[2]
	case 2:
		return ranks[1]
	default:
		return make([]int, 0)
	}
}
func (r Ranks) GetFeePayer(ranks [][]int) []int {
	if len(ranks) == 3 {
		return ranks[1]
	}
	return make([]int, 0)
}
func (r Ranks) ValidateBasic() sdk.Error {
	for i := 2; i <= MaxGrid; i++ {
		ranks := r.GetRank(i)
		if len(ranks) == 0 {
			continue
		}
		if len(ranks) != 2 && len(ranks) != 3 {
			return sdk.ErrInternal("wrong rank1")
		}
		if len(ranks[0]) > 3 {
			return sdk.ErrInternal("wrong rank2")
		}
		temps := make([]int, 0)
		for _, v1 := range ranks {
			temps = append(temps, v1...)
		}
		if len(temps) != i || temps[0] != 1 || temps[len(temps)-1] != len(temps) {
			return sdk.ErrInternal("wrong rank3")
		}
		for i, vt := range temps {
			if i > 0 && vt-temps[i-1] != 1 {
				return sdk.ErrInternal("wrong rank4")
			}
		}
	}
	return nil
}
