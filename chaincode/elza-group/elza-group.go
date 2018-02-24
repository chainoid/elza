// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports
* 5 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Group structure, with several properties.
Structure tags are used by encoding/json library
*/
type Group struct {
	GroupId string `json:"groupId"`
	GroupName string `json:"groupName"`
	GroupDesc string `json:"groupDesc"`
}


/*
 * The Init method *
 called when the Smart Contract "elza-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "posta-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryGroupById" {
		return s.queryGroupById(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createGroup" {
		return s.createGroup(APIstub, args)
	} else if function == "queryAllGroups" {
		return s.queryAllGroups(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryParsel method *
 Used to view the records of one particular parsel
 It takes one argument -- the key for the parsel in question
 */
func (s *SmartContract) queryGroupById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	groupAsBytes, _ := APIstub.GetState(args[0])
	if groupAsBytes == nil {
		return shim.Error("Could not locate parsel")
	}
	return shim.Success(groupAsBytes)
}

/*
 * The initLedger method *
Will add test data to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	groups := []Group{
		Group{GroupId: "001",  GroupName: "AC17",  GroupDesc: "Description for AB17"},
		Group{GroupId: "002",  GroupName: "AC18",  GroupDesc: "Description for AB18"},
	}

	i := 0
	for i < len(groups) {
		fmt.Println("i is ", i)
		groupAsBytes, _ := json.Marshal(groups[i])
		APIstub.PutState(strconv.Itoa(i+1), groupAsBytes)
		fmt.Println("Added", groups[i])
		i = i + 1
	}

	return shim.Success(nil)
}


/*
 * The createGroup method *
   This method takes in four arguments (attributes to be saved in the ledger).
 */

  func (s *SmartContract) createGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var group = Group{ GroupId: args[1], GroupName: args[2], GroupDesc: args[3] }

	groupAsBytes, _ := json.Marshal(group)
	err := APIstub.PutState(args[0], groupAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record new group: %s", args[0]))
	}

	return shim.Success(nil)
}


/*
 * The queryAllGroups method *
allows for assessing all the records added to the ledger(all groups in the delivery system)
This method does not take any arguments. Returns JSON string containing results.
 */
func (s *SmartContract) queryAllGroups(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllGroups:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


/*
 * The querySender method *
allows for assessing all the records from selected sender

Returns JSON string containing results.
 */
/*
 * func (s *SmartContract) querySender(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		parsel := Parsel{}
		// Unmarshal record to parsel
		json.Unmarshal(queryResponse.Value, &parsel)

		// Add only filtered ny sender records
		if parsel.Sender == args[0] {

		// Add comma before array members,suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}

			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
	    }
	}
	buffer.WriteString("]")

	if bArrayMemberAlreadyWritten == false {
		return shim.Error("No parsels for sender")
	}

	fmt.Printf("- querySender:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
*/

/*
 * The deliveryParsel method *
The data in the world state can be updated with who has possession.
This function takes in 2 arguments, parsel id and timestamp of delivery.
 */
/*func (s *SmartContract) deliveryParsel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	parselAsBytes, _ := APIstub.GetState(args[0])
	if parselAsBytes == nil {
		return shim.Error("Could not locate parsel")
	}
	parsel := Parsel{}

	json.Unmarshal(parselAsBytes, &parsel)
	// Normally check that the specified argument is a valid holder of parsel
	// we are skipping this check for this example
	parsel.ReceiverTS = time.Now().Format(time.RFC1123Z)

	parselAsBytes, _ = json.Marshal(parsel)
	err := APIstub.PutState(args[0], parselAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change status of parsel: %s", args[0]))
	}

	return shim.Success(nil)
}
* /

/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
