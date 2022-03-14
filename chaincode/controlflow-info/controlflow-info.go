/*
@Author : LiHongchao
@Description ：
@File : 
@Software: 
@Version： 1.0.0
@Date : 2021/12/22 22:44
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
	"time"
)

type CFInfo struct{
	//程序ID
	ProgramID string `json:"ProgramID"`
	//输入
	Input string `json:"Input"`
	//控制流哈希值
	CfgHash string `json:"CfgHash"`
	//插入该数据的时间
	Time string `json:"Time"`
}

type SmartContract struct{ //合约即链码
}

func (b *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response{
	return shim.Success(nil)
}

//提供函数入口
func (b *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addCFInfo" { //添加程序控制流
		stub.SetEvent("addCFInfo", []byte{})
		return b.addCFInfo(stub, args)
	}else if function == "queryBasicByProgramID" { //根据ProgramID查询程序控制流信息
		stub.SetEvent("queryBasicByProgramID", []byte{})
		return b.queryBasicByProgramID(stub,args)
	}else if function == "delCFInfoByProgramID"{ // 根据ProgramID删除
		stub.SetEvent("delCFInfoByProgramID", []byte{})
		return b.delCFInfoByProgramID(stub, args)
	}else if function == "updateCFInfoByProgramID"{ // 根据ProgramID更新
		stub.SetEvent("updateCFInfoByProgramID", []byte{})
		return b.updateCFInfoByProgramID(stub, args)
	}else if function == "getHistoryCFInfo"{ // 查询历史程序控制流信息
		stub.SetEvent("getHistoryCFInfo", []byte{})
		return b.getHistoryCFInfo(stub, args)
	}/*else if function == "queryCFInfoByIDCard"{ // 富查询-根据身份证查询信息
		stub.SetEvent("queryCFInfoByIDCard", []byte{})
		return b.queryCFInfoByIDCard(stub, args)
	}*/

	return shim.Error("Invalid CFInfo Chaincode Smart Contract function name.")
}


//添加程序控制流信息(链码函数)
func (b *SmartContract) addCFInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		s_err := fmt.Sprintf("Incorrect number of arguments. Expecting 1, the length of args : %v,args:%v",len(args),args)
		return shim.Error(s_err)
	}
	insert_time :=time.Now().Format("2006-01-02 15:04:05")
	CFInfo := CFInfo{
		//程序id
		ProgramID:args[0],
		//输入
		Input:args[1],
		//控制流哈希值
		CfgHash:args[2],
		//插入该数据的时间
		Time:insert_time, //time.Now().Format("2006-01-02 15:04:05")
	}
	CFInfoAsBytes, _ := json.Marshal(CFInfo)
	stub.PutState(CFInfo.ProgramID, CFInfoAsBytes)
	fmt.Printf("\nAdd:%s\n\n", CFInfo)
	//txid:=stub.GetTxID()  //获取当前交易的tx_id
	return shim.Success(nil)
}


//根据程序ID查询程序控制流信息(链码函数)
// args: ProgramID
func (b *SmartContract) queryBasicByProgramID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//args传入函数多个参数，这里如果只传入ProgramID也可以不使用args []string 只用UserID string
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	CFInfoAsBytes, _ := stub.GetState(args[0])
	var CFInfo CFInfo
	err := json.Unmarshal(CFInfoAsBytes, &CFInfo)
	if err != nil {
		return  shim.Error("反序列化cf_info信息失败")
	}
	//fmt.Printf("\nQuery:%s\n\n", CFInfo)
	return shim.Success(CFInfoAsBytes)
}


// 根据ProgramID删除程序控制流信息
// args: ProgramID
func (b *SmartContract) delCFInfoByProgramID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	//根据ProgramID查询信息
	CFInfoAsBytes, err:=stub.GetState(args[0])
	if err!=nil{
		return shim.Error("未找到要删除的信息")
	}
	var CFInfo CFInfo
	err = json.Unmarshal(CFInfoAsBytes,&CFInfo)
	if err!=nil{
		return shim.Error("反序列化CFInfo失败")
	}
	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}
	fmt.Printf("\nDelete:%s\n\n", CFInfo)
	return shim.Success(nil)
}


//  根据ProgramID修改程序控制流信息
//  args: ProgramID
func (b *SmartContract) updateCFInfoByProgramID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		s_err := fmt.Sprintf("Incorrect number of arguments. Expecting 1, the length of args : %v,args:%v",len(args),args)
		return shim.Error(s_err)
	}
	insert_time :=time.Now().Format("2006-01-02 15:04:05")
	info := CFInfo{
		//程序id
		ProgramID:args[0],
		//输入
		Input:args[1],
		//控制流哈希值
		CfgHash:args[2],
		//插入该数据的时间
		Time:insert_time, //time.Now().Format("2006-01-02 15:04:05")
	}
	//根据ProgramID查询信息
	CFInfoAsBytes, err:=stub.GetState(args[0])
	if err!=nil{
		return shim.Error("未找到要修改的信息")
	}
	CFInfoAsBytes, err = json.Marshal(info)
	if err != nil {
		return shim.Error("序列化cfg_info失败")
	}

	// 保存cfg_info状态
	err = stub.PutState(info.ProgramID, CFInfoAsBytes)
	if err != nil {
		return shim.Error("更新cfg_info失败")
	}
	fmt.Printf("\nChanged to:%s\n\n", info)
	return shim.Success(nil) //	return shim.Success([]byte("修改信息成功"))
}

func (b *SmartContract) getHistoryCFInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProgramID := args[0]

	fmt.Printf("- start getHistoryCFInfo: %s\n", ProgramID)

	resultsIterator, err := stub.GetHistoryForKey(ProgramID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryCFInfo returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (b *SmartContract) queryCFInfoByIDCard(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	IDCard := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"IDCard\":\"%s\"}}", IDCard)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned asKey a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
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

	return &buffer, nil
}



func main(){
	err := shim.Start(new(SmartContract))
	if err != nil{
		fmt.Printf("启动SmartContract时发生错误: %s\n", err)
	}
}


