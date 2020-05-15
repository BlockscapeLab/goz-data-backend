package types

// DataProvider is the iterface the rest package needs for fetching its data
type DataProvider interface {
	GetScoreboardJSON() ([]byte, error)
	GetTeamDetailsJSON(chainID string) ([]byte, error)
	GetTeamChartDataJSON(chainID string) ([]byte, error)
}
