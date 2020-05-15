package types

import (
	"time"
)

type TeamDetails struct {
	Clients              []Client
	TeamName             string
	ChainID              string
	AvailableDoubloons   AvailableDoubloons
	NumberOfTransactions int
	AvgPayedGas          int
	AvgRequiredGas       int
	Address              string
}

type AvailableDoubloons struct {
	Balance int
	Height  int
}

type Client struct {
	ClientID             string
	StartBlock           int
	EndBlock             int
	StartTime            time.Time
	EndTime              time.Time
	TrustPeriodInSeconds int
	NumberOfTransactions int
	AvgPayedGas          int
	AvgRequiredGas       int
}

type TeamChart struct {
	ClientData         map[string]ClientData
	LastUpdateByClient map[string]int
}

// ClientData is a map of updateTime as unix seconds mapping the seconds since last update time
type ClientData = map[int]int
