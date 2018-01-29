package db

import (
	"os"
	"regexp"

	"github.com/tidwall/gjson"

	"github.com/google/uuid"
	"github.com/leesper/couchdb-golang"
)

var (
	serverDB *couchdb.Server
	diRuleDB *couchdb.Database
)

/*
InitDatabase initialise the db. create it if does not exist and load in memory Rules
*/
func InitDatabase(url string, dbname string) {
	exitCode := 1
	initServer(url, exitCode)
	initDB(dbname, exitCode)
}

func initServer(url string, exitCode int) {
	var err error
	serverDB, err = couchdb.NewServer(url)
	if err != nil {
		os.Exit(exitCode)
	}
	serverDB.Version()
}

func initDB(name string, exitCode int) {
	var err error
	diRuleDB, err = serverDB.Get(name)
	if err != nil {
		regResp, _ := regexp.MatchString("404", err.Error())
		if regResp {
			diRuleDB, err = serverDB.Create(name)
			if err != nil {
				os.Exit(exitCode)
			}
		} else {
			os.Exit(exitCode)
		}
	}
}

/*
CreateRule stores a rule in db
*/
func CreateRule(json []byte) (string, string, error) {
	result := gjson.ParseBytes(json)
	buffer := Buffer{}
	buffer.AddResult(result)
	doc := buffer.records //map[string]interface{}
	id := uuid.New().String()
	err := diRuleDB.Set(id, doc)
	rev := doc["_rev"].(string)
	return id, rev, err
}

/*
UpdateRule updates a rule in db
*/
func UpdateRule(id string, rev1 string, json []byte) (string, error) {
	result := gjson.ParseBytes(json)
	buffer := Buffer{}
	buffer.AddResult(result)
	buffer.records["_rev"] = rev1
	doc := buffer.records
	err := diRuleDB.Set(id, doc)
	rev2 := doc["_rev"].(string)
	return rev2, err
}

/*
DeleteRule deletes a rule from db
*/
func DeleteRule(id string) error {
	return diRuleDB.Delete(id)
}

/*
GetRule get rule from db
*/
func GetRules(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]map[string]interface{}, error) {
	return diRuleDB.Query(fields, selector, sorts, limit, skip, index)
}

/*
GetRules get rules from db
*/
func GetRule(id string) (map[string]interface{}, error) {
	return diRuleDB.Get(id, nil)
}

// ProcessRecord adds a message to the buffer.
func (buffer *Buffer) AddResult(result gjson.Result) {
	if buffer.records == nil {
		buffer.records = make(map[string]interface{})
	}
	result.ForEach(func(key, value gjson.Result) bool {

		//ADD RECORD
		attrName := key.Str
		switch value.Type {
		case gjson.Null:
		case gjson.True:
			buffer.records[attrName] = value.Bool
		case gjson.False:
			buffer.records[attrName] = value.Bool
		case gjson.JSON:
			if value.IsArray() {
				// Array
				s := len(value.Array())
				subElements := make([]map[string]interface{}, s)

				for index, subResult := range value.Array() {
					subBuffer := Buffer{}
					subBuffer.AddResult(subResult)
					subElements[index] = subBuffer.records
				}
				buffer.records[attrName] = subElements
			} else {
				// Simple
				subBuffer := Buffer{}
				subBuffer.AddResult(value)
				buffer.records[attrName] = subBuffer.records
			}

		case gjson.String:
			buffer.records[attrName] = value.Str
		case gjson.Number:
			buffer.records[attrName] = value.Num
		}

		return true // keep iterating
	})
}

type Buffer struct {
	records map[string]interface{}
}
