package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/jcantonio/di-rule/command"
	"github.com/jcantonio/di-rule/db"

	"github.com/gorilla/mux"
)

var mu sync.Mutex
var count int

/*
Init runs server to handle requests
*/
func Init(port int) {

	router := mux.NewRouter()
	//GET
	//http.HandleFunc("/rules", handler)
	//http.HandleFunc("/", handler)

	router.HandleFunc("/processes/run", ProcessRules).Methods("POST")
	router.HandleFunc("/rules", GetRules).Methods("GET")
	router.HandleFunc("/rules", CreateRule).Methods("POST")
	router.HandleFunc("/rules/{id}", GetRule).Methods("GET")
	router.HandleFunc("/rules/{id}", UpdateRule).Methods("PUT")
	router.HandleFunc("/rules/{id}", DeleteRule).Methods("DELETE")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func ProcessRules(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//entityType := params["entityType"]

	entityTypes, ok := r.URL.Query()["entityType"]

	if !ok || len(entityTypes) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Url Param 'entityType' is missing"))
		return
	}
	// Query()["entityType"] will return an array of items,
	// we only want the single item.
	entityType := entityTypes[0]

	var body []byte
	var err error
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	entityJSON := string(body)
	command.ProcessRules(&entityType, &entityJSON, command.ExecuteActions)
}
func GetRules(w http.ResponseWriter, r *http.Request) {
	/*
		json.NewEncoder(w).Encode(people)
	*/
	/*
		params := mux.Vars(r)
		    for _, item := range people {
		        if item.ID == params["id"] {
		            json.NewEncoder(w).Encode(item)
		            return
		        }
		    }
		    json.NewEncoder(w).Encode(&Person{})
	*/

	//fmt.Fprintf(w, "GetRules %s", "to")

	selector := `_id > nil`

	jsonResult, err := command.GetRulesAsJSON(nil, selector, nil, nil, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(jsonResult)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
func GetRule(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.String())
	ruleMap, err := db.GetRule(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json, _ := json.Marshal(ruleMap)
	w.Write(json)
}

func CreateRule(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	var rev string
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var id string
	id, rev, err = db.CreateRule(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	respBody := fmt.Sprintf(`{"_id": "%s","_rev":"%s"}`, id, rev)
	w.Write([]byte(respBody))
}
func UpdateRule(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	var rev1 string
	var rev2 string

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	id := path.Base(r.URL.String())

	revResult := gjson.Get(string(body), "_rev")
	rev1 = revResult.String()
	rev2, err = db.UpdateRule(id, rev1, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := fmt.Sprintf(`{"_id": "%s","_rev":"%s"}`, id, rev2)
	w.Write([]byte(respBody))
}
func DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.String())
	err := db.DeleteRule(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.path = %q\n", r.URL.Path)
}
