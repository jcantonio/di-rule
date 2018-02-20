package model

import (
	"errors"
	"fmt"
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
	Less = "<"
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
	ID        string
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

	if result.Type == gjson.String {
		switch comparator.Operator {
		case Equal:
			return comparator.Value == result.Str, nil
		case Regexp:
			return regexp.MatchString(comparator.Value, result.Str)
		}
		return false, fmt.Errorf("Unidentified Operator %s", comparator.Operator)
	}
	return false, fmt.Errorf("Unidentified Type %s", result.Type)
}

/*
NumberComparatorCondition compares number
*/
type NumberComparatorCondition struct {
	Path     string
	Operator string
	Value    float64
}

/*
IsMet verifies whether the condition is met
*/
func (comparator *NumberComparatorCondition) IsMet(currentEntityJSON *string) (bool, error) {
	result := gjson.Get(*currentEntityJSON, comparator.Path)

	if result.Type == gjson.Number {
		switch comparator.Operator {
		case Equal:
			return result.Num == comparator.Value, nil
		case Greater:
			return result.Num > comparator.Value, nil
		case Less:
			return result.Num < comparator.Value, nil
		}
		return false, fmt.Errorf("Unidentified Operator %s", comparator.Operator)
	}
	return false, fmt.Errorf("Unidentified Type %s", result.Type)
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
