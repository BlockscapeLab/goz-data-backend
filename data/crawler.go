package data

import (
	"log"
	"time"
)

func (dp *DataProvider) crawl() {
	for {
		max, err := dp.lcd.GetCurrentHeight()
		if err != nil {
			log.Println("Error while initiating crawl, retrying in 5s:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		from := 37823 // dp.lastSyncedBlock + 1
		log.Printf("Starting crawl from block %d to %d...\n", from, max)
		for h := from; h <= max; h++ {
			err := dp.getBlockData(h)
			if err != nil {
				log.Printf("Error while fetching data for block %d. Retrying in 5s: %s\n", h, err.Error())
				h--
				time.Sleep(5 * time.Second)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (dp *DataProvider) getBlockData(height int) error {
	//get Block, if txs: continue
	br, err := dp.lcd.GetBlockOfHeight(height)
	if err != nil {
		return err
	}

	if len(br.Block.Data.Txs) == 0 {
		dp.lastSyncedBlock = height
		return nil
	}

	// gettxs
	tr, err := dp.lcd.GetTxsOfHeight(height)
	if err != nil {
		return err
	}

	// analyze txs
	for _, tx := range tr.Txs {
		if err := dp.msgHandler(&tx); err != nil {
			return err
		}

	}

	// for each tx check involved accounts balance

	// update scoreboard
	return nil
}
