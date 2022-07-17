package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type finalTxnData struct {
	Count      string     `json:"count"`
	TotalCount string     `json:"total_count"`
	Txs        []finalTxn `json:"txs"`
}

type finalTxn struct {
	Height    string   `json:"height"`
	Txhash    string   `json:"txhash"`
	Data      string   `json:"data"`
	GasWanted string   `json:"gas_wanted"`
	GasUsed   string   `json:"gas_used"`
	Value     msgValue `json:"value"`
}

type txnResponse struct {
	Sent         finalTxnData `json:"sent"`
	Received     finalTxnData `json:"received"`
	Status       int          `json:"status"`
	ResponseText string       `json:"responseText"`
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
	var finalSentTxs []finalTxn
	var receivedResFromChain txnsFromChain
	var finalReceivedTxs []finalTxn
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
		if len(sentResFromChain.Txs) > 0 {
			for i := range sentResFromChain.Txs {
				finalSentTxs = append(finalSentTxs, finalTxn{
					Height:    sentResFromChain.Txs[i].Height,
					Txhash:    sentResFromChain.Txs[i].Txhash,
					Data:      sentResFromChain.Txs[i].Data,
					GasWanted: sentResFromChain.Txs[i].GasWanted,
					GasUsed:   sentResFromChain.Txs[i].GasUsed,
					Value: msgValue{
						FromAddress: sentResFromChain.Txs[i].Tx.Value.Msg[0].Value.FromAddress,
						ToAddress:   sentResFromChain.Txs[i].Tx.Value.Msg[0].Value.ToAddress,
						Amount:      sentResFromChain.Txs[i].Tx.Value.Msg[0].Value.Amount,
					},
				})
			}
		}
		if len(receivedResFromChain.Txs) > 0 {
			for j := range receivedResFromChain.Txs {
				finalReceivedTxs = append(finalReceivedTxs, finalTxn{
					Height:    receivedResFromChain.Txs[j].Height,
					Txhash:    receivedResFromChain.Txs[j].Txhash,
					Data:      receivedResFromChain.Txs[j].Data,
					GasWanted: receivedResFromChain.Txs[j].GasWanted,
					GasUsed:   receivedResFromChain.Txs[j].GasUsed,
					Value: msgValue{
						FromAddress: receivedResFromChain.Txs[j].Tx.Value.Msg[0].Value.FromAddress,
						ToAddress:   receivedResFromChain.Txs[j].Tx.Value.Msg[0].Value.ToAddress,
						Amount:      receivedResFromChain.Txs[j].Tx.Value.Msg[0].Value.Amount,
					},
				})
			}
		}
		finalRes = &txnResponse{
			Sent: finalTxnData{
				Count:      sentResFromChain.Count,
				TotalCount: sentResFromChain.TotalCount,
				Txs:        finalSentTxs,
			},
			Received: finalTxnData{
				Count:      receivedResFromChain.Count,
				TotalCount: receivedResFromChain.TotalCount,
				Txs:        finalReceivedTxs,
			},
			Status:       200,
			ResponseText: "successful",
		}
		resData, _ := json.Marshal(finalRes)
		_, _ = fmt.Fprintf(w, string(resData))
	}
}
