package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
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
func Init(port uint) {

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
	entities, ok := r.URL.Query()["entity"]
	if !ok || len(entities) < 1 {
		WriteError(w, "Url Param 'entity' is missing", http.StatusBadRequest)
		return
	}
	// Query()["entity"] will return an array of items,
	// we only want the single item.
	entity := entities[0]

	var body []byte
	var err error
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	entityJSON := string(body)
	execActions := command.ExecuteActions{}
	command.ProcessRules(&entity, &entityJSON, &execActions)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func GetRules(w http.ResponseWriter, r *http.Request) {
	var page, pageSize int

	pages := r.URL.Query()["page"]
	if len(pages) > 0 {
		page, _ = strconv.Atoi(pages[0])
	} else {
		page = 1
	}
	pageSizes := r.URL.Query()["size"]
	if len(pages) > 0 {
		pageSize, _ = strconv.Atoi(pageSizes[0])
	} else {
		pageSize = 1
	}
	rulesMaps, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, err := command.GetRulesAsMaps(nil, pageSize, page)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteListSuccess(w, rulesMaps, "", selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, pageSize)
}
func GetRule(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.String())
	ruleMap, err := db.GetRule(id)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteSuccess(w, ruleMap)
}

func CreateRule(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	var rev string
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var id string
	id, rev, err = command.CreateRule(body)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"_id":  id,
		"_rev": rev,
	}
	WriteSuccess(w, data)
}
func UpdateRule(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	var rev1 string
	var rev2 string

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := path.Base(r.URL.String())

	revResult := gjson.Get(string(body), "_rev")
	rev1 = revResult.String()
	rev2, err = command.UpdateRule(id, rev1, body)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"_id":  id,
		"_rev": rev2,
	}
	WriteSuccess(w, data)
}
func DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.String())
	err := db.DeleteRule(id)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteSuccess(w, nil)
}

func WriteSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		// TODO Convert map to bytes
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
func WriteListSuccess(w http.ResponseWriter, data interface{}, urlWithFilter string, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, pageSize int) {

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
		"meta": map[string]interface{}{
			"total-pages": totalPages,
		},
		"links": map[string]interface{}{
			"self":  fmt.Sprintf("%s?page[number]=%d&page[size]=%d", urlWithFilter, selfPage, pageSize),
			"first": fmt.Sprintf("%s?page[number]=%d&page[size]=%d", urlWithFilter, firstPage, pageSize),
			"prev":  fmt.Sprintf("%s?page[number]=%d&page[size]=%d", urlWithFilter, prevPage, pageSize),
			"next":  fmt.Sprintf("%s?page[number]=%d&page[size]=%d", urlWithFilter, nextPage, pageSize),
			"last":  fmt.Sprintf("%s?page[number]=%d&page[size]=%d", urlWithFilter, lastPage, pageSize),
		},
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		// TODO Convert map to bytes
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}

func WriteFail(w http.ResponseWriter, failMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	response := map[string]interface{}{
		"status": "fail",
		"data":   failMessage,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(failMessage))
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
func WriteError(w http.ResponseWriter, errorMessage string, status int) {
	response := map[string]interface{}{
		"status":  "error",
		"message": errorMessage,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorMessage, status)
	}
	w.Header().Add("Content-Type", "application/json")
	http.Error(w, string(responseJSON), status)
}
