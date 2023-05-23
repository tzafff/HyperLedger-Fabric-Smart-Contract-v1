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

type Trainer struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	University string `json:"university"`
}

type Vlab struct {
	Boxname     string
	Domain      string
	SystemType  string
	Description string
	ExpPoints   string
	Trainees    []Trainee
}

type Trainee struct {
	TraineeID  string
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	University string `json:"university"`
	VlabPoints string
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
	} else if function == "getIdentity" {
		return t.getIdentity(stub, args)
	} else if function == "updateTraineeUniversity" {
		return t.updateTraineeUniversity(stub, args)
	} else if function == "addTraineeToVLab" {
		return t.addTraineeToVLab(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "ScoreTheVlab" {
		return t.ScoreTheVlab(stub, args)
	} else if function == "createTrainer" {
		return t.createTrainer(stub, args)
	} else if function == "createVlab" {
		return t.createVlab(stub, args)
	} else if function == "getAllAsset" {
		return t.getAllAsset(stub, args)
	} else if function == "TransferTrainee" {
		return t.TransferTrainee(stub, args)
	}
	

	return shim.Error("Invalid invoke function name. Expecting \"createTrainee\" \"getTrainee\" \"updateTraineeUniversity\" \"addVLabToTrainee\"")
}

func (t *SimpleChaincode) createTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
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
		TraineeID:  role,
		Name:       name,
		Surname:    surname,
		University: university,
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

func (t *SimpleChaincode) createVlab(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// randomInt := rand.Intn(100)
	//  + strconv.Itoa(randomInt)
	VlabID := args[0]
	SystemType := args[1]
	ExpPoints := args[2]
	Domain := args[3]
	Description := args[4]
	Boxname := args[5]

	// Check if vlab already exists
	vlabsBytes, err := stub.GetState(VlabID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if vlabsBytes != nil {
		return shim.Error("Vlab already exists")
	}

	// Create a new Vlab object
	vlab := Vlab{
		ExpPoints:   ExpPoints,
		SystemType:  SystemType,
		Domain:      Domain,
		Description: Description,
		Boxname:     Boxname,
		Trainees:    []Trainee{},
	}

	// Convert trainee object to JSON
	vlabJSON, err := json.Marshal(vlab)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save trainee JSON to the ledger
	err = stub.PutState(VlabID, vlabJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTraineeToVLab(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	role := args[0]
	VlabID := args[1]

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

	// Check if vlab already exists
	vlabsBytes, err := stub.GetState(VlabID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if vlabsBytes == nil {
		return shim.Error("Vlab already exists")
	}
	vlab := Vlab{}
	err = json.Unmarshal(vlabsBytes, &vlab)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	for _, val := range vlab.Trainees {
		if val.TraineeID == role {
			return shim.Error("Traine is already exists in Vlab with Id %d")
		}
	}
	vlab.Trainees = append(vlab.Trainees, trainee)

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(vlab)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(VlabID, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) createTrainer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	role := args[0]
	name := args[1]
	surname := args[2]
	university := args[3]

	// Check if trainer already exists
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes != nil {
		return shim.Error("Trainer already exists")
	}

	// Create a new trainer object
	trainer := Trainer{
		Name:       name,
		Surname:    surname,
		University: university,
	}

	// Convert trainee object to JSON
	trainerJSON, err := json.Marshal(trainer)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save trainee JSON to the ledger
	err = stub.PutState(role, trainerJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) getIdentity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	role := args[0]

	// Retrieve trainee from the ledger
	identityBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if identityBytes == nil {
		return shim.Error("Identity does not exist")
	}

	return shim.Success(identityBytes)
}

func (t *SimpleChaincode) updateTraineeUniversity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	trainerID := args[0]
	role := args[1]
	newUniversity := args[2]

	// Add authorized validation
	if trainerID != "Trainer" {
		return shim.Error("Not authorized for that transaction.")
	}

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
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	trainerID := args[0]
	traineeID := args[1]
	vlabID := args[2]
	vlabPoints := args[3]

	// Add authorized validation
	if trainerID != "Trainer" {
		return shim.Error("Not authorized for that transaction.")
	}

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(traineeID)
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
	trainee.VlabPoints = vlabPoints



	// Retrieve the Vlab from the ledger
	vlabBytes, err := stub.GetState(vlabID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if vlabBytes == nil {
		return shim.Error("Vlab does not exist")
	}

	// Unmarshal the Vlab JSON
	vlab := Vlab{}
	err = json.Unmarshal(vlabBytes, &vlab)
	if err != nil {
		return shim.Error("Failed to unmarshal Vlab JSON")
	}

	// Find the trainee within the Vlab
	var found bool
	for i, trainee := range vlab.Trainees {
		if trainee.TraineeID == traineeID {
			vlab.Trainees[i].VlabPoints = vlabPoints
			found = true
			break
		}
	}

	if !found {
		return shim.Error("Trainee not found in Vlab")
	}

	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(traineeID, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert the updated Vlab object to JSON
	updatedVlabJSON, err := json.Marshal(vlab)
	if err != nil {
		return shim.Error("Failed to marshal updated Vlab to JSON")
	}

	// Save the updated Vlab JSON to the ledger
	err = stub.PutState(vlabID, updatedVlabJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) getAllAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Retrieve all assets from the ledger
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error("Failed to get assets")
	}
	defer resultsIterator.Close()

	// Iterate over the result set and collect the assets
	var assets []string
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Failed to iterate over assets")
		}
		assets = append(assets, string(queryResult.Value))
	}

	// Convert assets to JSON
	assetsJSON, err := json.Marshal(assets)
	if err != nil {
		return shim.Error("Failed to marshal assets to JSON")
	}

	return shim.Success(assetsJSON)
}


func (t *SimpleChaincode) TransferTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	trainerID := args[0]
	traineeID := args[1]
	from_Vlab := args[2]
	to_Vlab := args[3]

	// Add authorized validation
	if trainerID != "Trainer" {
		return shim.Error("Not authorized for that transaction.")
	}

	
	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(traineeID)
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


	// Retrieve the from Vlab from the ledger
	from_vlabBytes, err := stub.GetState(from_Vlab)
	if err != nil {
		return shim.Error(err.Error())
	}
	if from_vlabBytes == nil {
		return shim.Error("Vlab does not exist")
	}

	// Unmarshal the from Vlab JSON
	vlab1 := Vlab{}
	err = json.Unmarshal(from_vlabBytes, &vlab1)
	if err != nil {
		return shim.Error("Failed to unmarshal from Vlab JSON")
	}


	// Retrieve the to Vlab from the ledger
	to_vlabBytes, err := stub.GetState(to_Vlab)
	if err != nil {
		return shim.Error(err.Error())
	}
	if to_vlabBytes == nil {
		return shim.Error("Vlab does not exist")
	}

	// Unmarshal the to Vlab JSON
	vlab2 := Vlab{}
	err = json.Unmarshal(to_vlabBytes, &vlab2)
	if err != nil {
		return shim.Error("Failed to unmarshal to Vlab JSON")
	}

	// Find the trainee within the Vlab
	var found bool
	var indexToRemove int 
	for i, trainee := range vlab1.Trainees {
		if trainee.TraineeID == traineeID {
			indexToRemove = i
			found = true
			break
		}
	}

	vlab1.Trainees = append(vlab1.Trainees[:indexToRemove],vlab1.Trainees[indexToRemove+1:]... )


	vlab2.Trainees = append(vlab2.Trainees, trainee)
	if !found {
		return shim.Error("Trainee not found in Vlab")
	}


	// Convert the updated from Vlab object to JSON
	updatedVlab1JSON, err := json.Marshal(vlab1)
	if err != nil {
		return shim.Error("Failed to marshal updated Vlab to JSON")
	}

	// Save the updated from Vlab JSON to the ledger
	err = stub.PutState(from_Vlab, updatedVlab1JSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert the updated to Vlab object to JSON
	updatedVlab2JSON, err := json.Marshal(vlab2)
	if err != nil {
		return shim.Error("Failed to marshal updated Vlab to JSON")
	}

	// Save the updated to Vlab JSON to the ledger
	err = stub.PutState(to_Vlab, updatedVlab2JSON)
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
