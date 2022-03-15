/*
@Author : LiHongchao
@Description ：
@File : 
@Software: GoLand
@Version： 1.0.0
@Date : 2021/12/21 15:28
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

type BasicInfo struct{
	//个人基本信息ID
	UserID string `json:"UserID"`
	//身份证
	IDCard string `json:"IDCard"`
	//姓名
	Name string `json:"Name"`
	//性别
	Sex string `json:"Sex"`
	//民族
	Nation string `json:"Nation"`
	//籍贯
	Native string `json:"Native"`
	//生日
	Birthday string `json:"Birthday"`
	//手机
	Phone string `json:"Phone"`
	//邮箱
	Email string `json:"Email"`
	//政治面貌
	PoliticalLook string `json:"PoliticalLook"`
	//家庭住址
	HomeAddress string `json:"HomeAddress"`
	//信息上传账号
	LoginUserID string `json:"LoginUserID"`
	//个人信息文件hash
	UserInfoFileHash string `json:"UserInfoFileHash"`
	//插入该数据的时间
	Time string `json:"Time"`
	//TxId（交易ID） 项目后面增加的，写这篇文档时还没有增加
	TxId string `json:"TxId"`
}

type SmartContract struct{ //合约即链码
}

func (b *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response{
	return shim.Success(nil)
}

//提供函数入口
func (b *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addBasicInfo" { //添加个人基本信息
		stub.SetEvent("addBasicInfo", []byte{})
		return b.addBasicInfo(stub, args)
	}else if function == "queryBasicByUserID" { //根据UserID查询个人基本信息
		stub.SetEvent("queryBasicByUserID", []byte{})
		return b.queryBasicByUserID(stub,args)
	}else if function == "delBasicInfoByUserID"{ // 根据UserID删除个人基础信息
		stub.SetEvent("delBasicInfoByUserID", []byte{})
		return b.delBasicInfoByUserID(stub, args)
	}else if function == "updateBasicInfoByUserID"{ // 根据UserID更新个人基础信息
		stub.SetEvent("updateBasicInfoByUserID", []byte{})
		return b.updateBasicInfoByUserID(stub, args)
	}else if function == "getHistoryBasicInfo"{ // 查询历史个人信息
		stub.SetEvent("getHistoryBasicInfo", []byte{})
		return b.getHistoryBasicInfo(stub, args)
	}else if function == "queryBasicInfoByIDCard"{ // 富查询-根据身份证查询信息
		stub.SetEvent("queryBasicInfoByIDCard", []byte{})
		return b.queryBasicInfoByIDCard(stub, args)
	}

	return shim.Error("Invalid BasicInfo Chaincode Smart Contract function name.")
}


//添加个人基本信息(链码函数)
func (b *SmartContract) addBasicInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 13 {
		s_err := fmt.Sprintf("Incorrect number of arguments. Expecting 1, the length of args : %v,args:%v",len(args),args)
		return shim.Error(s_err)
	}
	insert_time :=time.Now().Format("2006-01-02 15:04:05")
	basicInfo := BasicInfo{
		//个人基本信息id
		UserID:args[0],
		//身份证
		IDCard:args[1],
		//姓名
		Name:args[2],
		//性别
		Sex:args[3],
		//民族
		Nation:args[4],
		//籍贯
		Native:args[5],
		//生日
		Birthday:args[6],
		//手机
		Phone:args[7],
		//邮箱
		Email:args[8],
		//政治面貌
		PoliticalLook:args[9],
		//家庭住址
		HomeAddress:args[10],
		//信息上传账号
		LoginUserID:args[11],
		//个人信息文件hash
		UserInfoFileHash:args[12],
		//插入该数据的时间
		Time:insert_time, //time.Now().Format("2006-01-02 15:04:05")
	}
	basicInfoAsBytes, _ := json.Marshal(basicInfo)
	stub.PutState(basicInfo.UserID, basicInfoAsBytes)
	fmt.Printf("\nAdd:%s\n\n", basicInfo)
	//txid:=stub.GetTxID()  //获取当前交易的tx_id
	return shim.Success(nil)
}


//根据用户ID进行查找人的基本信息(链码函数)
// args: UserID
func (b *SmartContract) queryBasicByUserID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//args传入函数多个参数，这里如果只传入UserID也可以不使用args []string 只用user_id string
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	basicInfoAsBytes, _ := stub.GetState(args[0])
	var basicInfo BasicInfo
	err := json.Unmarshal(basicInfoAsBytes, &basicInfo)
	if err != nil {
		return  shim.Error("反序列化basic_info信息失败")
	}
	//fmt.Printf("\nQuery:%s\n\n", basicInfo)
	return shim.Success(basicInfoAsBytes)
}


// 根据UserID删除个人信息
// args: UserID
func (b *SmartContract) delBasicInfoByUserID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	//根据userid查询信息
	basicInfoAsBytes, err:=stub.GetState(args[0])
	if err!=nil{
		return shim.Error("未找到要删除的信息")
	}
	var basicInfo BasicInfo
	err = json.Unmarshal(basicInfoAsBytes,&basicInfo)
	if err!=nil{
		return shim.Error("反序列化basicInfo失败")
	}
	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}
	fmt.Printf("\nDelete:%s\n\n", basicInfo)
	return shim.Success(nil)
}


//  根据UserID修改个人信息
//  args: UserID
func (b *SmartContract) updateBasicInfoByUserID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	//反序列化传入参数
	insertTime :=time.Now().Format("2006-01-02 15:04:05")
	info := BasicInfo{
		//个人基本信息id
		UserID:args[0],
		//身份证
		IDCard:args[1],
		//姓名
		Name:args[2],
		//性别
		Sex:args[3],
		//民族
		Nation:args[4],
		//籍贯
		Native:args[5],
		//生日
		Birthday:args[6],
		//手机
		Phone:args[7],
		//邮箱
		Email:args[8],
		//政治面貌
		PoliticalLook:args[9],
		//家庭住址
		HomeAddress:args[10],
		//信息上传账号
		LoginUserID:args[11],
		//个人信息文件hash
		UserInfoFileHash:args[12],
		//插入该数据的时间
		Time: insertTime, //time.Now().Format("2006-01-02 15:04:05")
	}

	//根据userid查询信息
	basicInfoAsBytes, err:=stub.GetState(args[0])
	if err!=nil{
		return shim.Error("未找到要修改的信息")
	}

	var basicInfo BasicInfo
	err = json.Unmarshal(basicInfoAsBytes,&basicInfo)
	if err!=nil{
		return shim.Error("反序列化basic_info失败")
	}

	basicInfo.IDCard=info.IDCard
	basicInfo.Name=info.Name
	basicInfo.Sex=info.Sex
	basicInfo.Nation=info.Nation
	basicInfo.Native=info.Native
	basicInfo.Birthday=info.Birthday
	basicInfo.Phone=info.Phone
	basicInfo.Email=info.Email
	basicInfo.PoliticalLook=info.PoliticalLook
	basicInfo.HomeAddress=info.HomeAddress
	basicInfo.LoginUserID=info.LoginUserID
	basicInfo.UserInfoFileHash=info.UserInfoFileHash
	basicInfo.Time=info.Time

	basicInfoAsBytes, err = json.Marshal(basicInfo)
	if err != nil {
		return shim.Error("序列化basic_info失败")
	}

	// 保存basic_info状态
	err = stub.PutState(basicInfo.UserID, basicInfoAsBytes)
	if err != nil {
		return shim.Error("更新basic_info失败")
	}
	fmt.Printf("\nChanged to:%s\n\n", basicInfo)
	return shim.Success(nil) //	return shim.Success([]byte("修改信息成功"))
}

func (b *SmartContract) getHistoryBasicInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[0]

	fmt.Printf("- start getHistoryBasicInfo: %s\n", userId)

	resultsIterator, err := stub.GetHistoryForKey(userId)
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

	fmt.Printf("- getHistoryBasicInfo returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (b *SmartContract) queryBasicInfoByIDCard(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
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


