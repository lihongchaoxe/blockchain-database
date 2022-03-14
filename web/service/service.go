package service

import (
	"github.com/lihongchaoxe/blockchain-database/web/info"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
)

func Query(ProgramID string) (string, error){
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn: "queryBasicByProgramID",
		Args: [][]byte{
			[]byte(ProgramID),
		},
	}
	respone, err := info.Client.Query(req)
	if err != nil {
		return "", errors.WithMessage(err, "查询链码失败！")
	}
	return string(respone.Payload), nil
}
func QueryHistory(ProgramID string) (string, error){
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn: "getHistoryCFInfo",
		Args: [][]byte{
			[]byte(ProgramID),
		},
	}
	respone, err := info.Client.Query(req)
	if err != nil {
		return "", errors.WithMessage(err, "查询链码失败！")
	}
	return string(respone.Payload), nil
}
/*func Transfer(name1, name2, num string) (string, error) {
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn: "transfer",
		Args: [][]byte{
			[]byte(name1),
			[]byte(name2),
			[]byte(num),
		},
	}
	respone, err := info.Client.Execute(req)
	if err != nil {
		return "", errors.WithMessage(err, "调用链码失败！")
	}
	return string(respone.TransactionID), nil
}*/

func Register(ProgramID, Input, CfgHash string) (string, error) {
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn: "addCFInfo",
		Args: [][]byte{
			[]byte(ProgramID),
			[]byte(Input),
			[]byte(CfgHash),
		},
	}
	respone, err := info.Client.Execute(req)
	if err != nil {
		return "", errors.WithMessage(err, "调用链码失败！")
	}
	return string(respone.TransactionID), nil
}

func Delete(ProgramID string) (string, error) {
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn: "delCFInfoByProgramID",
		Args: [][]byte{
			[]byte(ProgramID),
		},
	}
	respone, err := info.Client.Execute(req)
	if err != nil {
		return "", errors.WithMessage(err, "调用链码失败！")
	}
	return string(respone.TransactionID), nil
}
