package command

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	dirule "github.com/jcantonio/di-rule"
	"github.com/jcantonio/di-rule/converter"
	"github.com/jcantonio/di-rule/db"
)

var rulesInMem map[string](map[string]dirule.Rule)

type ExecuteRule func(rule *dirule.Rule, entityJSON *string) error

type Action interface {
	Execute(rule *dirule.Rule, entityJSON *string) error
}

type ExecuteActions struct {
}

type ExecuteGatherRules struct {
	rules []dirule.Rule
}

func (exe *ExecuteActions) Execute(rule *dirule.Rule, entityJSON *string) error {
	println("PASSED", rule.Name)

	return nil
}
func (exe *ExecuteGatherRules) Execute(rule *dirule.Rule, entityJSON *string) error {
	exe.rules = append(exe.rules, *rule)
	return nil
}
func ProcessRules(entityType *string, entityJSON *string, action Action) error {
	if rulesInMem == nil {
		LoadRulesInMem()
	}
	rulesForEntityType := rulesInMem[*entityType]
	for _, rule := range rulesForEntityType {
		conditionIsMet, err := rule.Condition.IsMet(entityJSON)
		if err != nil {
			return err
		}
		if conditionIsMet {
			action.Execute(&rule, entityJSON)
		}
	}
	return nil
}

func LoadRulesInMem() error {
	if rulesInMem == nil {
		rulesInMem = make(map[string](map[string]dirule.Rule))
	}
	selector := `_id > nil`
	rules, err := GetRules(nil, selector, nil, nil, nil, nil)

	for _, rule := range rules {
		addRuleInMem(&rule)
	}
	println(rules)

	if err != nil {
		return err
	}
	return nil
}
func addRuleInMem(rule *dirule.Rule) {
	rulesPerEntity := rulesInMem[rule.Entity]
	if rulesPerEntity == nil {
		rulesPerEntity = make(map[string]dirule.Rule)
		rulesInMem[rule.Entity] = rulesPerEntity
	}
	rulesPerEntity[rule.Name] = *rule
}

func getDoc(jsonDoc []byte) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := json.Unmarshal(jsonDoc, &result)
	if err != nil {
		return result, err
	}
	if _, ok := result["error"]; ok {
		reason := result["reason"].(string)
		return result, errors.New(reason)
	}
	return result, nil
}

func CreateRule(json []byte) (string, string, error) {
	doc, err := getDoc(json)
	if err != nil {
		return "", "", err
	}

	id := uuid.New().String()

	//create and Validate Rule
	rule, err := converter.GetRule(doc)

	if err != nil {
		return "", "", err
	}
	// Store
	ver, err := db.CreateRule(id, doc)

	if err != nil {
		return "", "", err
	}

	// Update Cache
	addRuleInMem(&rule)

	return id, ver, err
}

func UpdateRule(id string, rev1 string, json []byte) (string, error) {
	doc, err := getDoc(json)
	if err != nil {
		return "", err
	}

	doc["_rev"] = rev1

	//create and Validate Rule
	rule, err := converter.GetRule(doc)

	if err != nil {
		return "", err
	}
	// Store
	ver, err := db.UpdateRule(id, doc)

	if err != nil {
		return "", err
	}

	// Update Cache
	addRuleInMem(&rule)

	return ver, err
}

func GetRules(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]dirule.Rule, error) {

	var rulesMap []map[string]interface{}
	var err error

	rulesMap, err = db.GetRules(fields, selector, sorts, limit, skip, index)
	if err != nil {
		return nil, err
	}
	rules := []dirule.Rule{}
	for _, ruleMap := range rulesMap {
		rule, err := converter.GetRule(ruleMap)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
func GetRulesAsJSON(fields []string, selector string, sorts []string, limit, skip, index interface{}) ([]byte, error) {

	var rules []map[string]interface{}
	var err error

	rules, err = db.GetRules(fields, selector, sorts, limit, skip, index)
	if err != nil {
		return nil, err
	}

	var jsonResult []byte
	jsonResult, err = json.Marshal(rules)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}
