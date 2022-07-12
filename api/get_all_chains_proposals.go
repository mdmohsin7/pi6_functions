package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pi6_functions/api_utils"
)

type allChainsPropsResponse struct {
	Status       int                            `json:"status,omitempty"`
	Proposals    []api_utils.ProposalsFromChain `json:"proposals,omitempty"`
	ResponseText string                         `json:"responseText,omitempty"`
}

var proposals []api_utils.ProposalsFromChain

func AllChainsPropsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	var api string
	path := "/cosmos/gov/v1beta1/proposals?proposal_status="
	propStatus := "2"
	proposals = nil
	for _, v := range api_utils.Prefixes {
		if v == "cosmos" {
			api = "https://rest-cosmoshub.ecostake.com" + path + propStatus
			res, _ := http.Get(api)
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			var prop api_utils.ProposalsFromChain
			prop.Chain = "cosmos"
			_ = json.Unmarshal([]byte(resJson), &prop)
			if len(prop.Proposals) > 0 {
				proposals = append(proposals, prop)
			}

		} else if v == "cre" {
			api = "https://mainnet.crescent.network:1317" + path + propStatus
			res, _ := http.Get(api)
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			var prop api_utils.ProposalsFromChain
			prop.Chain = "cre"
			_ = json.Unmarshal([]byte(resJson), &prop)
			if len(prop.Proposals) > 0 {
				proposals = append(proposals, prop)
			}
		} else if v == "akash" {
			api = "https://akash.c29r3.xyz:443/api" + path + propStatus
			res, _ := http.Get(api)
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			var prop api_utils.ProposalsFromChain
			prop.Chain = "akash"
			_ = json.Unmarshal([]byte(resJson), &prop)
			if len(prop.Proposals) > 0 {
				proposals = append(proposals, prop)
			}
		} else if v == "evmos" {
			api = "https://rest.bd.evmos.org:1317" + path + propStatus
			res, _ := http.Get(api)
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			var prop api_utils.ProposalsFromChain
			prop.Chain = "evmos"
			_ = json.Unmarshal([]byte(resJson), &prop)
			if len(prop.Proposals) > 0 {
				proposals = append(proposals, prop)
			}
		} else if v == "stars" {
			api = "https://api-stargaze-ia.notional.ventures" + path + propStatus
			res, _ := http.Get(api)
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			var prop api_utils.ProposalsFromChain
			prop.Chain = "stars"
			_ = json.Unmarshal([]byte(resJson), &prop)
			fmt.Println(resJson)
			if len(prop.Proposals) > 0 {
				proposals = append(proposals, prop)
			}
		}

	}
	resData, _ := json.Marshal(allChainsPropsResponse{
		Status:       200,
		Proposals:    proposals,
		ResponseText: "successful",
	})
	_, _ = fmt.Fprintf(w, string(resData))
}
