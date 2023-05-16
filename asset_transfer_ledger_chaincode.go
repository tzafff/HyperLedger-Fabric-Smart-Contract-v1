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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"
	"strings"

	// "github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric-chaincode-go/shim"

	// pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")

	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "create" {
		// Creates a new account with an initial balance
		return t.create(stub, args)
	} else if function == "GetAllAsset" {
		// Retrieves the names of all assets
		return t.GetAllAsset(stub)
	} else if function == "ChangeUniversity" {
		return t.ChangeUniversity(stub,args)
	}
	

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Role\":\"" + A + "\",\"Personal Data\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

// Transaction creates a new account with an initial balance
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var role string
	var name string
	var surname string
	var university string
	var vlab string
	var result string

	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	role = args[0]

	// Check if account already exists
	accountBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountBytes != nil {
		return shim.Error("Account already exists")
	}
	name = args[1]
	surname = args[2]
	university = args[3]
	vlab = args[4]
	result = args[5]
	// Parse the balance argument as an integer

	data := []byte(name+ "," + surname + "," + university + "," + vlab + "," + result)
	// Write the new account to the ledger
	err = stub.PutState(role, data)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (t *SimpleChaincode) GetAllAsset(stub shim.ChaincodeStubInterface) pb.Response {
    iterator, err := stub.GetStateByRange("", "")
    if err != nil {
        return shim.Error("Failed to get assets")
    }
    defer iterator.Close()

    var assetNames []string

    for iterator.HasNext() {
        result, err := iterator.Next()
        if err != nil {
            return shim.Error("Failed to iterate over assets")
        }

        array, err := stub.GetState(result.GetKey())
        assetNames = append(assetNames,result.Key ,string(array))
        // stringArray := strings.Split(args, ",")
        // t.query(stub, stringArray)

    }

    // Convert the asset names to JSON
    assetNamesJSON, err := json.Marshal(assetNames)
    if err != nil {
        return shim.Error("Failed to marshal asset names to JSON")
    }

    return shim.Success(assetNamesJSON)
}

func (t *SimpleChaincode) ChangeUniversity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	role := args[0]
	newUniversity := args[1]

	// Check if account exists
	accountBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountBytes == nil {
		return shim.Error("Account does not exist")
	}

	// Parse the existing data
	data := string(accountBytes)

	// Split the data into fields
	fields := strings.Split(data, ",")

	// Update the university field
	fields[2] = newUniversity

	// Join the fields back into a string
	updatedData := strings.Join(fields, ",")

	// Write the updated data to the ledger
	err = stub.PutState(role, []byte(updatedData))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}




func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
