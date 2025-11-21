package requests

import (
	"io"
	"log"
	"net/http"
	"bytes"
	"time"
	"encoding/json"
)
 
func GetTransactionStatus (url string, apiKey string) {

	retries := 1
for {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "Token "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
        if err != nil {
                log.Fatalln(err)
        }
	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
        if err != nil {
                log.Fatalln(err)
        }

	//Convert the body to type string
/*	sb := string(body)
	log.Printf("Transaction Status: %v", sb)*/

	//Extract and return Status
        var t map[string]interface{}
        json.Unmarshal(body, &t)
        status, ok := t["status"].(string)
	if ok {

	switch status {
	case "SUCCESSFUL":
		log.Println("\n Transaction Status: Payment completed successfully")
		return

	case "FAILED":
		log.Println("\n Transaction Status: Payment failed")
		return
	}
	}else {
		log.Println("status does not exist")
	}

	if retries<5{
                log.Println("Payment still pending... checking again in 5 seconds")
                time.Sleep(5 * time.Second)
		retries ++
		continue
	} else {
		log.Println("\n Tansaction Status: Payment still pending")
		return
	}
}
}

func PaymentRequest (url string, apiKey string, content map[string]string) string {
	//Encoded the data
	postBody, _ := json.Marshal(content)
	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make a request
	req, err := http.NewRequest(http.MethodPost, url, responseBody)
        if err != nil {
                log.Fatalln(err)
        }

	req.Header.Add("Authorization", "Token "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer response.Body.Close()

	//Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
/*	sb := string(body)
	log.Printf(sb)*/

	//Extract and return reference
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if tid, ok := result["reference"].(string); ok {
		return tid
	}

	return ""
}

