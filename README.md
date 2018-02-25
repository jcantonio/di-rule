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


### overview 
Di-Rule saves and keeps in memory rules created via its REST API.
Each Rule is composed of 
* name: name of the rule
* description: description of the rule
* entity: entity related to the rule
* condition: condition of application of the rule
* actions: actions to be executed is condition is met

Once the rules are loaded, use the run end-point to pass an entity and get all actions to be executed for it.

### Create Rule
POST http://localhost:8000/rules 
```
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


response
```
{"data":{"_id":%ruleID%,"_rev":%revision%},"status":"success"}
```
for example
```
{"data":{"_id":"2f8198bb-7b02-4c08-9c9f-9ba02680a416","_rev":"1-c7809e5a70f4b1fb313b21df966f1eff"},"status":"success"}
```

### Update Rule
PUT http://localhost:8000/rules/%ruleID% 

```
{
  	"_id": %ruleID%,
  	"_rev": %rev%,
    "name": "R3",
    "entity": "CUSTOMER",
    "description": "R3 checks if the shipment is to go to the US or FR",
    	"actions": [{
		"name": "DoThatUpdate"
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

### Delete Rule
DELETE http://localhost:8000/rules/%ruleID% 

### Run Rules
POST http://localhost:8000/processes/run?entity=%ENTITY%
e.g.
POST http://localhost:8000/processes/run?entity=CUSTOMER
```
{ "name"   : "John Smith",
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