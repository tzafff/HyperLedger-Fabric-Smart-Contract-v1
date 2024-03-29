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

type Trainer struct {
	TrainerID    	string
	FirstName		string
	LastName		string
	EmailAddress	string
	City			string
	Description 	string
	Nickname 		string
}

type Platform struct {
	PlatformID  	string
	PlatformName 	string
	EmailAddress 	string				
	Description 	string
	Trainees    	[]Trainee
	Vlabs			[]Vlab
}

type Trainee struct {
	TraineeID     		string
	FirstName    		string 
	LastName     		string 
	EmailAddress 		string 
	City  				string
	Description         string
	Nickname 			string
	ActivePlatform    	string
	Total_Exp_Points	string
	VlabPointsMap2 map[string]Vlab    `json:"Trainee_vlabs"`
}

type Vlab struct {
	VlabID 			string `json:"vlabID"`
	BoxName 		string
    Domain 			string
    SystemType 		string
    Description 	string
    ExpPoints 		string
    BoxDifficulty  	string
    TimeNeeded 		string
	Result 			string
	// Add other fields as needed
}

type Administrator struct{
	AdministratorID string
	FirstName		string
	LastName		string
	EmailAddress	string			
	City			string
	Description 	string
	Nickname 		string
}

type VlabOwner struct{
	VLabOwnerID		string
	FirstName		string
	LastName		string
    EmailAddress	string
    City			string
	Description 	string
	Nickname 		string
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
	} else if function == "TransferTrainee1" {
		return t.TransferTrainee1(stub, args)
	} else if function == "createVlab" {
		return t.createVlab(stub, args)
	} else if function == "addVlabToTrainee" {
		return t.addVlabToTrainee(stub, args)
	} else if function == "deleteTraineeFromPlatform" {
		return t.deleteTraineeFromPlatform(stub, args)
	} else if function == "createAdministrator" {
		return t.createAdministrator(stub, args)
	} else if function == "createVlabOwner" {
		return t.createVlabOwner(stub, args)
	} else if function == "addVlabToPlatform" {
		return t.addVlabToPlatform(stub, args)
	} else if function == "calculateExpPoints" {
		return t.calculateExpPoints(stub, args)
	}

	

	return shim.Error("Invalid invoke function name. Expecting \"createTrainee\" \"getTrainee\" \"updateTraineeUniversity\" \"addVLabToTrainee\"")
}

