package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"pi6_functions/api_utils"
	"strings"
)

type propsResponse struct {
	Status       int                          `json:"status,omitempty"`
	ResponseText string                       `json:"responseText,omitempty"`
	Proposals    api_utils.ProposalsFromChain `json:"proposals,omitempty"`
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
