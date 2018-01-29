package api

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {

	var jsonStr = []byte(`{
		"key": "R1",
		"description": "R1 IS US or FR",
		"action": {
			"name": "DoThis"
		},
		"condition": {
			"op": "OR",
			"conditions": [
				{
					"path": "shipTo.countryCode",
					"op": "=",
					"value": "FR"
				},
				{
					"path": "shipTo.countryCode",
					"op": "=",
					"value": "US"
				}
			]
		}
	}`)
	urlAPI := "http://localhost:5984/rules"

	client := &http.Client{}

	r, _ := http.NewRequest("POST", urlAPI, bytes.NewBuffer(jsonStr)) // <-- URL-encoded payload
	//    r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Set("X-Custom-Header", "myvalue")
	r.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

}
