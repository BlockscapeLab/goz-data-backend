package types

import (
	"fmt"
	"strconv"
)

const DUBLOON_DENOM = "dubloons"

type BalanceResponse struct {
	Height string          `json:"height"`
	Result []BalanceResult `json:"result"`
}
type BalanceResult struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (br BalanceResponse) GetDubloonsAtCurrentHeight() (dubloons, height int) {
	var err error
	height, err = strconv.Atoi(br.Height)
	if err != nil {
		fmt.Println("Height of retrieved balance could not be parsed !!!", br)
		return 0, 0
	}
	for _, r := range br.Result {
		if r.Denom == DUBLOON_DENOM {
			dubloons, err = strconv.Atoi(r.Amount)
			if err != nil {
				fmt.Println("Couldn't parse dubloon amount !!!", br)
			}
			return
		}
	}
	return
}
