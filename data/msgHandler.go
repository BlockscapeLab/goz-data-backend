package data

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	chainTypes "github.com/BlockscapeLab/goz-data-backend/chain/types"
	"github.com/BlockscapeLab/goz-data-backend/data/types"
)

// Cosmos Message Types
const (
	CREATE_CLIENT_MESSAGE = "ibc/client/tendermint/MsgCreateClient"
	UPDATE_CLIENT_MESSAGE = "ibc/client/tendermint/MsgUpdateClient"
	SEND_MESSAGE          = "cosmos-sdk/MsgSend"
)

func txDataFromTx(tx *chainTypes.Txs) types.TxData {
	h, err := strconv.Atoi(tx.Height)
	if err != nil {
		log.Println("Couldn't cast tx height to int:", err)
		h = 0
	}
	gu, err := strconv.Atoi(tx.GasUsed)
	if err != nil {
		log.Println("Couldn't cast gasUsed to int:", err)
		gu = 0
	}
	gw, err := strconv.Atoi(tx.GasWanted)
	if err != nil {
		log.Println("Couldn't cast gasWanted height to int:", err)
		gw = 0
	}

	return types.TxData{
		GasUsed:   gu,
		GasWanted: gw,
		Timestamp: tx.Timestamp,
		Height:    h,
	}
}

func (dp *DataProvider) msgHandler(tx *chainTypes.Txs) error {
	txData := txDataFromTx(tx)
	teamChainID := ""

	for _, msg := range tx.Tx.Value.Msg {
		switch msg.Type {
		case CREATE_CLIENT_MESSAGE:
			ccm := chainTypes.MsgCreateClient{}
			var err error
			if err = json.Unmarshal(msg.Value, &ccm); err != nil {
				return err
			}
			if teamChainID, err = dp.handleCreateClientMsg(ccm, txData); err != nil {
				return err
			}
		case UPDATE_CLIENT_MESSAGE:
			var err error
			ucm := chainTypes.MsgUpdateClient{}
			if err = json.Unmarshal(msg.Value, &ucm); err != nil {
				return err
			}
			if teamChainID, err = dp.handleUpdateClientMsg(ucm, txData); err != nil {
				return err
			}
		case SEND_MESSAGE:
			sm := chainTypes.MsgSend{}
			if err := json.Unmarshal(msg.Value, &sm); err != nil {
				return err
			}
			if err := dp.handleSendMsg(sm, txData); err != nil {
				return err
			}
		default:
			log.Println("Encountered unhandled message type", msg.Type)
		}
	}

	if err := dp.applyTxDataToTeam(txData, teamChainID); err != nil {
		log.Printf("couldn't update team %s. Will try at a later time: %s\n", teamChainID, err)
	}

	return nil
}

func (dp *DataProvider) handleCreateClientMsg(msg chainTypes.MsgCreateClient, txData types.TxData) (string, error) {
	chainID := msg.Header.SignedHeader.Header.ChainID
	clientID := msg.ClientID
	creationTime := int(txData.Timestamp.Unix())

	teamChart, ok := dp.chartData[chainID]
	if !ok {
		log.Println("Instantiating team chart for team", chainID)
		lubc := make(map[string]int)
		lubc[clientID] = creationTime

		cdm := make(map[string]types.ClientData)
		cd := make(map[int]int)
		cd[creationTime] = 0

		cdm[clientID] = cd

		teamChart = types.TeamChart{
			ClientData:         cdm,
			LastUpdateByClient: lubc,
		}

		dp.chartData[chainID] = teamChart
	} else {
		log.Printf("Adding new client (%s) to team %s\n", clientID, chainID)
		cd := make(map[int]int)
		cd[creationTime] = 0
		teamChart.ClientData[clientID] = cd
		teamChart.LastUpdateByClient[clientID] = creationTime
	}

	team, ok := dp.teams[chainID]
	if !ok {
		//create  new team
		team = types.TeamDetails{
			Address: msg.Address,
			AvailableDoubloons: types.AvailableDoubloons{
				Balance: 1250000,
				Height:  0,
			},
			ChainID:              chainID,
			NumberOfTransactions: 1,
			Clients:              make([]types.Client, 0),
		}
	}

	trustPerNano, _ := strconv.Atoi(msg.TrustingPeriod)

	team.Clients = append(team.Clients, types.Client{
		ClientID:             clientID,
		TrustPeriodInSeconds: trustPerNano / 1000000000,
		StartTime:            txData.Timestamp,
		StartBlock:           txData.Height,
	})

	dp.teams[chainID] = team

	return chainID, nil
}

func (dp *DataProvider) handleUpdateClientMsg(msg chainTypes.MsgUpdateClient, txData types.TxData) (string, error) {
	clientID := msg.ClientID
	chainID := msg.Header.SignedHeader.Header.ChainID
	updateTime := txData.Timestamp

	teamChart, ok := dp.chartData[chainID]
	if !ok {
		return "", fmt.Errorf("Received an update_client for unknown team with chainID '%s'. update_client is never the first message of a team.", chainID)
	}

	clientData, ok := teamChart.ClientData[clientID]
	if !ok {
		return "", fmt.Errorf("Received an update_client for unknown client '%s'. A create-client should have happened before.", clientID)
	}

	timeOfLastUpdate, ok := teamChart.LastUpdateByClient[clientID]
	if !ok {
		return "", fmt.Errorf("Received an update_client for client '%s' with no recorded previous update.", clientID)
	}

	team, ok := dp.teams[chainID]
	if !ok {
		return "", fmt.Errorf("Received update_client but couldn't find team details")
	}

	found := false
	for i, c := range team.Clients {
		if c.ClientID == clientID {
			found = true
			c.EndBlock = txData.Height
			c.EndTime = txData.Timestamp
			team.Clients[i] = c
			break
		}
	}
	if !found {
		return "", fmt.Errorf("Received update_client but couldn't find client in team details")
	}

	unix := int(updateTime.Unix())
	timeSinceLastUpdate := unix - timeOfLastUpdate

	clientData[unix] = timeSinceLastUpdate
	teamChart.LastUpdateByClient[clientID] = unix

	return chainID, nil
}

func (dp *DataProvider) handleSendMsg(msg chainTypes.MsgSend, txData types.TxData) error {
	return nil
}

func (dp *DataProvider) applyTxDataToTeam(txd types.TxData, teamChainID string) error {
	team, ok := dp.teams[teamChainID]
	if !ok {
		return fmt.Errorf("Couldn't update team %s. Team not found in map", teamChainID)
	}

	team.TotalPayedGas = team.TotalPayedGas + txd.GasWanted
	team.TotalRequiredGas = team.TotalRequiredGas + txd.GasUsed
	team.NumberOfTransactions = team.NumberOfTransactions + 1

	bal, height, err := dp.lcd.GetDoubloonsOfAccount(team.Address)
	if err != nil {
		return err
	}
	team.AvailableDoubloons = types.AvailableDoubloons{
		Balance: bal,
		Height:  height,
	}

	dp.teams[teamChainID] = team
	return nil
}
