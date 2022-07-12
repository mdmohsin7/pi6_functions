package api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
	"pi6_functions/api_utils"
)

type accResponse struct {
	Status         int               `json:"status,omitempty"`
	Addresses      map[string]string `json:"addresses,omitempty"`
	AddressesSaved bool              `json:"addresses_saved,omitempty"`
	ResponseText   string            `json:"responseText,omitempty"`
}

func saveAccounts(uid string, addresses map[string]string) (int, error) {
	// Use a service account

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
	if err != nil {
		return 0, err
	}
	res := 500
	var resErr error = nil
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
			_, err2 := client.Collection("users").Doc(uid).Create(ctx, map[string]interface{}{"addresses": addresses})
			if err2 != nil {
				resErr = err
				res = 500
			} else {
				resErr = nil
				res = 200
			}
		}
		log.Printf("Failed saving addresses: %v", err)
	} else {
		log.Println("Document does exist")
		_, err := client.Collection("users").Doc(uid).Update(ctx, []firestore.Update{{Path: "addresses", Value: addresses}})
		if err != nil {
			log.Println(err)
			resErr = err
			res = 500
		} else {
			resErr = nil
			res = 200
		}
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)
	return res, resErr
}

func FindAccounts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("addr")
	uid := r.URL.Query().Get("uid")
	var finalRes *accResponse
	accounts := make(map[string]string)
	_, b64, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return
	}
	for k, v := range api_utils.Prefixes {
		addr, e := bech32.ConvertAndEncode(v, b64)
		if e != nil {
			log.Println(k, e)
		}
		accounts[v] = addr
	}
	if err != nil {
		return
	}

	result, err := saveAccounts(uid, accounts)
	if err != nil {
		return
	}
	if result == 200 {
		finalRes = &accResponse{
			Addresses:      accounts,
			Status:         200,
			AddressesSaved: true,
			ResponseText:   "success",
		}
	} else {
		finalRes = &accResponse{
			Addresses:      accounts,
			Status:         200,
			AddressesSaved: false,
			ResponseText:   "success",
		}
	}
	resData, _ := json.Marshal(finalRes)
	_, _ = fmt.Fprintf(w, string(resData))
}
