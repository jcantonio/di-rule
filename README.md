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

## Get di-rule
Get dependencies
```bash
go get -u github.com/tidwall/gjson
go get -u github.com/jcantonio/couchdb-golang
```
Get di-rule
```bash
go get -u github.com/jcantonio/di-rule
```
## REST API
### Start server
Create a file config.yml
```
db:
    name:     di-rule
    Address:  "http://localhost"
    port:     5984
server:
    port:     8000
```
Start server.

### Create Rule
```
POST http://localhost:8000/rules 
{
    "name": "R3",
    "entity": "CUSTOMER",
    "description": "R3 checks if the shipment is to go to the US or FR",
    	"actions": [{
		"name": "DoThat"
	}],
    "condition": {
        "op": "or",
        "conditions": [
            {
                "path": "shipTo.countryCode",
                "op": "=",
                "value": "FR"
            },
            {
                "path": "shipTo.countryCode",
                "op": "=",
                "value": "US"
            }
        ]
    }
}
```
## Golang API
To fill

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
* REST : Add Rule in CouchDB	  	- DONE 
* REST : Execute Rule				- DONE