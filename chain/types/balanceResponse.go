package types

import (
	"fmt"
	"strconv"
)

const DOUBLOON_DENOM = "doubloons"

type BalanceResponse struct {
	Height string          `json:"height"`
	Result []BalanceResult `json:"result"`
}
type BalanceResult struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (br BalanceResponse) GetDoubloonsAtCurrentHeight() (doubloons, height int) {
	var err error
	height, err = strconv.Atoi(br.Height)
	if err != nil {
		fmt.Println("Height of retrieved balance could not be parsed !!!", br)
		return 0, 0
	}
	for _, r := range br.Result {
		if r.Denom == DOUBLOON_DENOM {
			doubloons, err = strconv.Atoi(r.Amount)
			if err != nil {
				fmt.Println("Couldn't parse doubloon amount !!!", br)
			}
			return
		}
	}
	return
}
