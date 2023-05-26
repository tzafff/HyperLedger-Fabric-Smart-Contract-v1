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

type Platform struct {
	PlatformID  string
	Boxname     string
	Domain      string
	SystemType  string
	Description string
	ExpPoints   string
	Trainees    []Trainee
}

type Trainee struct {
	TraineeID     		string
	Name          		string `json:"name"`
	Surname       		string `json:"surname"`
	University    		string `json:"university"`
	ActivePlatform    	string
	//VlabPointsMap 		map[string]string `json:"vlabPoints"`
	VlabPointsMap2 map[string]Vlab    `json:"Trainee_vlabs"`
}

type Vlab struct {
	VlabID string `json:"vlabID"`
	Ects string
	Result string
	// Add other fields as needed
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
	} else if function == "addTraineeToPlatform" {
		return t.addTraineeToPlatform(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "ScoreTheVlab" {
		return t.ScoreTheVlab(stub, args)
	} else if function == "createTrainer" {
		return t.createTrainer(stub, args)
	} else if function == "createPlatform" {
		return t.createPlatform(stub, args)
	} else if function == "getAllAsset" {
		return t.getAllAsset(stub, args)
	} else if function == "TransferTrainee" {
		return t.TransferTrainee(stub, args)
	} else if function == "createVlab" {
		return t.createVlab(stub, args)
	} else if function == "addVlabToTrainee" {
		return t.addVlabToTrainee(stub, args)
	} else if function == "deleteTraineeFromPlatform" {
		return t.deleteTraineeFromPlatform(stub, args)
	}

	

	return shim.Error("Invalid invoke function name. Expecting \"createTrainee\" \"getTrainee\" \"updateTraineeUniversity\" \"addVLabToTrainee\"")
}

func (t *SimpleChaincode) createTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	traineeID := args[0]
	name := args[1]
	surname := args[2]
	university := args[3]

	// Check if trainee already exists
	traineeBytes, err := stub.GetState(traineeID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes != nil {
		return shim.Error("TraineeID already exists")
	}

	// Create a new trainee object
	trainee := Trainee{
		TraineeID:     traineeID,
		Name:          name,
		Surname:       surname,
		University:    university,
		//VlabPointsMap: make(map[string]string),
		VlabPointsMap2: make(map[string]Vlab),
	}

	// Convert trainee object to JSON
	traineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save trainee JSON to the ledger
	err = stub.PutState(traineeID, traineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) createPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	
	PlatformID := args[0]
	SystemType := args[1]
	ExpPoints := args[2]
	Domain := args[3]
	Description := args[4]
	Boxname := args[5]

	// Check if platform exists
	platformBytes, err := stub.GetState(PlatformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if platformBytes != nil {
		return shim.Error("PlatformID already exists")
	}

	// Create a new Platform object
	platform := Platform{
		PlatformID:  PlatformID,
		ExpPoints:   ExpPoints,
		SystemType:  SystemType,
		Domain:      Domain,
		Description: Description,
		Boxname:     Boxname,
		Trainees:    []Trainee{},
	}

	// Convert Platform object to JSON
	platformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save Platform JSON to the ledger
	err = stub.PutState(PlatformID, platformJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTraineeToPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	role := args[0]
	PlatformID := args[1]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(role)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("TraineeID does not exist")
	}
	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	// Check if Platform exists
	platformBytes, err := stub.GetState(PlatformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if platformBytes == nil {
		return shim.Error("PlatformID does not exists")
	}
	platform := Platform{}
	err = json.Unmarshal(platformBytes, &platform)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	for _, val := range platform.Trainees {
		if val.TraineeID == role {
			return shim.Error("TraineeID is already registered in the platform "+PlatformID)
		}
	}

	if trainee.ActivePlatform == "" {
		trainee.ActivePlatform = PlatformID
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
		platform.Trainees = append(platform.Trainees, trainee)
	} else {
		return shim.Error("TraineeID is already registered in ActivePlatform, you need to transfer first")

	}

	// Convert Platform object to JSON
	updatedPlatformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error("Failed to marshal updated Platform to JSON")
	}

	// Save updated Platform JSON to the ledger
	err = stub.PutState(PlatformID, updatedPlatformJSON)
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
		return shim.Error("TrainerID already exists")
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
		return shim.Error("TraineeID does not exist")
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
	vlabResult := args[3]

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

	// Check if the vlabID exists in VlabPointsMap2
	if _, exists := trainee.VlabPointsMap2[vlabID]; !exists {
		return shim.Error("Trainee does not have that vlabID")
	}


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

	// Update trainee's vlab points
	vlab.Result = vlabResult
	trainee.VlabPointsMap2[vlabID] = vlab

	// Get the trainee's platform from the ledger
	platformBytes, err := stub.GetState(trainee.ActivePlatform)
	if err != nil {
		return shim.Error(err.Error())
	}


	// Unmarshal the platform JSON
	var platform Platform
	err = json.Unmarshal(platformBytes, &platform)
	if err != nil {
		return shim.Error(err.Error())
	}

	for i, trainee := range platform.Trainees {
		if trainee.TraineeID == traineeID {
			platform.Trainees[i].VlabPointsMap2[vlabID] = vlab
			break
		}
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


	// Convert the updated platform to JSON
	updatedPlatformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Store the updated trainee in the ledger
	err = stub.PutState(trainee.ActivePlatform, updatedPlatformJSON)
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
	from_Platform := args[2]
	to_Platform := args[3]

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
		return shim.Error("TraineeID does not exist")
	}

	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	if trainee.ActivePlatform == "" {
		return shim.Error("Trainee does not belong in any Platform")
	}

	// Retrieve the from Platform from the ledger
	from_platformBytes, err := stub.GetState(from_Platform)
	if err != nil {
		return shim.Error(err.Error())
	}
	if from_platformBytes == nil {
		return shim.Error("PlatformID does not exist")
	}

	// Unmarshal the from Platform JSON
	platform1 := Platform{}
	err = json.Unmarshal(from_platformBytes, &platform1)
	if err != nil {
		return shim.Error("Failed to unmarshal from Platform JSON")
	}

	// Retrieve the to Platform from the ledger
	to_platformBytes, err := stub.GetState(to_Platform)
	if err != nil {
		return shim.Error(err.Error())
	}
	if to_platformBytes == nil {
		return shim.Error("PlatformID does not exist")
	}

	// Unmarshal the to Platform JSON
	platform2 := Platform{}
	err = json.Unmarshal(to_platformBytes, &platform2)
	if err != nil {
		return shim.Error("Failed to unmarshal to Vlab JSON")
	}

	// Find the trainee within the Platform
	var found bool
	var indexToRemove int
	for i, trainee := range platform1.Trainees {
		if trainee.TraineeID == traineeID {
			indexToRemove = i
			found = true
			break
		}
	}

	platform1.Trainees = append(platform1.Trainees[:indexToRemove], platform1.Trainees[indexToRemove+1:]...)

	platform2.Trainees = append(platform2.Trainees, trainee)
	if !found {
		return shim.Error("Trainee not found in Vlab")
	}

	// Convert the updated from Platform object to JSON
	updatedPlatform1JSON, err := json.Marshal(platform1)
	if err != nil {
		return shim.Error("Failed to marshal updated Vlab to JSON")
	}

	for i, trainee := range platform2.Trainees {
		if trainee.TraineeID == traineeID {
			platform2.Trainees[i].ActivePlatform = to_Platform
			break
		}
	}
	trainee.ActivePlatform = to_Platform
	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}
	// Save updated trainee JSON to the ledger
	err = stub.PutState(traineeID, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save the updated from Vlab JSON to the ledger
	err = stub.PutState(from_Platform, updatedPlatform1JSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert the updated to Vlab object to JSON
	updatedPlatform2JSON, err := json.Marshal(platform2)
	if err != nil {
		return shim.Error("Failed to marshal updated Vlab to JSON")
	}

	// Save the updated to Vlab JSON to the ledger
	err = stub.PutState(to_Platform, updatedPlatform2JSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) createVlab(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check the number of arguments
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting vlabID")
	}

	// Extract the vlabID from the argument
	vlabID := args[0]
	ects :=args[1]

	// Check if the Vlab already exists in the ledger
	vlabBytes, err := stub.GetState(vlabID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if vlabBytes != nil {
		return shim.Error("Vlab with ID "+vlabID+" already exists")
	}

	// Create a new Vlab instance
	vlab := Vlab{
		VlabID: vlabID,
		Ects:  ects,
		// Set other fields as needed
	}

	// Convert the Vlab instance to JSON
	vlabJSON, err := json.Marshal(vlab)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Store the Vlab JSON in the ledger
	err = stub.PutState(vlabID, vlabJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addVlabToTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check the number of arguments
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting traineeID and vlabID")
	}

	// Extract the traineeID and vlabID from arguments
	traineeID := args[0]
	vlabID := args[1]

	// Get the existing trainee from the ledger
	traineeBytes, err := stub.GetState(traineeID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Check if the trainee exists
	if traineeBytes == nil {
		return shim.Error("Trainee with ID "+traineeID+" does not exist")
	}

	// Unmarshal the trainee JSON
	var trainee Trainee
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error(err.Error())
	}

	if trainee.ActivePlatform == "" {
		return shim.Error("Trainee need to be registered in platform in order to add a vlab")
	}
	

	// Get the existing vlab from the ledger
	vlabBytes, err := stub.GetState(vlabID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Check if the vlab exists
	if vlabBytes == nil {
		return shim.Error("Vlab with ID "+vlabID+" does not exist")
	}

	// Unmarshal the vlab JSON
	var vlab Vlab
	err = json.Unmarshal(vlabBytes, &vlab)
	if err != nil {
		return shim.Error(err.Error())
	}



	// Get the trainee's platform from the ledger
	platformBytes, err := stub.GetState(trainee.ActivePlatform)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	// Unmarshal the platform JSON
	var platform Platform
	err = json.Unmarshal(platformBytes, &platform)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// Check if the Vlab already exists in the Trainee's VlabPointsMap
	if _, exists := trainee.VlabPointsMap2[vlabID]; exists {
		return shim.Error("Vlab with ID "+vlabID+" already exists for trainee with ID "+ vlabID)
	}

	// Add the Vlab to the Trainee's VlabPointsMap
	trainee.VlabPointsMap2[vlabID] = vlab


	for i, trainee := range platform.Trainees {
		if trainee.TraineeID == traineeID {
			platform.Trainees[i].VlabPointsMap2[vlabID] = vlab
			break
		}
	}


	// Convert the updated trainee to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Store the updated trainee in the ledger
	err = stub.PutState(traineeID, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}


	// Convert the updated platform to JSON
	updatedPlatformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Store the updated trainee in the ledger
	err = stub.PutState(trainee.ActivePlatform, updatedPlatformJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) deleteTraineeFromPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	traineeID := args[0]
	platformID := args[1]

	// Retrieve trainee from the ledger
	traineeBytes, err := stub.GetState(traineeID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes == nil {
		return shim.Error("TraineeID does not exist")
	}

	// Unmarshal trainee JSON
	trainee := Trainee{}
	err = json.Unmarshal(traineeBytes, &trainee)
	if err != nil {
		return shim.Error("Failed to unmarshal trainee JSON")
	}

	// Check if the trainee is registered in the specified platform
	if trainee.ActivePlatform != platformID {
		return shim.Error("Trainee is not registered in the specified platform")
	}

	// Retrieve the platform from the ledger
	platformBytes, err := stub.GetState(platformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if platformBytes == nil {
		return shim.Error("PlatformID does not exist")
	}

	// Unmarshal platform JSON
	platform := Platform{}
	err = json.Unmarshal(platformBytes, &platform)
	if err != nil {
		return shim.Error("Failed to unmarshal platform JSON")
	}

	// Remove the trainee from the platform's Trainees list
	updatedTrainees := make([]Trainee, 0)
	for _, trainee := range platform.Trainees {
		if trainee.TraineeID != traineeID {
			updatedTrainees = append(updatedTrainees, trainee)
		}
	}
	platform.Trainees = updatedTrainees

	// Update the trainee's active platform to empty
	trainee.ActivePlatform = ""

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(traineeID, updatedTraineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert platform object to JSON
	updatedPlatformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error("Failed to marshal updated platform to JSON")
	}

	// Save updated platform JSON to the ledger
	err = stub.PutState(platformID, updatedPlatformJSON)
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