package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pi6_functions/api_utils"
	"strings"
)

type balResponse struct {
	Status       int                 `json:"status,omitempty"`
	Balances     []api_utils.Balance `json:"balances,omitempty"`
	Chain        string              `json:"chain,omitempty"`
	Address      string              `json:"address,omitempty"`
	ResponseText string              `json:"responseText,omitempty"`
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
		} else if chain == "cosmos" {
			api = "https://rest-cosmoshub.ecostake.com" + path + addr
		} else if chain == "akash" {
			api = "https://akash.c29r3.xyz:443/api" + path + addr
		} else {
			wrongChain = true
		}
		if wrongChain == false {
			res, apiErr := http.Get(api)
			if apiErr != nil {
				fmt.Println(apiErr)
			}

			if res.StatusCode == 200 {
				resBody, _ := ioutil.ReadAll(res.Body)
				resJson := string(resBody)
				var balanceFromChain api_utils.BalanceFromChain
				err := json.Unmarshal([]byte(resJson), &balanceFromChain)
				if err != nil {
					return
				}
				finalRes = &balResponse{
					Balances:     balanceFromChain.Balances,
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
