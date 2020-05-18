package data

import (
	"encoding/json"
	"time"

	"github.com/BlockscapeLab/goz-data-backend/data/types"
)

type MockProvider struct{}

func (mp MockProvider) GetScoreboardJSON() ([]byte, error) {
	sb := make(types.Scoreboard)
	sb[1] = types.RankData{
		ChainID:              "chainpions",
		ClientID:             "abc123",
		TeamName:             "Block-Shrooms",
		TrustPeriodInSeconds: 600,
		UptimeInBlocks:       467,
		UptimeInSeconds:      2802,
	}
	sb[2] = types.RankData{
		ChainID:              "Chaina-Town",
		ClientID:             "def456",
		TeamName:             "Chainese Team",
		TrustPeriodInSeconds: 600,
		UptimeInBlocks:       412,
		UptimeInSeconds:      2472,
	}
	sb[3] = types.RankData{
		ChainID:              "Chainless-Steel",
		ClientID:             "ghi789",
		TeamName:             "Chain-Forge",
		TrustPeriodInSeconds: 1000,
		UptimeInBlocks:       500,
		UptimeInSeconds:      3000,
	}
	return json.MarshalIndent(sb, "", " ")
}

func (mp MockProvider) GetTeamDetailsJSON(chainID string) ([]byte, error) {
	clients := []types.Client{}
	clients = append(clients, types.Client{
		ClientID:             "abc123",
		EndBlock:             3000,
		EndTime:              time.Now(),
		StartBlock:           2000,
		StartTime:            time.Now(),
		TrustPeriodInSeconds: 600,
	},
		types.Client{
			ClientID:             "def456",
			EndBlock:             4500,
			EndTime:              time.Now(),
			StartBlock:           3000,
			StartTime:            time.Now(),
			TrustPeriodInSeconds: 550,
		})

	td := types.TeamDetails{
		AvailableDoubloons: types.AvailableDoubloons{
			Balance: 104883,
			Height:  58320,
		},
		TotalPayedGas:        45839,
		TotalRequiredGas:     42949,
		ChainID:              chainID,
		Clients:              clients,
		NumberOfTransactions: 200,
		TeamName:             "MockTeam",
	}

	return json.MarshalIndent(td, "", " ")
}

func (mp MockProvider) GetTeamChartDataJSON(chainID string) ([]byte, error) {
	cd := make(map[string]types.ClientData)
	numbers := make(types.ClientData)
	numbers[1589533067] = 599
	numbers[1589532468] = 583
	numbers[1589531885] = 0

	cd["abc123"] = numbers
	cd["def456"] = numbers

	tc := types.TeamChart{
		ClientData: cd,
	}

	return json.MarshalIndent(tc, "", " ")
}
