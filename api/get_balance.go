package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type balanceFromChain struct {
	Balances []balance `json:"balances"`
}

type balResponse struct {
	Status       int       `json:"status,omitempty"`
	Balances     []balance `json:"balances,omitempty"`
	Chain        string    `json:"chain,omitempty"`
	Address      string    `json:"address,omitempty"`
	ResponseText string    `json:"responseText,omitempty"`
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("addr")
	chain := strings.ToLower(r.URL.Query().Get("chain"))
	var api string
	var finalRes *balResponse
	var wrongChain = false
	if addr != "" && chain != "" {
		path := "/cosmos/bank/v1beta1/balances/"
		if chain == "cre" {
			api = "https://mainnet.crescent.network:1317" + path + addr
		} else if chain == "atom" {
			api = "https://rest-cosmoshub.ecostake.com" + path + addr
		} else if chain == "akt" {
			api = "https://akash.c29r3.xyz:443/api" + path + addr
		} else {
			wrongChain = true
		}
		if wrongChain == false {
			res, apiErr := http.Get(api)
			if apiErr != nil {
				return
			}

			if res.StatusCode == 200 {
				resBody, _ := ioutil.ReadAll(res.Body)
				resJson := string(resBody)
				var BalanceFromChain balanceFromChain
				err := json.Unmarshal([]byte(resJson), &BalanceFromChain)
				if err != nil {
					return
				}
				finalRes = &balResponse{
					Balances:     BalanceFromChain.Balances,
					Chain:        chain,
					Address:      addr,
					Status:       res.StatusCode,
					ResponseText: "success",
				}
			} else {
				finalRes = &balResponse{
					Status:       res.StatusCode,
					ResponseText: "something went wrong",
				}

			}
		} else {
			finalRes = &balResponse{
				Status:       501,
				ResponseText: "The chain isn't supported or does not exist",
			}
		}

	} else {
		finalRes = &balResponse{
			Status:       501,
			ResponseText: "Address or chain name can't be empty",
		}
	}

	resData, _ := json.Marshal(finalRes)
	_, _ = fmt.Fprintf(w, string(resData))

}
