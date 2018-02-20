package converter

import (
	"errors"

	"github.com/jcantonio/di-rule/model"
)

func GetRule(ruleMap map[string]interface{}) (model.Rule, error) {
	var rule model.Rule
	name := ruleMap["name"].(string)
	entity := ruleMap["entity"].(string)
	actions := []model.Action{}

	conditionMap := ruleMap["condition"]

	condition, err := getCondition(conditionMap.(map[string]interface{}))
	if err != nil {
		return rule, err
	}
	var id string
	idInterface := ruleMap["_id"]
	if idInterface != nil {
		id = idInterface.(string)
	}

	rule = model.Rule{
		ID:        id,
		Name:      name,
		Entity:    entity,
		Actions:   actions,
		Condition: condition,
	}

	return rule, nil
}

func getCondition(conditionMap map[string]interface{}) (model.Condition, error) {
	operation := conditionMap["op"]

	if operation == nil {
		return nil, errors.New("No op found")
	}

	// Logical Condition
	switch operation {
	case "or", "and":
		condition := &model.LogicalCondition{
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

	// Value Condition
	path := conditionMap["path"]

	if path == nil {
		return nil, errors.New("No path found")
	}

	value := conditionMap["value"]

	if value == nil {
		return nil, errors.New("No value found")
	}

	switch value.(type) {
	case int, float64:
		condition := &model.NumberComparatorCondition{
			Path:     path.(string),
			Operator: operation.(string),
			Value:    value.(float64),
		}
		return condition, nil
	case string:
		condition := &model.StringComparatorCondition{
			Path:     path.(string),
			Operator: operation.(string),
			Value:    value.(string),
		}
		return condition, nil
	}

	return nil, errors.New("Type not handled yet")
}
