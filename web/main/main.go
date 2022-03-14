package main

import (
	"fmt"
	"github.com/lihongchaoxe/blockchain-database/tree/master/web/controller"
	"github.com/lihongchaoxe/blockchain-database/tree/master/web/info"
	"github.com/lihongchaoxe/blockchain-database/tree/master/web/sdkInit"
	"github.com/lihongchaoxe/blockchain-database/tree/master/web/service"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func org1Start(sdk *fabsdk.FabricSDK)  {

	fmt.Println("开始Org1.......")
	//创建Org1通道
	err := sdkInit.CreateChannel(sdk)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !info.Initialized {
		//Org1加入通道
		err = sdkInit.PeerJoinChannel(sdk, info.Org1Name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//安装链码
		err = sdkInit.InstallCC()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//实例化链码
		err = sdkInit.InstantiateCC()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		info.Initialized = !info.Initialized
	}
}

func main() {

	//初始化SDK
	sdk, err := sdkInit.SetupSDK()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer sdk.Close()

	org1Start(sdk)
	
	channelClient, err := sdkInit.CreatChannelClient(sdk)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)
	info.Client = channelClient

	WebStart()
	
}

func  WebStart()  {

	router := gin.Default()
	router.GET("/", controller.Index)
	router.GET("/query/:ProgramID", controller.Query)
	router.GET("/queryhistory/:ProgramID", controller.QueryHistory)
	//router.POST("/transfer",controller.Transfer)
	router.POST("/delete/:ProgramID", controller.Delete)
	router.POST("/register", controller.Register)
	fmt.Println("启动Web服务, 监听端口号: 9000")
	router.Run(":9000")
	
}
