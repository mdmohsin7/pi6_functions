package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type txnResponse struct {
	Sent         txnsFromChain `json:"sent"`
	Received     txnsFromChain `json:"received"`
	Status       int           `json:"status"`
	ResponseText string        `json:"responseText"`
}

type txnsFromChain struct {
	Count      string `json:"count"`
	TotalCount string `json:"total_count"`
	Txs        []txn  `json:"txs"`
}

type txn struct {
	Height    string  `json:"height"`
	Txhash    string  `json:"txhash"`
	Data      string  `json:"data"`
	GasWanted string  `json:"gas_wanted"`
	GasUsed   string  `json:"gas_used"`
	Tx        txnData `json:"tx"`
}

type txnData struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

type value struct {
	Msg []msg `json:"msg"`
}

type msg struct {
	Type  string   `json:"type"`
	Value msgValue `json:"value"`
}

type msgValue struct {
	FromAddress string        `json:"from_address"`
	ToAddress   string        `json:"to_address"`
	Amount      []valueAmount `json:"amount"`
}

type valueAmount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("addr")
	chain := strings.ToLower(r.URL.Query().Get("chain"))
	var finalRes *txnResponse
	var sentResFromChain txnsFromChain
	var receivedResFromChain txnsFromChain
	var sentAPI string
	var receivedAPI string
	if addr != "" && chain != "" {
		path := "/txs?"
		if chain == "cre" {
			sentAPI = "https://api.crescent.pupmos.network" + path + "transfer.sender=" + addr + "&message.action=send"
			receivedAPI = "https://api.crescent.pupmos.network" + path + "transfer.recipient=" + addr + "&message.action=send"
		} else if chain == "cosmos" {
			sentAPI = "https://api.cosmos.network" + path + "transfer.sender=" + addr + "&message.action=send"
			receivedAPI = "https://api.cosmos.network" + path + "transfer.recipient=" + addr + "&message.action=send"
		} else if chain == "akash" {
			sentAPI = "https://akash.c29r3.xyz:443/api" + path + "transfer.sender=" + addr + "&message.action=send"
			receivedAPI = "https://akash.c29r3.xyz:443/api" + path + "transfer.recipient=" + addr + "&message.action=send"
		} else {
		}

		sentRes, apiErr := http.Get(sentAPI)
		if apiErr != nil {
			fmt.Println(apiErr)
		}

		resBody, _ := ioutil.ReadAll(sentRes.Body)
		resJson := string(resBody)
		err := json.Unmarshal([]byte(resJson), &sentResFromChain)
		if err != nil {
			return
		}
		fmt.Println(resJson)

		receivedRes, apiErr := http.Get(receivedAPI)
		if apiErr != nil {
			fmt.Println(apiErr)
		}

		resBody, _ = ioutil.ReadAll(receivedRes.Body)
		resJson = string(resBody)
		err = json.Unmarshal([]byte(resJson), &receivedResFromChain)
		if err != nil {
			return
		}
		fmt.Println(resJson)

		finalRes = &txnResponse{
			Sent:         sentResFromChain,
			Received:     receivedResFromChain,
			Status:       200,
			ResponseText: "successful",
		}
		resData, _ := json.Marshal(finalRes)
		_, _ = fmt.Fprintf(w, string(resData))
	}
}
