package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var CLIENT_ID string
var APP_SECRET string
var base = "https://api-m.sandbox.paypal.com"

type CreateOrder struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Links  []Links `json:"links"`
}

type Links struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// Work in Progress
func createOrder(w http.ResponseWriter, r *http.Request) {
	accessToken := generateAccessToken()
	body := []byte(`{
		"intent":"CAPTURE",
		"purchase_units":[
		   {
			  "amount":{
				 "currency_code":"USD",
				 "value":"100.00"
			  }
		   }
		]
	 }`)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	req, err := http.NewRequest("POST", base+"/v2/checkout/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println("request failed with status:", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func capturePayment(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["orderID"]
	accessToken := generateAccessToken()
	body := []byte(`{}`)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	req, err := http.NewRequest("POST", base+"/v2/checkout/orders/"+orderID+"/capture", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("An Error Occured:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println("request failed with status:", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		return
	}

	// copy response from external to frontend
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}

}

func generateAccessToken() string {

	auth := base64.StdEncoding.EncodeToString([]byte(CLIENT_ID + ":" + APP_SECRET))

	body := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", base+"/v1/oauth2/token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response into a map
	var jsonMap map[string]interface{}
	json.Unmarshal(bodyBytes, &jsonMap)
	access_token := fmt.Sprint(jsonMap["access_token"])

	return access_token
}

func main() {
	log.Println("Running on port: 8000")

	//Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	CLIENT_ID = os.Getenv("CLIENT_ID")
	APP_SECRET = os.Getenv("APP_SECRET")

	//Init Mux Router
	router := mux.NewRouter()
	router.HandleFunc("/api/orders", createOrder).Methods("POST")
	router.HandleFunc("/api/orders/{orderID}/capture", capturePayment).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
