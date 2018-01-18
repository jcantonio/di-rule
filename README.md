# di-rule
A Go Rules Engine

Project from a Golang learner. Please give feedbacks. 

Create Rules based on Conditions (Condition) which can be
Logical Conditions (LogicalCondition) or Value Comparator Conditions (ComparatorCondition)

Value Comparator Conditions are Value comparators with 
	Path     string path to json element
	Operator string 
	Value    string¡bool¡numeric 

Input json entity

```bash
go get -u github.com/jcantonio/di-rule
```

# Roadmap
string comparator: equal        - DONE
logical comparators: And        - DONE
logical comparators: Or         - DONE
string comparator: regex        - DONE
nil comparator
number comparator: equal
number comparator: greater
number comparator: less
REST : Add Rule