package types

import (
	"time"
)

type TxData struct {
	Heigth    int
	Timestamp time.Time
	GasWanted int
	GasUsed   int
}
