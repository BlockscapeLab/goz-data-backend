package types

type MsgSend struct {
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Amount      []Amount `json:"amount"`
}

type MsgUpdateClient struct {
	ClientID string   `json:"client_id"`
	Header   UCHeader `json:"header"`
	Address  string   `json:"address"`
}

type MsgCreateClient struct {
	ClientID        string   `json:"client_id"`
	Header          UCHeader `json:"header"`
	TrustingPeriod  string   `json:"trusting_period"`
	UnbondingPeriod string   `json:"unbonding_period"`
	MaxClockDrift   string   `json:"max_clock_drift"`
	Address         string   `json:"address"`
}

type Commit struct {
	Height     string       `json:"height"`
	Round      string       `json:"round"`
	BlockID    BlockID      `json:"block_id"`
	Signatures []Signatures `json:"signatures"`
}
type SignedHeader struct {
	Header Header `json:"header"`
	Commit Commit `json:"commit"`
}
type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Validators struct {
	Address          string `json:"address"`
	PubKey           PubKey `json:"pub_key"`
	VotingPower      string `json:"voting_power"`
	ProposerPriority string `json:"proposer_priority"`
}
type Proposer struct {
	Address          string `json:"address"`
	PubKey           PubKey `json:"pub_key"`
	VotingPower      string `json:"voting_power"`
	ProposerPriority string `json:"proposer_priority"`
}
type ValidatorSet struct {
	Validators []Validators `json:"validators"`
	Proposer   Proposer     `json:"proposer"`
}
type UCHeader struct {
	SignedHeader SignedHeader `json:"signed_header"`
	ValidatorSet ValidatorSet `json:"validator_set"`
}
