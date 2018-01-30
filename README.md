# di-rule
A Go Rules Engine

## Work in Progress, NOT Ready to be used yet

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

## Roadmap
* string comparator: equal        	- DONE
* logical comparators: And        	- DONE
* logical comparators: Or         	- DONE
* string comparator: equal        	- DONE
* string comparator: regex        	- DONE
* number comparator: equal        	- DONE
* number comparator: greater        - DONE
* number comparator: less        	- DONE
* nil comparator
* date comparator: equal         
* date comparator: greater        
* date comparator: less        	
* compare when changes
* REST : Add Rule in CouchDB	  	- PART 
* REST : Execute Rule				- PART