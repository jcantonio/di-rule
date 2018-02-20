package db

import (
	"testing"

	"github.com/google/uuid"
)

var jsonStr1 = `{
	"name": "R1",
    "entity": "CUSTOMER",	
	"description": "R1 IS US or FR",
	"action": [{
		"name": "DoThis"
	}],
	"condition": {
		"op": "or",
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
		"op": "or",
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
	InitDB("di-rule", 1)
	id := uuid.New().String()

	doc := map[string]interface{}{
		"name": "nameA",
	}

	rev, err := CreateDocument(id, doc)
	if err != nil {
		t.Error(err)
		return
	}

	doc = map[string]interface{}{
		"_id":  id,
		"name": "nameB",
	}
	rev, err = UpdateDocument(rev, doc)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetDocument(id)
	if err != nil {
		t.Error(err)
		return
	}

	var rules []map[string]interface{}
	rules, _, _, err = GetDocuments(nil, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	for _, rule := range rules {
		t.Log(rule)
	}

	DeleteDocument(id)
}
