package dirule

import (
	"testing"
)

const entity1 = `{ "name"   : "John Smith",
	"sku"    : "20223",
	"price"  : 23.95,
	"shipTo" : { "name" : "Jane Smith",
				 "address" : "123 Maple Street",
				 "city" : "Pretendville",
				 "state" : "NY",
				 "zip"   : "12345",
				 "countryCode"   : "FR" },
	"billTo" : { "name" : "John Smith",
				 "address" : "123 Maple Street",
				 "city" : "Pretendville",
				 "state" : "NY",
				 "zip"   : "12345",
				 "countryCode"   : "FR" }
  }`

var countries = [...]string{"GB", "DE", "IT", "SE", "SI", "FI", "FR", "VN", "ES", "BE", "CN", "RO", "US"}

func TestRuleStringComparator(t *testing.T) {

	entityA := entity1

	checkComparatorRuleFR := &ComparatorCondition{
		Path:     "shipTo.countryCode",
		Operator: "eq",
		Value:    "FR"}

	result, _ := checkComparatorRuleFR.IsMet(&entityA)
	if !result {
		t.Errorf("Str match %t \n", result)
	}
	checkComparatorRuleUS := &ComparatorCondition{
		Path:     "shipTo.countryCode",
		Operator: "eq",
		Value:    "US"}

	result, _ = checkComparatorRuleUS.IsMet(&entityA)
	if result {
		t.Errorf("Str match %t \n", result)
	}
	checkComparatorRuleRegex := &ComparatorCondition{
		Path:     "shipTo.name",
		Operator: "regexp",
		Value:    "Smith.*"}

	result, _ = checkComparatorRuleRegex.IsMet(&entityA)
	if !result {
		t.Errorf("Str match %t \n", result)
	}
	checkComparatorRuleRegex = &ComparatorCondition{
		Path:     "shipTo.name",
		Operator: "regexp",
		Value:    "Trump.*"}

	result, _ = checkComparatorRuleRegex.IsMet(&entityA)
	if result {
		t.Errorf("Str match %t \n", result)
	}
}
func TestRuleEq(t *testing.T) {
	entityA := entity1
	checkComparatorRuleFR := &ComparatorCondition{
		Path:     "shipTo.countryCode",
		Operator: "eq",
		Value:    "FR"}
	checkComparatorRuleUS := &ComparatorCondition{
		Path:     "shipTo.countryCode",
		Operator: "eq",
		Value:    "US"}

	orCOndition := LogicalCondition{
		Operator:   "or",
		Conditions: []Condition{}}
	orCOndition.Add(checkComparatorRuleFR)
	orCOndition.Add(checkComparatorRuleUS)

	result, _ := orCOndition.IsMet(&entityA)

	if !result {
		t.Errorf("Or match %t \n", result)
	}

	andCOndition := LogicalCondition{
		Operator:   "and",
		Conditions: []Condition{}}
	andCOndition.Add(checkComparatorRuleFR)
	andCOndition.Add(checkComparatorRuleUS)

	result, _ = andCOndition.IsMet(&entityA)
	if result {
		t.Errorf("And match %t \n", result)
	}
}
