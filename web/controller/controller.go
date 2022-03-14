package controller

import (
	"github.com/lihongchaoxe/blockchain-database/web/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(context *gin.Context) {
	context.String(http.StatusOK, "Welcome！")
}

func Query(context *gin.Context) {
	ProgramID := context.Param("ProgramID")
	msg, err := service.Query(ProgramID)
	if err != nil {
		msg = "没有查询到对应的信息"
	}
	context.JSON(http.StatusOK, gin.H{
		"ProgramID": ProgramID,
		"Message": msg,
	})
}
func QueryHistory(context *gin.Context) {
	ProgramID := context.Param("ProgramID")
	msg, err := service.QueryHistory(ProgramID)
	if err != nil {
		msg = "没有查询到对应的信息"
	}
	context.JSON(http.StatusOK, gin.H{
		"ProgramID": ProgramID,
		"Message": msg,
	})
}

/*func Transfer(context *gin.Context)  {
	var msg string
	transferer := context.Query("name1")
	beneficiary := context.Query("name2")
	count := context.Query("count")
	transactionID, err := service.Transfer(transferer, beneficiary, count)
	if err != nil {
		msg = err.Error()
	}else {
		msg = transactionID
	}
	context.JSON(http.StatusOK, gin.H{
		"Transferer": transferer,
		"Beneficiary": beneficiary,
		"Count": count,
		"Message": msg,
	})
}*/

func Delete(context *gin.Context)  {
	var msg string
	ProgramID := context.Param("ProgramID")
	transactionID, err := service.Delete(ProgramID)
	if err != nil {
		msg = err.Error()
	}else {
		msg = transactionID
	}
	context.JSON(http.StatusOK, gin.H{
		"ProgramID": ProgramID,
		"Message": msg,
	})
}

func Register(context *gin.Context)  {
	var msg string
	ProgramID := context.Query("ProgramID")
	Input := context.Query("Input")
	CfgHash := context.Query("CfgHash")
	transactionID, err := service.Register(ProgramID, Input, CfgHash)
	if err != nil {
		msg = err.Error()
	}else {
		msg = transactionID
	}
	context.JSON(http.StatusOK, gin.H{
		"ProgramID": ProgramID,
		"Input": Input,
		"CfgHash": CfgHash,
		"Message": msg,
	})
}
