package db

import (
	"fmt"
	"math"
	"os"
	"regexp"

	couchdb "github.com/jcantonio/couchdb-golang"
)

var (
	serverDB *couchdb.Server
	diRuleDB *couchdb.Database
)

func InitServer(url string, exitCode int) {
	var err error
	serverDB, err = couchdb.NewServer(url)
	if err != nil {
		os.Exit(exitCode)
	}
}

func InitDB(name string, exitCode int) {
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
CreateDocument stores a rule in db
*/
func CreateDocument(id string, doc map[string]interface{}) (string, error) {
	err := diRuleDB.Set(id, doc)
	rev := doc["_rev"].(string)
	return rev, err
}

/*
UpdateDocument updates a rule in db
*/
func UpdateDocument(id string, doc map[string]interface{}) (string, error) {
	err := diRuleDB.Set(id, doc)
	rev2 := doc["_rev"].(string)
	return rev2, err
}

/*
DeleteDocument deletes a rule from db
*/
func DeleteDocument(id string) error {
	return diRuleDB.Delete(id)
}

/*
GetDocuments get rule from db
*/
func GetDocuments(sorts []string, limit, skip int) ([]map[string]interface{}, int, int, error) {

	option := map[string]interface{}{"limit": limit, "skip": skip} //"descending": true,

	results, err := diRuleDB.View("rules/all", nil, option)
	if err != nil {
		println(err.Error())
	}

	// couchdb-golang only fetches data when calling .Rows()
	rows, err := results.Rows()
	if err != nil {
		println(err.Error())
	}

	if rows == nil {
		println(err.Error())
	}

	totalRows, _ := results.TotalRows()
	if totalRows == -1 {
		println(err.Error())

	}

	offset, _ := results.Offset()
	if offset == -1 {
		println(err.Error())
	}
	rules := []map[string]interface{}{}

	for _, row := range rows {
		rules = append(rules, row.Val.(map[string]interface{}))
	}

	return rules, offset, totalRows, nil //diRuleDB.Query(fields, selector, sorts, limit, skip, index)
}
func docFromNum(num int) map[string]interface{} {
	return map[string]interface{}{
		"_id": fmt.Sprintf("%d", num),
		"num": int(num / 2),
	}
}

func docFromRow(row couchdb.Row) map[string]interface{} {
	return map[string]interface{}{
		"_id": row.ID,
		"num": int(row.Key.(float64)),
	}
}

func iterateSlice(begin, end, incr int) []int {
	s := []int{}
	if begin <= end {
		for i := begin; i < end; i += incr {
			s = append(s, i)
		}
	} else {
		for i := begin; i > end; i += incr {
			s = append(s, i)
		}
	}
	return s
}
func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

/*
GetDocument get rules from db
*/
func GetDocument(id string) (map[string]interface{}, error) {
	return diRuleDB.Get(id, nil)
}
