package command

import (
	"testing"
)

func TestGetRules(t *testing.T) {

	InitDatabase("http://localhost:5984", "di-rule")

	t.Run("TEST1", func(t *testing.T) {
		rules, _, _, _, _, _, _, _, err := GetRules(nil, 1, 1)
		println(rules)
		if err != nil {
			t.Errorf("GetRules() error = %v", err)
			return
		}
	})
}

func TestLoadRules(t *testing.T) {
	InitDatabase("http://localhost:5984", "di-rule")
	if err := LoadRulesInMem(); err != nil {
		t.Errorf("LoadRules() error = %v", err)
	}
}

func TestRuleCRUD(t *testing.T) {
	InitDatabase("http://localhost:5984", "di-rule")
	ruleJSON := `{
		"action": [
		  {
			"name": "DoThis"
		  }
		],
		"condition": {
		  "conditions": [
			{
			  "op": "=",
			  "path": "shipTo.countryCode",
			  "value": "FR"
			},
			{
			  "op": "=",
			  "path": "shipTo.countryCode",
			  "value": "US"
			}
		  ],
		  "op": "or"
		},
		"description": "R3 IS US or FR",
		"entity": "CUSTOMER",
		"name": "R3"
	  }`

	id, rev, err := CreateRule([]byte(ruleJSON))
	if err != nil {
		t.Error(err)
		return
	}
	rev, err = UpdateRule(id, rev, []byte(ruleJSON))
	if err != nil {
		t.Error(err)
		return
	}
	err = DeleteRule(id)
	if err != nil {
		t.Error(err)
		return
	}
}