func (t *SimpleChaincode) createTrainee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	traineeID := args[0]
	firstName := args[1]
	lastName := args[2]
	emailAddress := args[3]
	city := args[4]
	description := args[5]
	nickname := args[6]

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
		TraineeID:     	traineeID,
		FirstName:      firstName,
		LastName:       lastName,
		EmailAddress:   emailAddress,
		City:    		city,
		Description:    description,
		Nickname:    	nickname,
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
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	
	platformID := args[0]
	platformName := args[1]
	emailAddress := args[2]
	description := args[3]

	// Check if platform exists
	platformBytes, err := stub.GetState(platformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if platformBytes != nil {
		return shim.Error("PlatformID already exists")
	}

	// Create a new Platform object
	platform := Platform{
		PlatformID:  platformID,
		PlatformName:   platformName,
		EmailAddress:  emailAddress,
		Description: description,
		Trainees:    []Trainee{},
		Vlabs: []Vlab{},
	}

	// Convert Platform object to JSON
	platformJSON, err := json.Marshal(platform)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save Platform JSON to the ledger
	err = stub.PutState(platformID, platformJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTraineeToPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	administratorID := args[0]
	role := args[1]
	PlatformID := args[2]

	// Check if administratorID starts with "admin"
	if !strings.HasPrefix(administratorID, "admin") {
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
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	trainerID := args[0]
	firstName := args[1]
	lastName := args[2]
	emailAddress := args[3]
	city := args[4]
	description := args[5]
	nickname := args[6]

	// Check if trainer already exists
	traineeBytes, err := stub.GetState(trainerID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if traineeBytes != nil {
		return shim.Error("TrainerID already exists")
	}

	// Create a new trainer object
	trainer := Trainer{
		TrainerID:     		trainerID,
		FirstName:    		firstName,
		LastName: 			lastName,
		EmailAddress:   	emailAddress,
		City: 				city,
		Description:    	description,
		Nickname: 			nickname,
	}

	// Convert trainee object to JSON
	trainerJSON, err := json.Marshal(trainer)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save trainee JSON to the ledger
	err = stub.PutState(trainerID, trainerJSON)
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

	// Check if trainerID starts with "trainer"
	if !strings.HasPrefix(trainerID, "Trainer") {
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


func (t *SimpleChaincode) createVlab(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check the number of arguments
	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting vlabID")
	}
	vlabOwnerId := args[0]

	// Check if administratorID starts with "admin"
	if !strings.HasPrefix(vlabOwnerId, "vlabowner") {
		return shim.Error("Not authorized for that transaction.")
	}

	// Extract the vlabID from the argument
	vlabID := args[1]
	boxName :=args[2]
	domain := args[3]
	systemType :=args[4]
	description :=args[5]
	expPoints :=args[6]
	boxDifficulty :=args[7]
	timeNeeded :=args[8]

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
		BoxName:  boxName,
		Domain: domain,
		SystemType:  systemType,
		Description:  description,
		ExpPoints:  expPoints,
		BoxDifficulty:  boxDifficulty,
		TimeNeeded:  timeNeeded,
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
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	administratorID := args[0]
	traineeID := args[1]
	platformID := args[2]

	// Check if administratorID starts with "admin"
	if !strings.HasPrefix(administratorID, "admin") {
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


func (t *SimpleChaincode) createAdministrator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	administratorID := args[0]
	firstName := args[1]
	lastName := args[2]
	emailAddress := args[3]
	city := args[4]
	description := args[5]
	nickname := args[6]
	

	// Check if Administrator already exists
	adminBytes, err := stub.GetState(administratorID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if adminBytes != nil {
		return shim.Error("AdministratorID already exists")
	}

	// Create a new Administrator object
	administrator := Administrator{
		AdministratorID:    administratorID,
		FirstName:    		firstName,
		LastName: 			lastName,
		EmailAddress:    	emailAddress,
		City: 				city,
		Description:    	description,
		Nickname: 			nickname,
	}

	// Convert Administrator object to JSON
	administratorJSON, err := json.Marshal(administrator)
	if err != nil {
		return shim.Error("Failed to marshal Administrator to JSON")
	}

	// Save Administrator JSON to the ledger
	err = stub.PutState(administratorID, administratorJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}



func (t *SimpleChaincode) createVlabOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	vlabOwnerID := args[0]
	firstName := args[1]
	lastName := args[2]
	emailAddress := args[3]
	city := args[4]
	description := args[5]
	nickname := args[6]
	

	// Check if VlabOwner already exists
	VlabOwnerBytes, err := stub.GetState(vlabOwnerID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if VlabOwnerBytes != nil {
		return shim.Error("VlabOwnerID already exists")
	}

	// Create a new VlabOwner object
	vlabOwner := VlabOwner{
		VLabOwnerID:     	vlabOwnerID,
		FirstName:    		firstName,
		LastName: 			lastName,
		EmailAddress:   	emailAddress,
		City: 				city,
		Description:    	description,
		Nickname: 			nickname,
	}

	// Convert VlabOwner object to JSON
	administratorJSON, err := json.Marshal(vlabOwner)
	if err != nil {
		return shim.Error("Failed to marshal VlabOwner to JSON")
	}

	// Save VlabOwner JSON to the ledger
	err = stub.PutState(vlabOwnerID, administratorJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (t *SimpleChaincode) addVlabToPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	administratorID := args[0]
	vlabID := args[1]
	PlatformID := args[2]

	// Check if administratorID starts with "admin"
	if !strings.HasPrefix(administratorID, "admin") {
		return shim.Error("Not authorized for that transaction.")
	}


	// Retrieve vlab from the ledger
	vlabIDBytes, err := stub.GetState(vlabID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if vlabIDBytes == nil {
		return shim.Error("vlabID does not exist")
	}
	// Unmarshal Vlab JSON
	vlab := Vlab{}
	err = json.Unmarshal(vlabIDBytes, &vlab)
	if err != nil {
		return shim.Error("Failed to unmarshal Vlab JSON")
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

	for _, val := range platform.Vlabs {
		if val.VlabID == vlabID {
			return shim.Error("vlabID is already added in the platform "+PlatformID)
		}
	}

	
	platform.Vlabs = append(platform.Vlabs, vlab)

	
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

func (t *SimpleChaincode) TransferTrainee1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	administratorID := args[0]
	traineeID := args[1]
	newPlatformID := args[2]

	// Check if administratorID starts with "admin"
	if !strings.HasPrefix(administratorID, "admin") {
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

	

	// Get the current platform of the trainee
	currentPlatformID := trainee.ActivePlatform

	
	// Retrieve the current platform from the ledger
	currentPlatformBytes, err := stub.GetState(currentPlatformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if currentPlatformBytes == nil {
		return shim.Error("Current platform does not exist")
	}

	if trainee.ActivePlatform == newPlatformID {
		return shim.Error("Trainee already signed in platform "+currentPlatformID)
	}

	// Unmarshal the current platform JSON
	currentPlatform := Platform{}
	err = json.Unmarshal(currentPlatformBytes, &currentPlatform)
	if err != nil {
		return shim.Error("Failed to unmarshal current platform JSON")
	}

	// Retrieve the new platform from the ledger
	newPlatformBytes, err := stub.GetState(newPlatformID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if newPlatformBytes == nil {
		return shim.Error("New platform does not exist")
	}

	// Unmarshal the new platform JSON
	newPlatform := Platform{}
	err = json.Unmarshal(newPlatformBytes, &newPlatform)
	if err != nil {
		return shim.Error("Failed to unmarshal new platform JSON")
	}

	// Find common VLabs between platforms
	commonVLabs := findCommonVLabs(trainee.VlabPointsMap2, newPlatform.Vlabs)

	// Remove non-common VLabs from trainee's VlabPointsMap2
	for vlabID := range trainee.VlabPointsMap2 {
		if !contains(commonVLabs, vlabID) {
			delete(trainee.VlabPointsMap2, vlabID)
		}
	}

	// Remove the trainee from the current platform
	for i, trainee := range currentPlatform.Trainees {
		if trainee.TraineeID == traineeID {
			currentPlatform.Trainees = append(currentPlatform.Trainees[:i], currentPlatform.Trainees[i+1:]...)
			break
		}
	}

	// Add the trainee to the new platform
	newPlatform.Trainees = append(newPlatform.Trainees, trainee)

	// Update the trainee's active platform
	trainee.ActivePlatform = newPlatformID

	for i, trainee := range newPlatform.Trainees {
		if trainee.TraineeID == traineeID {
			newPlatform.Trainees[i].ActivePlatform = newPlatformID
			break
		}
	}


	
	// Convert trainee object to JSON
	traineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal trainee to JSON")
	}

	// Save updated trainee JSON to the ledger
	err = stub.PutState(traineeID, traineeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert current platform object to JSON
	currentPlatformJSON, err := json.Marshal(currentPlatform)
	if err != nil {
		return shim.Error("Failed to marshal current platform to JSON")
	}

	// Save updated current platform JSON to the ledger
	err = stub.PutState(currentPlatformID, currentPlatformJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Convert new platform object to JSON
	newPlatformJSON, err := json.Marshal(newPlatform)
	if err != nil {
		return shim.Error("Failed to marshal new platform to JSON")
	}

	// Save updated new platform JSON to the ledger
	err = stub.PutState(newPlatformID, newPlatformJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Helper function to find common VLabs between trainee and platform
func findCommonVLabs(traineeVLabs map[string]Vlab, platformVLabs []Vlab) []string {
	commonVLabs := []string{}
	traineeVlabIDs := make(map[string]bool)

	// Get the keys (VLab IDs) from the trainee's VlabPointsMap2
	for vlabID := range traineeVLabs {
		traineeVlabIDs[vlabID] = true
	}

	// Find the common VLab IDs between the trainee and platform
	for _, vlab := range platformVLabs {
		if traineeVlabIDs[vlab.VlabID] {
			commonVLabs = append(commonVLabs, vlab.VlabID)
		}
	}

	return commonVLabs
}

// Helper function to check if a string slice contains a given string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}



func (t *SimpleChaincode) calculateExpPoints(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	traineeID := args[0]

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

	// Calculate the total experience points
	expPoints := 0
	for _, vlab := range trainee.VlabPointsMap2 {
		if vlab.Result != "" {
			// Convert the result to an integer and add it to expPoints
			result, err := strconv.Atoi(vlab.Result)
			if err != nil {
				return shim.Error("Failed to convert vlab result to integer")
			}
			expPoints += result
		}
	}

	// Update the trainee's experience points
	trainee.Total_Exp_Points = strconv.Itoa(expPoints)

	// Convert trainee object to JSON
	updatedTraineeJSON, err := json.Marshal(trainee)
	if err != nil {
		return shim.Error("Failed to marshal updated trainee to JSON")
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

	for i, trainee := range platform.Trainees {
		if trainee.TraineeID == traineeID {
			platform.Trainees[i].Total_Exp_Points = strconv.Itoa(expPoints)
			break
		}
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

	// Save updated trainee JSON to the ledger
	err = stub.PutState(traineeID, updatedTraineeJSON)
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