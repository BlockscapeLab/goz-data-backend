package chain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/BlockscapeLab/goz-data-backend/chain/types"
)

type LCDConnector struct {
	IP         string
	Port       int
	httpClient *http.Client
}

// CreateLCDConnector creates connector for rpc and also checks whether rpc endpoint can be reached
func CreateLCDConnector(lcdIP string, lcdPort int) (*LCDConnector, error) {
	c := &LCDConnector{IP: lcdIP, Port: lcdPort}
	c.httpClient = http.DefaultClient

	_, err := c.GetCurrentHeight()
	return c, err
}

// GetDubloonsOfAccount returns the current amount of dubloons the specified account owns. Should only be called when account is executing a transaction in phase 1
func (c *LCDConnector) GetDoubloonsOfAccount(bech32Addr string) (doubloons, height int, err error) {
	res, err := c.httpClient.Get(fmt.Sprintf("%s/bank/balances/%s", c.lcdAddress(), bech32Addr))
	if err != nil {
		return 0, 0, err
	}

	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, 0, err
	}
	br := types.BalanceResponse{}
	err = json.Unmarshal(bz, &br)
	doubloons, height = br.GetDubloonsAtCurrentHeight()
	return
}

func (c *LCDConnector) GetTxsOfHeight(height int) (*types.TransactionResponse, error) {
	tr := types.TransactionResponse{}
	res, err := c.httpClient.Get(fmt.Sprintf("%s/txs?tx.height=%d&limit=100", c.lcdAddress(), height))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bz, &tr)
	return &tr, err
}

// GetBlockOfHeight returns block with all necessary data for analysing
func (c *LCDConnector) GetBlockOfHeight(height int) (types.BlockResponse, error) {
	return c.getBlock(strconv.Itoa(height))
}

func (c *LCDConnector) GetLatestBlock() (types.BlockResponse, error) {
	return c.getBlock("latest")
}

func (c *LCDConnector) getBlock(b string) (types.BlockResponse, error) {
	br := types.BlockResponse{}
	res, err := c.httpClient.Get(fmt.Sprintf("%s/blocks/%s", c.lcdAddress(), b))
	if err != nil {
		return br, err
	}

	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return br, err
	}

	err = json.Unmarshal(bz, &br)

	return br, err
}

// GetCurrentHeight via status endpoint of chain rpc
func (c *LCDConnector) GetCurrentHeight() (int, error) {
	b, err := c.GetLatestBlock()
	if err != nil {
		return 0, err
	}

	lh := b.Block.Header.Height
	if lh == "" {
		return 0, nil
	} else {
		h, err := strconv.Atoi(lh)
		return h, err
	}
}

func (c *LCDConnector) lcdAddress() string {
	return fmt.Sprintf("http://%s:%d", c.IP, c.Port)
}
