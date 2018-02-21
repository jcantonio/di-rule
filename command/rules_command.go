package command

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jcantonio/di-rule/converter"
	"github.com/jcantonio/di-rule/db"
	"github.com/jcantonio/di-rule/model"
)

var rulesInMem map[string](map[string]model.Rule) = make(map[string](map[string]model.Rule))

type ExecuteRule func(rule *model.Rule, entityJSON *string) error

type Action interface {
	Execute(rule *model.Rule, entityJSON *string) error
}

type ExecuteActions struct {
}

type ExecuteGatherRules struct {
	rules []model.Rule
}

func (exe *ExecuteActions) Execute(rule *model.Rule, entityJSON *string) error {
	println("PASSED", rule.Name)
	return nil
}
func (exe *ExecuteGatherRules) Execute(rule *model.Rule, entityJSON *string) error {
	exe.rules = append(exe.rules, *rule)
	return nil
}
func ProcessRules(entityType *string, entityJSON *string) ([]interface{}, error) {
	var rulesMet []interface{}
	if rulesInMem == nil {
		LoadRulesInMem()
	}
	rulesForEntityType := rulesInMem[*entityType]
	for _, rule := range rulesForEntityType {
		conditionIsMet, err := rule.Condition.IsMet(entityJSON)
		if err != nil {
			return nil, err
		}
		if conditionIsMet {
			ruleMet := map[string]interface{}{
				"id":      rule.ID,
				"name":    rule.Name,
				"actions": rule.Actions,
			}
			rulesMet = append(rulesMet, ruleMet)
		}
	}
	return rulesMet, nil
}

/*
InitDatabase initialise the db. create it if does not exist and load in memory Rules
*/
func InitDatabase(url string, dbname string) {
	exitCode := 1
	db.InitServer(url, exitCode)
	db.InitDB(dbname, exitCode)
	//create view if does not exist
	viewId := "_design/rules"
	view, _ := db.GetDocument(viewId)
	if view == nil {
		view = map[string]interface{}{
			"language": "javascript",
			"views": map[string]interface{}{
				"all": map[string]interface{}{
					"map": "function(doc) { emit(doc._id, doc)}",
				},
			},
		}
		db.CreateDocument(viewId, view)
	}
}

func LoadRulesInMem() error {
	rules, _, _, _, _, _, _, _, err := GetRules(nil, 10000, 1)

	for _, rule := range rules {
		addRuleInMem(&rule)
	}

	if err != nil {
		return err
	}
	return nil
}
func addRuleInMem(rule *model.Rule) {
	if rulesInMem == nil {
		rulesInMem = make(map[string](map[string]model.Rule))
	}
	rulesPerEntity := rulesInMem[rule.Entity]
	if rulesPerEntity == nil {
		rulesPerEntity = make(map[string]model.Rule)
		rulesInMem[rule.Entity] = rulesPerEntity
	}
	rulesPerEntity[rule.ID] = *rule
}

func removeRuleFromMem(entity, ruleId string) {
	rulesPerEntity := rulesInMem[entity]
	if rulesPerEntity != nil {
		delete(rulesPerEntity, ruleId)
	}
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
	doc["_id"] = id
	//create and Validate Rule
	rule, err := converter.GetRule(doc)

	if err != nil {
		return "", "", err
	}
	// Store
	ver, err := db.CreateDocument(id, doc)

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
	ver, err := db.UpdateDocument(id, doc)

	if err != nil {
		return "", err
	}

	// Update Cache
	addRuleInMem(&rule)

	return ver, err
}
func DeleteRule(id string) error {
	ruleMap, err := db.GetDocument(id)
	if err != nil {
		return err
	}
	db.DeleteDocument(id)
	entity := ruleMap["entity"]
	removeRuleFromMem(entity.(string), id)
	return nil
}
func GetRulesAsMaps(sorts []string, pageSize int, page int) ([]map[string]interface{}, int, int, int, int, int, int, int, error) {

	var selfPage, firstPage, prevPage, nextPage, lastPage, totalPages int

	limit := pageSize
	skip := pageSize * (page - 1)
	rulesAsMaps, _, total, err := db.GetDocuments(sorts, limit, skip)

	totalPages = (total + pageSize - 1) / pageSize
	selfPage = page

	if totalPages == 0 {
		firstPage, prevPage, nextPage, lastPage = 0, 0, 0, 0
	} else if page == 1 {
		if page >= totalPages {
			firstPage, prevPage, nextPage, lastPage = 1, 1, 1, totalPages
		} else {
			firstPage, prevPage, nextPage, lastPage = 1, 1, 2, totalPages
		}
	} else if page >= totalPages {
		firstPage, prevPage, nextPage, lastPage = 1, page-1, page, totalPages
	}
	return rulesAsMaps, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, err
}
func GetRules(sorts []string, pageSize int, page int) ([]model.Rule, int, int, int, int, int, int, int, error) {

	var rulesMap []map[string]interface{}
	var err error

	var selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total int
	rulesMap, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, err = GetRulesAsMaps(sorts, pageSize, page)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, 0, 0, err
	}
	rules := []model.Rule{}
	for _, ruleMap := range rulesMap {
		rule, err := converter.GetRule(ruleMap)
		if err != nil {
			return nil, 0, 0, 0, 0, 0, 0, 0, err
		}
		rules = append(rules, rule)
	}
	return rules, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, nil
}
func GetRulesAsJSON(sorts []string, pageSize int, page int) ([]byte, int, int, int, int, int, int, int, error) {

	var rulesMap []map[string]interface{}
	var err error
	var selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total int
	rulesMap, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, err = GetRulesAsMaps(sorts, pageSize, page)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, 0, 0, err
	}

	var jsonResult []byte
	jsonResult, err = json.Marshal(rulesMap)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, 0, 0, err
	}
	return jsonResult, selfPage, firstPage, prevPage, nextPage, lastPage, totalPages, total, nil
}
