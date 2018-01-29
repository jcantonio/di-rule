package command

import (
	"encoding/json"
	"errors"

	dirule "github.com/jcantonio/di-rule"
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
		LoadRules()
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

func LoadRules() error {
	if rulesInMem == nil {
		rulesInMem = make(map[string](map[string]dirule.Rule))
	}
	selector := `_id > nil`
	rules, err := GetRules(nil, selector, nil, nil, nil, nil)

	for _, rule := range rules {
		rulesPerEntity := rulesInMem[rule.Entity]
		if rulesPerEntity == nil {
			rulesPerEntity = make(map[string]dirule.Rule)
			rulesInMem[rule.Entity] = rulesPerEntity
		}
		rulesPerEntity[rule.Name] = rule
	}
	println(rules)

	if err != nil {
		return err
	}
	return nil
}

func CreateRule(json []byte) (string, string, error) {
	//Validate
	_, err := validateRule(json)
	if err != nil {
		return "", "", err
	}
	return db.CreateRule(json)
}

func UpdateRule(id string, rev1 string, json []byte) (string, error) {
	_, err := validateRule(json)
	if err != nil {
		return "", err
	}
	return db.UpdateRule(id, rev1, json)
}
func validateRule(json []byte) (bool, error) {
	return true, nil
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
		rule, err := GetRule(ruleMap)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
func GetRule(ruleMap map[string]interface{}) (dirule.Rule, error) {
	var rule dirule.Rule
	name := ruleMap["name"].(string)
	entity := ruleMap["entity"].(string)
	actions := []dirule.Action{}

	conditionMap := ruleMap["condition"]

	condition, err := getCondition(conditionMap.(map[string]interface{}))
	if err != nil {
		return rule, err
	}

	rule = dirule.Rule{
		Name:      name,
		Entity:    entity,
		Actions:   actions,
		Condition: condition,
	}

	return rule, nil
}

func getCondition(conditionMap map[string]interface{}) (dirule.Condition, error) {
	operation := conditionMap["op"]
	if operation == nil {
		return nil, errors.New("No op found")
	}
	switch operation {
	case "or", "and", "OR", "AND":
		condition := &dirule.LogicalCondition{
			Operator: operation.(string)}

		subconditions := conditionMap["conditions"].([]interface{})
		for _, subcondition := range subconditions {
			subconditionMap := subcondition.(map[string]interface{})
			subcondition, err := getCondition(subconditionMap)
			if err != nil {
				return nil, err
			}
			condition.Add(subcondition)
		}
		return condition, nil
	}

	path := conditionMap["path"]
	value := conditionMap["value"]

	/*
		switch v := value.(type) {
		case int:
			// v is an int here, so e.g. v + 1 is possible.
			fmt.Printf("Integer: %v", v)
		case float64:
			// v is a float64 here, so e.g. v + 1.0 is possible.
			fmt.Printf("Float64: %v", v)
		case string:
			// v is a string here, so e.g. v + " Yeah!" is possible.
			fmt.Printf("String: %v", v)
		default:
			// And here I'm feeling dumb. ;)
			fmt.Printf("I don't know, ask stackoverflow.")
		}
	*/

	condition := &dirule.ComparatorCondition{
		Path:     path.(string),
		Operator: operation.(string),
		Value:    value.(string),
	}
	return condition, nil
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
