package types

import (
	"fmt"
)

type GenerateParams struct {
	Owners string `json:"owners" yaml:"owners"`
	Dapp   Dapp   `json:"dapp" yaml:"dapp"`
}

func DefaultGenerateParams() GenerateParams {
	return GenerateParams{
		Dapp: DefaultDapp(),
	}
}
func (c GenerateParams) String() string {
	return fmt.Sprintf(`GenerateParams:
		Owner:    %s
		Dapp:    %s`, c.Owners, c.Dapp.String())
}
