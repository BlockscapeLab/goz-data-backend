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
		Heigth:    h,
	}
}

func (dp *DataProvider) msgHandler(tx *chainTypes.Txs) error {
	txData := txDataFromTx(tx)

	for _, msg := range tx.Tx.Value.Msg {
		switch msg.Type {
		case CREATE_CLIENT_MESSAGE:
			ccm := chainTypes.MsgCreateClient{}
			if err := json.Unmarshal(msg.Value, &ccm); err != nil {
				return err
			}
			if err := dp.handleCreateClientMsg(ccm, txData); err != nil {
				return err
			}
		case UPDATE_CLIENT_MESSAGE:
			ucm := chainTypes.MsgUpdateClient{}
			if err := json.Unmarshal(msg.Value, &ucm); err != nil {
				return err
			}
			if err := dp.handleUpdateClientMsg(ucm, txData); err != nil {
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
	return nil
}

func (dp *DataProvider) handleCreateClientMsg(msg chainTypes.MsgCreateClient, txData types.TxData) error {
	chainID := msg.Header.SignedHeader.Header.ChainID
	clientID := msg.ClientID
	creationTime := int(txData.Timestamp.Unix())
	// address := msg.Address

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
	// TODO tema details

	return nil
}

func (dp *DataProvider) handleUpdateClientMsg(msg chainTypes.MsgUpdateClient, txData types.TxData) error {
	// addr := msg.Address
	clientID := msg.ClientID
	chainID := msg.Header.SignedHeader.Header.ChainID
	updateTime := txData.Timestamp

	teamChart, ok := dp.chartData[chainID]
	if !ok {
		return fmt.Errorf("Received an update_client for unknown team with chainID '%s'. update_client is never the first message of a team.", chainID)
	}

	clientData, ok := teamChart.ClientData[clientID]
	if !ok {
		return fmt.Errorf("Received an update_client for unknown client '%s'. A create-client should have happened before.", clientID)
	}

	timeOfLastUpdate, ok := teamChart.LastUpdateByClient[clientID]
	if !ok {
		return fmt.Errorf("Received an update_client for client '%s' with no recorded previous update.", clientID)
	}

	unix := int(updateTime.Unix())
	timeSinceLastUpdate := unix - timeOfLastUpdate

	clientData[unix] = timeSinceLastUpdate
	teamChart.LastUpdateByClient[clientID] = unix

	// TODO update team details
	return nil
}

func (dp *DataProvider) handleSendMsg(msg chainTypes.MsgSend, txData types.TxData) error {
	return nil
}
