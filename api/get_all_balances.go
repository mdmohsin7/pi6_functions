package api

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pi6_functions/api_utils"
	"strconv"
)

var balances []bal

type bal struct {
	Balances []api_utils.Balance `json:"balances"`
	Chain    string              `json:"chain"`
}

type rateRes struct {
	USD float64 `json:"USD"`
}

type allBalResponse struct {
	TotalBalInUSD float64 `json:"totalBalInUSD,omitempty"`
	ChainBalances []bal   `json:"chainBalances,omitempty"`
	Status        int     `json:"status,omitempty"`
	ResponseText  string  `json:"responseText,omitempty"`
}

func getAccountAddresses(uid string) (res int, err error, addresses interface{}) {
	res = 500
	var addrs interface{} = nil
	var resErr error = nil
	fc := api_utils.FirebaseConfig{
		Type:                    os.Getenv("FC_TYPE"),
		ProjectID:               os.Getenv("FC_PROJECTID"),
		PrivateKeyID:            os.Getenv("FC_PRIVATEKEYID"),
		PrivateKey:              os.Getenv("FC_PRIVATEKEY"),
		ClientEmail:             os.Getenv("FC_CLIENTEMAIL"),
		ClientID:                os.Getenv("FC_CLIENTID"),
		AuthUri:                 os.Getenv("FC_AUTHURI"),
		TokenUri:                os.Getenv("FC_TOKENURI"),
		AuthProviderX509CertUrl: os.Getenv("FC_APX509CU"),
		ClientX509CertUrl:       os.Getenv("FC_CX509CU"),
	}
	bytes, err := json.Marshal(fc)
	ctx := context.Background()
	sa := option.WithCredentialsJSON(bytes)
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("FC_PROJECTID")}, sa)
	if err != nil {
		resErr = err
		log.Fatalln(err)

	}

	client, err := app.Firestore(ctx)
	if err != nil {
		resErr = err
		log.Fatalln(err)
	}
	_, err = client.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Println("Document does not exist")

			resErr = err
			res = 500

		}
		log.Printf("Failed saving addresses: %v", err)
	} else {
		log.Println("Document does exist")
		docSnap, err := client.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			log.Println(err)
			resErr = err
			res = 500
		} else {
			addrs = docSnap.Data()["addresses"]
			resErr = nil
			res = 200
		}
	}
	return res, resErr, addrs
}

func AllBalancesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	uid := r.URL.Query().Get("uid")
	var finalRes *allBalResponse
	totalInUSD := 0.0
	balances = nil
	res, err, addresses := getAccountAddresses(uid)
	api := "https://min-api.cryptocompare.com/data/price?fsym="
	if err != nil {
		return
	}
	if res == 200 {
		for k, v := range addresses.(map[string]interface{}) {
			_, err := getBalance(k, fmt.Sprint(v))
			if err != nil {
				return
			}

		}
		for i := range balances {
			res, err := http.Get(api + api_utils.Tickers[balances[i].Chain] + "&tsyms=USD")
			fmt.Println(api_utils.Tickers[balances[i].Chain])
			if err != nil {
				return
			}
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			fmt.Println(resJson)
			var rateResponse rateRes
			err = json.Unmarshal([]byte(resJson), &rateResponse)
			fmt.Println(rateResponse)
			balanceInFloat, err := strconv.ParseFloat(balances[i].Balances[0].Amount, 64)
			if err != nil {
				return
			}
			totalInUSD += (balanceInFloat / 1000000) * rateResponse.USD
		}
		fmt.Println(totalInUSD)
		finalRes = &allBalResponse{
			TotalBalInUSD: totalInUSD,
			ChainBalances: balances,
			Status:        200,
			ResponseText:  "successful",
		}
	} else {
		finalRes = &allBalResponse{
			Status:       401,
			ResponseText: "account does not exist",
		}
	}
	resData, _ := json.Marshal(finalRes)
	_, _ = fmt.Fprintf(w, string(resData))
}

func getBalance(k string, v string) (resStatus int, resErr error) {
	var s = 0
	var resError error = nil
	fmt.Printf("chain: %s, addr: %s\n", k, v)
	var api string
	var wrongChain = false
	path := "/cosmos/bank/v1beta1/balances/"
	if k == "cre" {
		api = "https://mainnet.crescent.network:1317" + path + v
	} else if k == "cosmos" {
		api = "https://rest-cosmoshub.ecostake.com" + path + v
	} else if k == "akash" {
		api = "https://akash.c29r3.xyz:443/api" + path + v
	} else {
		wrongChain = true
	}
	if wrongChain == false {
		res, apiErr := http.Get(api)
		if apiErr != nil {
			s = res.StatusCode
			resError = apiErr
			fmt.Println(apiErr)
		}

		if res.StatusCode == 200 {
			resBody, _ := ioutil.ReadAll(res.Body)
			resJson := string(resBody)
			fmt.Println(resJson)
			var balanceFromChain api_utils.BalanceFromChain
			err := json.Unmarshal([]byte(resJson), &balanceFromChain)
			if err != nil {
				return
			}
			if len(balanceFromChain.Balances) > 0 {
				balances = append(balances, bal{
					Balances: balanceFromChain.Balances,
					Chain:    k,
				})
			}

			s = 200
			resError = nil

		} else {
			s = res.StatusCode
			resError = apiErr
		}
	} else {
		//finalRes = &balResponse{
		//	Status:       501,
		//	ResponseText: "The chain isn't supported or does not exist",
		//}
	}
	return s, resError
}
