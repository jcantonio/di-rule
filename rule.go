package dirule

import (
	"errors"
	"regexp"

	"github.com/tidwall/gjson"
)

const (
	// And is the Logical Operator And
	And = "and"
	// Or is the Logical Operator Or
	Or = "or"
	//Equal is the Equal value comparator
	Equal = "="
	//Regexp is the Equal value comparator
	Regexp = "regexp"
	// Greater is the Greater value comparator
	Greater = ">"
	// Lesser is the Lesser value comparator
	Lesser = "<"
	// Nil is the Greater value comparator
	Nil = "nil"
)

/*
Condition needs to be implemented
*/
type Condition interface {
	IsMet(currentEntityJSON *string) (bool, error)
	//	IsMetWhenChanged(currentEntityJSON, previousEntityJSON string) (bool, error)
}

/*
Rule compare
*/
type Rule struct {
	Name      string
	Entity    string
	Condition Condition
	Actions   []Action
}

/*
Action TBD
*/
type Action struct {
	Name string
}

/*
StringComparatorCondition compares string
*/
type StringComparatorCondition struct {
	Path     string
	Operator string
	Value    string
}

/*
IsMet verifies whether the condition is met
*/
func (comparator *StringComparatorCondition) IsMet(currentEntityJSON *string) (bool, error) {
	result := gjson.Get(*currentEntityJSON, comparator.Path)
	switch result.Type {
	case gjson.Null:
		return false, errors.New("Not implemented")
	case gjson.True:
		return false, errors.New("Not implemented")
	case gjson.False:
		return false, errors.New("Not implemented")
	case gjson.JSON:
		return false, errors.New("Not implemented")
	case gjson.String:
		return compareString(comparator, result.Str)
	case gjson.Number:
		return false, errors.New("Not implemented")
	}
	return false, errors.New("Unidentified type")
}

func compareString(comparator *StringComparatorCondition, value string) (bool, error) {
	switch comparator.Operator {
	case Equal:
		if comparator.Value == value {
			return true, nil
		}
		return false, nil
	case Regexp:
		return regexp.MatchString(comparator.Value, value)
	}
	return false, errors.New("Unidentified Comp")
}

/*
LogicalCondition
*/
type LogicalCondition struct {
	Operator   string
	Conditions []Condition
}

/*
Add func implements interface, byt adding condition
*/
func (comparator *LogicalCondition) Add(condition Condition) {
	comparator.Conditions = append(comparator.Conditions, condition)
}

/*
Remove func implements interface, by Removing condition
*/
func (comparator *LogicalCondition) Remove(condition Condition) {
	comparator.Conditions = append(comparator.Conditions, condition)
}

/*
IsMet func implements IsMet
*/
func (comparator *LogicalCondition) IsMet(currentEntityJSON *string) (bool, error) {
	switch comparator.Operator {
	case Or:
		for _, condition := range comparator.Conditions {
			result, err := condition.IsMet(currentEntityJSON)
			if err != nil {
				return result, err
			}
			if result == true {
				return true, nil
			}
		}
		return false, nil
	case And:
		for _, condition := range comparator.Conditions {
			result, err := condition.IsMet(currentEntityJSON)
			if err != nil {
				return result, err
			}
			if result == false {
				return false, nil
			}
		}
		return true, nil
	}
	return false, errors.New("Unidentified operator")
}
