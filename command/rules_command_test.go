package command

import (
	"reflect"
	"testing"

	dirule "github.com/jcantonio/di-rule"
	"github.com/jcantonio/di-rule/db"
)

func TestGetRules(t *testing.T) {

	db.InitDatabase("http://localhost:5984", "di-rule")

	selector := `entity == "CUSTOMER"`

	selector = `_id > nil`

	t.Run("TEST1", func(t *testing.T) {
		rules, err := GetRules(nil, selector, nil, nil, nil, nil)
		println(rules)
		if err != nil {
			t.Errorf("GetRules() error = %v", err)
			return
		}
	})

}

func TestGetRuleAsObject(t *testing.T) {
	type args struct {
		ruleMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    dirule.Rule
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRule(tt.args.ruleMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConditionAsObject(t *testing.T) {
	type args struct {
		conditionMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    dirule.Condition
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCondition(tt.args.conditionMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadRules(t *testing.T) {
	db.InitDatabase("http://localhost:5984", "di-rule")
	if err := LoadRules(); err != nil {
		t.Errorf("LoadRules() error = %v", err)
	}
}
