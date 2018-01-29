package db

import (
	"testing"
)

var jsonStr1 = `{
	"name": "R1",
    "entity": "CUSTOMER",	
	"description": "R1 IS US or FR",
	"action": [{
		"name": "DoThis"
	}],
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
}`

var jsonStr2 = `{
	"name": "R1",
    "entity": "CUSTOMER",
	"description": "R1 IS GB or FR",
	"action": [{
		"name": "DoThis"
	}],
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
				"value": "GB"
			}
		]
	}
}`

func TestCRUD(t *testing.T) {
	InitDatabase("http://localhost:5984", "di-rule")
	id, rev, err := CreateRule([]byte(jsonStr1))
	if err != nil {
		t.Error(err)
		return
	}
	rev, err = UpdateRule(id, rev, []byte(jsonStr2))
	if err != nil {
		t.Error(err)
		return
	}
	var rule map[string]interface{}
	rule, err = GetRule(id)
	if err != nil {
		t.Error(err)
		return
	}
	println(rule)
	/*
		{
		    "selector": {
		        "year": {"$gt": 2010}
		    },
		    "fields": ["_id", "_rev", "year", "title"],
		    "sort": [{"year": "asc"}],
		    "limit": 2,
		    "skip": 0
		}
	*/

	selector := `entity == "CUSTOMER"`
	var rules []map[string]interface{}
	rules, err = GetRules(nil, selector, nil, nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for _, rule := range rules {
		t.Log(rule)
	}

	DeleteRule(id)
}
