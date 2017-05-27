/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (

	"errors"
	"fmt"
	"math/rand"
	"time"
	"encoding/json"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/crypto/primitives"
)

// MetaTagger is simple chaincode implementing a basic Asset Management system
// with access control enforcement at chaincode level.
// Look here for more information on how to implement access control at chaincode level:
// https://github.com/hyperledger/fabric/blob/master/docs/tech/application-ACL.md
// An asset is simply represented by a string.
type MetaTagger struct {
}

// Init method will be called during deployment.
// The deploy transaction metadata is supposed to contain the administrator cert

func (t *MetaTagger) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	// Create ownership table
	err := stub.CreateTable("MetaTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Index", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Doc", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Tag", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Note", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating MetaTable table.")
	}

	fmt.Printf("Init Chaincode...done")

	return nil, nil
}

//
//
//


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
		rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

//
//
//
func (t *MetaTagger) create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("create...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	index := RandStringRunes(32)
	doc := args[0]
	tag := args[1]
	note := args[2]


	ok, err := stub.InsertRow("MetaTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: index}},
			&shim.Column{Value: &shim.Column_String_{String_: doc}},
			&shim.Column{Value: &shim.Column_String_{String_: tag}},
			&shim.Column{Value: &shim.Column_String_{String_: note}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Meta was already created.")
	}

	fmt.Printf("create...done!")

	return nil, err
}



// Invoke will be called for every transaction.
// Supported functions are the following:
// "create(asset, owner)": to create ownership of assets. An asset can be owned by a single entity.
// Only an administrator can call this function.
// "transfer(asset, newOwner)": to transfer the ownership of an asset. Only the owner of the specific
// asset can call this function.
// An asset is any string to identify it. An owner is representated by one of his ECert/TCert.
func (t *MetaTagger) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "create" {
		// create ownership
		return t.create(stub, args)
	}
	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
// Supported functions are the following:
// "query(asset)": returns the owner of the asset.
// Anyone can invoke this function.
func (t *MetaTagger) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Printf("Query [%s]", function)

	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting 'query' but found '" + function + "'")
	}

	var err error

	if len(args) != 1 {
		fmt.Printf("Incorrect number of arguments. Expecting name of an user to query")
		return nil, errors.New("Incorrect number of arguments. Expecting name of an user to query")
	}

	var columns []shim.Column

	col1Val := args[0]
	col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
	columns = append(columns, col1)

	rowChannel, err := stub.GetRows("MetaTable", columns)
	if err != nil {
		return nil, fmt.Errorf("get Rows failed. %s", err)
	}
	var string valueTest = ""
	var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				valueTest = valueTest+"A"
				rowChannel = nil
			} else {
				rows = append(rows, row)
				valueTest = valueTest+"B"
			}
		}
		if rowChannel == nil {
			valueTest+"C"
			break
		}
	}

	jsonRows, err := json.Marshal(rows)
	if err != nil {
		return nil, fmt.Errorf("getRowsTableFour operation failed. Error marshaling JSON: %s", err)
	}

	return []byte(valueTest), nil
}

func main() {
	err := shim.Start(new(MetaTagger))
	if err != nil {
		fmt.Printf("Error starting MetaTagger: %s", err)
	}
}
