package sdkInit

import (
	"fmt"
	"github.com/lihongchaoxe/blockchain-database/web/info"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"log"
)

func SetupSDK() (*fabsdk.FabricSDK, error) {

	sdk, err := fabsdk.New(config.FromFile(info.ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK初始化成功")
	return sdk, nil
}


func CreateResMgmtClient(sdk *fabsdk.FabricSDK, orgAdmin, orgName string) (*resmgmt.Client, error) {

	clientContext := sdk.Context(
		fabsdk.WithUser(orgAdmin),
		fabsdk.WithOrg(orgName))
	if clientContext == nil {
		return nil, fmt.Errorf(orgName + "创建资源管理客户端Context失败")
	}
	fmt.Println(orgName + "创建资源管理客户端成功...")

	// 返回资源管理客户端实例
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return nil, fmt.Errorf(orgName + "创建通道管理客户端失败: %v", err)
	}
	fmt.Println(orgName + "创建通道管理客户端成功...")
	return resMgmtClient, nil
}


func CreateChannel(sdk *fabsdk.FabricSDK) error {

	resMgmtClient, err := CreateResMgmtClient(sdk, info.OrgAdmin, info.OrdererOrgName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if !info.Initialized{
		fmt.Println("保存带有交易ID的通道")
		resMgmtClient, err = saveChannel(sdk, resMgmtClient)
		if err != nil {
			return fmt.Errorf("保存带有交易ID的通道: %v", err)
		}
	}else {
		fmt.Println("加入通道已完成...")
	}
	info.ResMgmtClient = resMgmtClient
	return nil
}


func saveChannel(sdk *fabsdk.FabricSDK, resMgmtClient *resmgmt.Client) (*resmgmt.Client,error) {

	mspClient, err := mspclient.New(
		sdk.Context(),
		mspclient.WithOrg(info.Org1Name))
	if err != nil {
		log.Panicf("创建 msp 客户端失败: %s", err)
	}

	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		log.Panicf("获取标识失败: %s", err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:info.ChannelID,
		ChannelConfigPath:info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	// 保存带有交易ID的渠道响应
	_, err = resMgmtClient.SaveChannel(channelReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererID))
	if err != nil {
		return nil,fmt.Errorf("创建应用通道失败: %v", err)
	}
	fmt.Println("通道创建成功...")
	return resMgmtClient,nil

}


func PeerJoinChannel(sdk *fabsdk.FabricSDK, orgName string) error {

	orgResMgmt, err := CreateResMgmtClient(sdk, info.OrgAdmin, orgName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
      reqPeers := resmgmt.WithTargetEndpoints("peer0.org2.cfginfo.com")
      err = info.OrgResMgmt.JoinChannel(info.ChannelID, reqPeers, resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
      if err != nil {
	      return fmt.Errorf("peer0.org2.cfginfo.com 加入通道失败: %v", err)
        }
      fmt.Println("peer0.org2.cfginfo.com 已成功加入通道...")
      reqPeers1 := resmgmt.WithTargetEndpoints("peer1.org2.cfginfo.com")
      err = info.OrgResMgmt.JoinChannel(info.ChannelID, reqPeers, resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
      if err != nil {
	      return fmt.Errorf("peer1.org2.cfginfo.com 加入通道失败: %v", err)
        }
      fmt.Println("peer1.org2.cfginfo.com 已成功加入通道...")
	err = orgResMgmt.JoinChannel(info.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererID))
	if  err != nil {
		return fmt.Errorf("节点加入通道失败: %v", err)
	}
	fmt.Println(orgName + " 加入通道成功...")
	info.OrgResMgmt = orgResMgmt
	return nil
}


func InstallCC() (error) {

	fmt.Println("开始安装链码......")
	ccPkg, err := gopackager.NewCCPackage(
		info.ChaincodePath,
		info.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v", err)
	}
	fmt.Println("创建链码包成功...")

	installCCReq := resmgmt.InstallCCRequest{
		Name: info.ChaincodeID,
		Path: info.ChaincodePath,
		Version:info.ChaincodeVersion,
		Package: ccPkg,
	}
      reqPeers := resmgmt.WithTargetEndpoints("peer0.org2.cfginfo.com")
      _, err = info.OrgResMgmt.InstallCC(installCCReq, reqPeers)
      if err != nil {
		return fmt.Errorf("安装链码失败: %v", err)
	}
	_, err = info.OrgResMgmt.InstallCC(
		installCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("安装链码失败: %v", err)
	}

	fmt.Println("指定链码安装成功")
	return nil
}


func InstantiateCC() (error) {

	fmt.Println("开始实例化链码......")

	org1OrOrg2 := "OR('Org1MSP.peer','Org2MSP.peer')"
	ccPolicy, err := policydsl.FromString(org1OrOrg2)
	if err != nil {
		return fmt.Errorf("获取策略失败！")
	}

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: info.ChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}
	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID,
		instantiateCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("实例化链码失败: %v", err)
	}
	fmt.Println("链码实例化成功")
	return nil
}


func CreatChannelClient(sdk *fabsdk.FabricSDK) (*channel.Client, error) {

	clientChannelContext := sdk.ChannelContext(
		info.ChannelID,
		fabsdk.WithUser("User1"))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}
	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}