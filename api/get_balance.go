package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("addr")
	chain := strings.ToLower(r.URL.Query().Get("chain"))
	_ = make(map[string]string)
	var api string
	path := "/cosmos/bank/v1beta1/balances/"
	if chain == "cre" {
		api = "https://mainnet.crescent.network:1317" + path + addr
	} else if chain == "atom" {
		api = "https://rest-cosmoshub.ecostake.com" + path + addr
	} else if chain == "akt" {
		api = "https://akash.c29r3.xyz:443/api" + path + addr
	}

	res, _ := http.Get(api)
	//status := res.Status
	//_, err := fmt.Fprintf(w, status)
	//if err != nil {
	//	return
	//}
	resBody, _ := ioutil.ReadAll(res.Body)
	resJson := string(resBody)
	_, _ = fmt.Fprintf(w, resJson)

}
