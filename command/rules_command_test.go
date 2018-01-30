package command

import (
	"testing"

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

func TestLoadRules(t *testing.T) {
	db.InitDatabase("http://localhost:5984", "di-rule")
	if err := LoadRulesInMem(); err != nil {
		t.Errorf("LoadRules() error = %v", err)
	}
}
