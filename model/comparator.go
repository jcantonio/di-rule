package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ComparatorCondition defines a condition for comparing JSON values
type ComparatorCondition struct {
	Path     string // JSON key path (e.g., "age", "details.price")
	Operator string // Comparison operator (==, !=, >, <, >=, <=)
	Value    string // Value to compare against (as string)
}

// IsMet checks if the condition is met based on the JSON input
func (c *ComparatorCondition) IsMet(currentEntityJSON *string) (bool, error) {
	// Parse JSON into a map
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(*currentEntityJSON), &data); err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}

	// Extract the actual value from JSON
	actualValue, err := extractJSONValue(data, c.Path)
	if err != nil {
		return false, err
	}

	// Compare the extracted value with the expected value
	return compareValues(actualValue, c.Value, c.Operator)
}

// extractJSONValue retrieves a value from JSON using dot notation (e.g., "details.price")
func extractJSONValue(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	var value interface{} = data

	for _, key := range keys {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path: %s", path)
		}
		value, ok = m[key]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", path)
		}
	}

	return value, nil
}

// compareValues performs the comparison based on the operator
func compareValues(actual interface{}, expectedStr, operator string) (bool, error) {
	actualStr := fmt.Sprintf("%v", actual) // Convert actual value to string

	switch operator {
	case "==":
		return actualStr == expectedStr, nil
	case "!=":
		return actualStr != expectedStr, nil
	case ">", "<", ">=", "<=":
		actualFloat, err := toFloat(actual)
		if err != nil {
			return false, errors.New("actual value is not a number")
		}
		expectedFloat, err := strconv.ParseFloat(expectedStr, 64)
		if err != nil {
			return false, errors.New("expected value is not a number")
		}
		return compareNumeric(actualFloat, expectedFloat, operator), nil
	default:
		return false, errors.New("unsupported operator")
	}
}

// toFloat converts an interface{} to float64 if possible
func toFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, errors.New("invalid numeric type")
	}
}

// compareNumeric compares two float64 values based on the operator
func compareNumeric(actual, expected float64, operator string) bool {
	switch operator {
	case ">":
		return actual > expected
	case "<":
		return actual < expected
	case ">=":
		return actual >= expected
	case "<=":
		return actual <= expected
	default:
		return false
	}
}
