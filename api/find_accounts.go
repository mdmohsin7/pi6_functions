package api

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"log"
	"net/http"
)

var Prefixes = map[string]string{
	"agoric":        "agoric",
	"akash":         "akash",
	"axelar":        "axelar",
	"bandchain":     "band",
	"cosmoshub":     "cosmos",
	"crescent":      "cre",
	"evmos":         "evmos",
	"injective":     "inj",
	"juno":          "juno",
	"kava":          "kava",
	"nomic":         "nomic",
	"osmosis":       "osmo",
	"regen":         "regen",
	"secretnetwork": "secret",
	"sommelier":     "somm",
	"stargaze":      "stars",
}

type accResponse struct {
	Status       int               `json:"status,omitempty"`
	Addresses    map[string]string `json:"addresses,omitempty"`
	ResponseText string            `json:"responseText,omitempty"`
}

func FindAccounts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("addr")
	var finalRes *accResponse
	accounts := make(map[string]string)
	_, b64, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return
	}
	for k, v := range Prefixes {
		addr, e := bech32.ConvertAndEncode(v, b64)
		if e != nil {
			log.Println(k, e)
		}
		accounts[v] = addr
	}
	if err != nil {
		return
	}

	finalRes = &accResponse{
		Addresses:    accounts,
		Status:       200,
		ResponseText: "success",
	}
	resData, _ := json.Marshal(finalRes)
	_, _ = fmt.Fprintf(w, string(resData))
}
