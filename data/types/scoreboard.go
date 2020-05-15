package types

// Scoreboard is a list of RankData where the key is the rank
type Scoreboard = map[int]RankData

// RankData is the most relevant data for a scoreboard
type RankData struct {
	TeamName             string
	ChainID              string
	ClientID             string
	UptimeInSeconds      int
	UptimeInBlocks       int
	TrustPeriodInSeconds int
}
