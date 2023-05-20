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
	
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric-chaincode-go/shim"

	// pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Trainee struct {
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	University string   `json:"university"`
	VLab       []string `json:"vlab"`
}

/*
...
*/

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createTrainee" {
		return t.createTrainee(stub, args)
	} else if function == "getTrainee" {
		return t.getTrainee(stub, args)
	} else if function == "updateTraineeUniversity" {
		return t.updateTraineeUniversity(stub, args)
	} else if function == "addVLabToTrainee" {
		return t.addVLabToTrainee(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "ScoreTheVlab" {
		return t.ScoreTheVlab(stub, args)
	}
	

	return shim.Error("Invalid invoke function name. Expecting \"createTrainee\" \"getTrainee\" \"updateTraineeUniversity\" \"addVLabToTrainee\"")
}

func (t *SimpleChaincode) createTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	role := args[0]
	name := args[1]
	surname := args[2]
	university := args[3]


	// Check if trainee already exists
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes != nil {
		return shim.Error("Trainee already exists")
	}

	// Create a new trainee object
	trainee := Trainee{
		Name:       name,
		Surname:    surname,
		University: university,
		VLab:       []string{},
	
	}

	// Convert trainee object to JSON
	traineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save trainee JSON to the ledger
	err = stub.PutState(role, traineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) getTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	role := args[0]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("Trainee does not exist")
	}

	return shim.Success(traineeBytes)
}

func (t *SimpleChaincode) updateTraineeUniversity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	role := args[0]
	newUniversity := args[1]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("Trainee does not exist")
	}

	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	// Update trainee's university
	trainee.University = newUniversity

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(role, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (t *SimpleChaincode) addVLabToTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	role := args[0]
	vLabName := args[1]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("Trainee does not exist")
	}

	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	// Add virtual lab name to trainee's list of virtual labs
	trainee.VLab = append(trainee.VLab, vLabName)

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(role, updatedTraineeJSON)
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


func (t *SimpleChaincode) ScoreTheVlab(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	role := args[0]
	vLabName := args[1]
	score := args[2]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("Trainee does not exist")
	}

	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	// Update the score of the specified virtual lab
	for i, lab := range trainee.VLab {
		if lab == vLabName {
			trainee.VLab[i] = vLabName + " = " + score
			break
		}
	}

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(role, updatedTraineeJSON)
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