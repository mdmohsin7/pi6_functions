package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SpecificProposalHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	propID := r.URL.Query().Get("propID")
	chain := strings.ToLower(r.URL.Query().Get("chain"))
	_ = make(map[string]string)
	var api string
	path := "/cosmos/gov/v1beta1/proposals/"
	if chain == "cre" {
		api = "https://mainnet.crescent.network:1317" + path + propID
	} else if chain == "atom" {
		api = "https://rest-cosmoshub.ecostake.com" + path + propID
	} else if chain == "akt" {
		api = "https://akash.c29r3.xyz:443/api" + path + propID
	}

	res, _ := http.Get(api)
	resBody, _ := ioutil.ReadAll(res.Body)
	resJson := string(resBody)
	_, _ = fmt.Fprintf(w, resJson)
}
