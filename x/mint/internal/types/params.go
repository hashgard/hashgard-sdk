package types

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyMintDenom           = []byte("MintDenom")
	KeyInflationRateChange = []byte("InflationRateChange")
	KeyInflationMax        = []byte("InflationMax")
	KeyInflationMin        = []byte("InflationMin")
	KeyGoalBonded          = []byte("GoalBonded")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
	// HashGard
	KeyPerBlockMint        = []byte("PerBlockMint")
	KeyMinterSupplyAddress = []byte("MinterSupplyAddress")
)

//HashGard
var (
	DefaultMinterSupplyAddress = sdk.AccAddress(crypto.AddressHash([]byte("initMinterSupplyAddress")))
	DefaultPerBlockMint        = ""
)

// mint parameters
type Params struct {
	MintDenom           string  `json:"mint_denom" yaml:"mint_denom"`                       // type of coin to mint
	InflationRateChange sdk.Dec `json:"inflation_rate_change" yaml:"inflation_rate_change"` // maximum annual change in inflation rate
	InflationMax        sdk.Dec `json:"inflation_max" yaml:"inflation_max"`                 // maximum inflation rate
	InflationMin        sdk.Dec `json:"inflation_min" yaml:"inflation_min"`                 // minimum inflation rate
	GoalBonded          sdk.Dec `json:"goal_bonded" yaml:"goal_bonded"`                     // goal of percent bonded atoms
	BlocksPerYear       uint64  `json:"blocks_per_year" yaml:"blocks_per_year"`             // expected blocks per year
	// HashGard
	PerBlockMint        string         `json:"per_block_mint" yaml:"per_block_mint"`
	MinterSupplyAddress sdk.AccAddress `json:"minter_supply_address" yaml:"minter_supply_address"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflationRateChange, inflationMax,
	inflationMin, goalBonded sdk.Dec, blocksPerYear uint64, perBlockMint string, minterSupplyAddress sdk.AccAddress) Params {

	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
		// HashGard
		PerBlockMint:        perBlockMint,
		MinterSupplyAddress: minterSupplyAddress,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom: sdk.DefaultBondDenom,
		// HashGard
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 3), // assuming 3 second block times
		// HashGard
		PerBlockMint:        DefaultPerBlockMint,
		MinterSupplyAddress: DefaultMinterSupplyAddress,
	}
}

// validate params
func ValidateParams(params Params) error {
	//if params.GoalBonded.LT(sdk.ZeroDec()) {
	//	return fmt.Errorf("mint parameter GoalBonded should be positive, is %s ", params.GoalBonded.String())
	//}
	//if params.GoalBonded.GT(sdk.OneDec()) {
	//	return fmt.Errorf("mint parameter GoalBonded must be <= 1, is %s", params.GoalBonded.String())
	//}
	//if params.InflationMax.LT(params.InflationMin) {
	//	return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	//}
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	// HashGard
	if len(params.PerBlockMint) > 0 {
		if strings.Index(params.PerBlockMint, ":") == -1 {
			return perBlockMintError()
		}
		blockMints := strings.Split(params.PerBlockMint, ",")
		//100:10
		preBlockHeight := int64(-1)
		for k, v := range blockMints {
			blockMint := strings.Split(v, ":")
			if len(blockMint) == 1 {
				return perBlockMintError()
			}
			block, err := strconv.ParseInt(blockMint[0], 10, 64)
			if k == 0 && block != 0 {
				return fmt.Errorf("mint parameter PerBlockMint start blockHeight should be zero")
			}
			if preBlockHeight >= block {
				return fmt.Errorf("mint parameter PerBlockMint blockHeight should be incremental sequence")
			}
			preBlockHeight = block
			if err != nil || block < 0 {
				return perBlockMintError()
			}
			mint, err := strconv.ParseInt(blockMint[1], 10, 64)
			if err != nil || mint < 0 {
				return perBlockMintError()
			}
		}
	}

	return nil
}

// HashGard
func perBlockMintError() error {
	return fmt.Errorf("mint parameter PerBlockMint should be blockHeight:mintAmount_percent_accAddress,blockHeight:mintAmount_percent_accAddress")
}

// HashGard
func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Blocks Per Year:        %d,
  Minter Supply Address:  %s,
`,
		p.MintDenom, p.BlocksPerYear, p.MinterSupplyAddress.String(),
	)
	//	return fmt.Sprintf(`Minting Params:
	//  Mint Denom:             %s
	//  Inflation Rate Change:  %s
	//  Inflation Max:          %s
	//  Inflation Min:          %s
	//  Goal Bonded:            %s,
	//`,
	//		p.MintDenom, p.InflationRateChange, p.InflationMax,
	//		p.InflationMin, p.GoalBonded,
	//	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMintDenom, &p.MintDenom},
		{KeyInflationRateChange, &p.InflationRateChange},
		{KeyInflationMax, &p.InflationMax},
		{KeyInflationMin, &p.InflationMin},
		{KeyGoalBonded, &p.GoalBonded},
		{KeyBlocksPerYear, &p.BlocksPerYear},
		// HashGard
		{KeyPerBlockMint, &p.PerBlockMint},
		{KeyMinterSupplyAddress, &p.MinterSupplyAddress},
	}
}
