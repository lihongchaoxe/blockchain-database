package info

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"os"
)

var (
	ChannelID = "mychannel"
	ChannelConfig = os.Getenv("GOPATH") + "/src/blockchain-database/fixtures/channel-artifacts/channel.tx"
	OrdererID = "orderer0.cfginfo.com"
	OrgAdmin = "Admin"
	OrdererOrgName = "OrdererOrg"
	Org1Name = "Org1"
	Org2Name = "Org2"
	ChaincodeGoPath	= os.Getenv("GOPATH")
	ChaincodePath = "blockchain-database/chaincode/controlflow-info"
	UserName = "User1"
	ChaincodeID = "mycc"
	ConfigFile = "../config/config.yaml"
	Initialized = false
	Installed = false
	ChaincodeVersion  = "1.0"
	ResMgmtClient *resmgmt.Client
	OrgResMgmt *resmgmt.Client
	Client *channel.Client
)
