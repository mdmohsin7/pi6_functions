package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type finalTallyResult struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	Abstain    string `json:"abstain"`
	NoWithVeto string `json:"no_with_veto"`
}

type totalDeposit struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type proposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type proposal struct {
	ProposalID       string           `json:"proposal_id"`
	Content          proposalContent  `json:"content"`
	Status           string           `json:"status"`
	FinalTallyResult finalTallyResult `json:"final_tally_result"`
	SubmitTime       string           `json:"submit_time"`
	DepositEndTime   string           `json:"deposit_end_time"`
	TotalDeposit     []totalDeposit   `json:"total_deposit"`
	VotingStartTime  string           `json:"voting_start_time"`
	VotingEndTime    string           `json:"voting_end_time"`
}
type pagination struct {
	NextKey string `json:"next_key"`
	Total   int    `json:"total"`
}
type proposalsFromChain struct {
	Proposals  []proposal `json:"proposals"`
	Pagination pagination `json:"pagination"`
}

type propsResponse struct {
	Status       int                `json:"status,omitempty"`
	ResponseText string             `json:"responseText,omitempty"`
	Proposals    proposalsFromChain `json:"proposals,omitempty"`
}

func ProposalHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	propStatus := r.URL.Query().Get("prop_status")
	chain := strings.ToLower(r.URL.Query().Get("chain"))
	var api string
	if chain != "" && propStatus != "" {

	}
	path := "/cosmos/gov/v1beta1/proposals?proposal_status="
	if chain == "cre" {
		api = "https://mainnet.crescent.network:1317" + path + propStatus
	} else if chain == "atom" {
		api = "https://rest-cosmoshub.ecostake.com" + path + propStatus
	} else if chain == "akt" {
		api = "https://akash.c29r3.xyz:443/api" + path + propStatus
	}

	res, _ := http.Get(api)
	resBody, _ := ioutil.ReadAll(res.Body)
	resJson := string(resBody)
	_, _ = fmt.Fprintf(w, resJson)
}
