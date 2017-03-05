package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

//AssetManagementChaincode APIs exposed to chaincode callers
type AssetManagementChaincode struct {
}

var myLogger = logging.MustGetLogger("asset_mgm")

// Init initialization, this method will create asset despository in the chaincode state
func (t *AssetManagementChaincode) Init(stub shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var columnDefsTableTwo []*shim.ColumnDefinition
	columnOneTableTwoDef := shim.ColumnDefinition{Name: nameKey,
		Type: shim.ColumnDefinition_STRING, Key: true}
	columnTwoTableTwoDef := shim.ColumnDefinition{Name: nameValue,
		Type: shim.ColumnDefinition_STRING, Key: false}

	columnDefsTableTwo = append(columnDefsTableTwo, &columnOneTableTwoDef)
	columnDefsTableTwo = append(columnDefsTableTwo, &columnTwoTableTwoDef)

	// Create asset depository table
	return nil, stub.CreateTable("TableString", columnDefsTableTwo)
}

// Invoke  method is the interceptor of all invocation transactions, its job is to direct
// invocation transactions to intended APIs
func (t *AssetManagementChaincode) Invoke(stub shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	fmt.Println(len(args))

	ok, err := stub.InsertRow("TableString", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}}},
	})

	// you can only assign balances to new account IDs
	if !ok && err == nil {
		myLogger.Errorf("system error %v", err)
		return nil, errors.New("Asset was already assigned." + strconv.Itoa(len(args)))
	}

	return nil, nil
}

// Query method is the interceptor of all invocation transactions, its job is to direct
// query transactions to intended APIs, and return the result back to callers
func (t *AssetManagementChaincode) Query(stub shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("TableString", columns)
	if err != nil {
		return nil, fmt.Errorf("getRowTableTwo operation failed. %s", err)
	}

	var strResult, strRow string
	strRow = "[" + row.Columns[0].GetString_() + ":"
	strRow += row.Columns[1].GetString_() + "]"
	strResult += strRow

	//	rowChannel, err := stub.GetAllRows("TableAsset")
	//	if err != nil {
	//		return nil, errors.New("Query operation fail")
	//	}

	//	var rows []shim.Row
	//	var strResult, strRow string
	//	for {
	//		select {
	//		case row, ok := <-rowChannel:
	//			if !ok {
	//				rowChannel = nil
	//			} else {
	//				rows = append(rows, row)

	//				strRow = "[" + row.Columns[0].GetString_() + ":"
	//				strRow += row.Columns[1].GetString_() + ":"
	//				strRow += strconv.FormatUint(row.Columns[2].GetUint64(), 10) + "]"
	//				strResult += strRow
	//			}
	//		}
	//		if rowChannel == nil {
	//			break
	//		}
	//	}

	return []byte(strResult), nil
}

func main() {

	//	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(AssetManagementChaincode))
	if err != nil {
		myLogger.Debugf("Error starting AssetManagementChaincode: %s", err)
	}
}
