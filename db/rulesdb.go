package db

import (
	"os"
	"regexp"

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
func CreateRule(id string, doc map[string]interface{}) (string, error) {
	err := diRuleDB.Set(id, doc)
	rev := doc["_rev"].(string)
	return rev, err
}

/*
UpdateRule updates a rule in db
*/
func UpdateRule(id string, doc map[string]interface{}) (string, error) {
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
