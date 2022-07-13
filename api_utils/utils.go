package api_utils

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type BalanceFromChain struct {
	Balances []Balance `json:"balances"`
}

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

var Prefixes = map[string]string{
	"akash":         "akash",
	"bandchain":     "band",
	"cosmoshub":     "cosmos",
	"crescent":      "cre",
	"evmos":         "evmos",
	"injective":     "inj",
	"juno":          "juno",
	"osmosis":       "osmo",
	"secretnetwork": "secret",
	"stargaze":      "stars",
}

var Tickers = map[string]string{
	"akash":   "akt",
	"band":    "band",
	"cosmos":  "atom",
	"cre":     "cre",
	"evmos":   "evmos",
	"juno":    "juno",
	"osmosis": "osmo",
	"secret":  "scrt",
	"stars":   "stars",
}

type FinalTallyResult struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	Abstain    string `json:"abstain"`
	NoWithVeto string `json:"no_with_veto"`
}

type TotalDeposit struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type ProposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Proposal struct {
	ProposalID       string           `json:"proposal_id"`
	Content          ProposalContent  `json:"content"`
	Status           string           `json:"status"`
	FinalTallyResult FinalTallyResult `json:"final_tally_result"`
	SubmitTime       string           `json:"submit_time"`
	DepositEndTime   string           `json:"deposit_end_time"`
	TotalDeposit     []TotalDeposit   `json:"total_deposit"`
	VotingStartTime  string           `json:"voting_start_time"`
	VotingEndTime    string           `json:"voting_end_time"`
}
type Pagination struct {
	NextKey string `json:"next_key"`
	Total   int    `json:"total"`
}
type ProposalsFromChain struct {
	Chain      string     `json:"chain,omitempty"`
	Proposals  []Proposal `json:"proposals"`
	Pagination Pagination `json:"pagination"`
}
