package types

import (
	"time"
)

type TxData struct {
	Height    int
	Timestamp time.Time
	GasWanted int
	GasUsed   int
}
