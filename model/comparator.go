package model

type ComparatorCondition struct {
	Path     string
	Operator string
	Value    string
}

func (c *ComparatorCondition) IsMet(currentEntityJSON *string) (bool, error) {
	return false, nil
}
