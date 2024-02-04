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
	"encoding/json"
	"fmt"
	"strconv"

	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("example_cc0")

type tutor struct {
	Id      string
	Name    string
	Surname string
}

type tutorial struct {
	Id     string
	Name   string
	Tutors []string
}

// Global variables for ID
var tutorId int
var tutorialId int

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Init ###########")

	var tutor1 = tutor{"tu1", "John", "Doe"}
	var tutor2 = tutor{"tu2", "Michel", "Green"}
	var tutor3 = tutor{"tu3", "Jova", "Jovanovic"}
	tutorId = 4

	var tutorsFor1 = make([]string, 0, 20)
	tutorsFor1 = append(tutorsFor1, tutor1.Id)
	tutorsFor1 = append(tutorsFor1, tutor2.Id)
	// var tutorial1 = tutorial{"t1","Blockcahin tutorial",tutorsFor1}
	var tutorial1 = tutorial{"t1", "Blockcahin tutorial", tutorsFor1}

	var tutorsFor2 = make([]string, 0, 20)
	tutorsFor2 = append(tutorsFor2, tutor3.Id)
	var tutorial2 = tutorial{"t2", "Spark tutorial", tutorsFor2}
	tutorialId = 3

	// Write the state to the ledger
	ajson, _ := json.Marshal(tutor1)
	err := stub.PutState("tu1", ajson)
	if err != nil {
		return shim.Error(err.Error())
	}
	ajson, _ = json.Marshal(tutor2)
	err = stub.PutState("tu2", ajson)
	if err != nil {
		return shim.Error(err.Error())
	}
	ajson, _ = json.Marshal(tutor3)
	err = stub.PutState("tu3", ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	ajson, _ = json.Marshal(tutorial1)
	err = stub.PutState("t1", ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	ajson, _ = json.Marshal(tutorial2)
	err = stub.PutState("t2", ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "delete" {
		return t.delete(stub, args)
	}
	if function == "query" {
		return t.query(stub, args)
	}
	if function == "addTutorial" {
		return t.addTutorial(stub, args)
	}
	if function == "addTutor" {
		return t.addTutor(stub, args)
	}
	if function == "addTutorToTutorial" {
		return t.addTutorToTutorial(stub, args)
	}
	if function == "removeTutorFromTutorial" {
		return t.removeTutorFromTutorial(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (t *SimpleChaincode) addTutor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, surname string // Entities

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	name = args[0]
	surname = args[1]

	tutorKey := "tu" + strconv.Itoa(tutorId)
	tutorId = tutorId + 1
	var newTutor = tutor{tutorKey, name, surname}

	ajson, _ := json.Marshal(newTutor)
	err := stub.PutState(newTutor.Id, ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTutorial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// TODO implement function
	// arg 0 - name, arg1,arg2,arg3,arg4... - tutorID (which is the same as tutorKey)
	// Check number of arguments
	// Check if tutors exist in ledger before adding them to tutorial
	var name, tutorKey string // Entities

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	name = args[0]
	var tutors = make([]string, 0, 20)
	for i := 1; i < len(args); i++ {
		tutorKey = args[i]
		tutorI, err := stub.GetState(tutorKey)
		logger.Info("Tutor " + tutorKey + "postoji")
		logger.Info(tutorI)

		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + tutorKey + "\"}"
			return shim.Error(jsonResp)
		}
		if tutorI == nil || len(tutorI) == 0 {
			jsonResp := "{\"Error\":\" " + tutorKey + " does not exit " + "\"}"
			return shim.Error(jsonResp)
		}

		tutors = append(tutors, tutorKey)
	}

	tutorialKey := "t" + strconv.Itoa(tutorialId)
	tutorialId = tutorialId + 1
	var newTutorial = tutorial{tutorialKey, name, tutors}

	ajson, _ := json.Marshal(newTutorial)
	err := stub.PutState(newTutorial.Id, ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTutorToTutorial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// TODO implement function
	// arg0 - tutrorialId (which is the same as tutorialKey), arg1 - tutorId
	// Check number of arguments
	// Check if tutor and tutorial exist in ledger
	// Make sure that tutor is not already listed in tutorial. If that is the case, return error
	var tutorialKey, tutorKey string // Entities

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	tutorialKey = args[0]
	tutorKey = args[1]

	// load tutorial
	tutorialB, err := stub.GetState(tutorialKey)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + tutorialKey + "\"}"
		return shim.Error(jsonResp)
	}
	if tutorialB == nil || len(tutorialB) == 0 {
		jsonResp := "{\"Error\":\" " + tutorialKey + " does not exit " + "\"}"
		return shim.Error(jsonResp)
	}
	tutorial := tutorial{}
	err = json.Unmarshal(tutorialB, &tutorial)
	if err != nil {
		return shim.Error("Failed to get state")
	}

	// load tutor which will be added to tutorial
	tutor, err := stub.GetState(tutorKey)
	logger.Info("Tutor " + tutorKey + "postoji")
	logger.Info(tutor)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + tutorKey + "\"}"
		return shim.Error(jsonResp)
	}
	if tutor == nil || len(tutor) == 0 {
		jsonResp := "{\"Error\":\" " + tutorKey + " does not exit " + "\"}"
		return shim.Error(jsonResp)
	}

	for i := 0; i < len(tutorial.Tutors); i++ {
		if tutorial.Tutors[i] == tutorKey {
			jsonResp := "{\"Error\":\" Tutor with id: " + tutorKey + " is already on list to tutors" + "\"}"
			return shim.Error(jsonResp)
		}
	}

	tutorial.Tutors = append(tutorial.Tutors, tutorKey)

	ajson, _ := json.Marshal(tutorial)
	err = stub.PutState(tutorial.Id, ajson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) removeTutorFromTutorial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// TODO implement function
	// arg0 - tutrorialId, arg1 - tutorId
	// Check number of arguments
	// Check if tutor and tutorial exist in ledger
	// If tutor (which we want to remove) is not listed in tutorial, return error
	var tutorialKey, tutorKey string // Entities

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	tutorialKey = args[0]
	tutorKey = args[1]

	// load tutorial
	tutorialB, err := stub.GetState(tutorialKey)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + tutorialKey + "\"}"
		return shim.Error(jsonResp)
	}
	if tutorialB == nil || len(tutorialB) == 0 {
		jsonResp := "{\"Error\":\" " + tutorialKey + " does not exit " + "\"}"
		return shim.Error(jsonResp)
	}
	tutorial := tutorial{}
	err = json.Unmarshal(tutorialB, &tutorial)
	if err != nil {
		return shim.Error("Failed to get state")
	}

	// load tutor which will be removed from tutorial
	tutor, err := stub.GetState(tutorKey)
	logger.Info("Tutor " + tutorKey + "postoji")
	logger.Info(tutor)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + tutorKey + "\"}"
		return shim.Error(jsonResp)
	}
	if tutor == nil || len(tutor) == 0 {
		jsonResp := "{\"Error\":\" " + tutorKey + " does not exit " + "\"}"
		return shim.Error(jsonResp)
	}

	for i := 0; i < len(tutorial.Tutors); i++ {
		if tutorial.Tutors[i] == tutorKey {

			tutorial.Tutors = append(tutorial.Tutors[:i], tutorial.Tutors[i+1:]...)
			ajson, _ := json.Marshal(tutorial)
			err = stub.PutState(tutorial.Id, ajson)
			if err != nil {
				return shim.Error(err.Error())
			}

			return shim.Success(nil)
		}
	}

	// If tutor is not removed then it does not exits in list of tutors for given tutorial
	jsonResp := "{\"Error\":\" Tutor with id: " + tutorKey + " is not on list of tutors" + "\"}"
	return shim.Error(jsonResp)
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

// Query callback representing the query of a chaincode
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

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
