package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)
import "github.com/hyperledger/fabric-sdk-go/pkg/core/config"

const DefaultChaincodeVersion = "1.0"

func InitSDK(configPath string, initialized bool, info *InitInfo) (*fabsdk.FabricSDK, error) {
	if initialized {
		return nil, fmt.Errorf("SDK has already been initialized")
	}
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sdk from file: %v, and error is %v", configPath, err.Error())
	}
	fmt.Println("successfully initialize sdk")
	clientContext := sdk.Context(fabsdk.WithOrg(info.OrgName), fabsdk.WithUser(info.OrgAdmin))
	if clientContext == nil {
		return nil, fmt.Errorf("failed to generate resource management context with OrgAdmin and OrgName")
	}
	resmgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manage client with given context: %v", err.Error())
	}
	info.OrgResMgmt = resmgmtClient
	if !info.ChannelExist {
		CreateChannel(sdk, info)
		JoinChannel(sdk, info)
	}
	if !info.ChaincodeExist {
		InstallCC(sdk, info)
		InstantiateCC(sdk, info)
	}
	info.ChannelClient, err = GetChannelClient(sdk, info)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	clientContext := sdk.Context(fabsdk.WithOrg(info.OrgName), fabsdk.WithUser(info.OrgAdmin))
	if clientContext == nil {
		return fmt.Errorf("failed to generate resource management context with OrgAdmin and OrgName")
	}
	resmgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("failed to create resource manage client with given context: %v", err.Error())
	}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("failed to create MSP client with given OrgName: %v", err.Error())
	}
	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("failed to get signature with given Admin: %v", err.Error())
	}
	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	_, err = resmgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("failed to create channel: %v", err.Error())
	}
	fmt.Println("successfully create channel")
	return nil
}

func JoinChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	clientContext := sdk.Context(fabsdk.WithOrg(info.OrgName), fabsdk.WithUser(info.OrgAdmin))
	if clientContext == nil {
		return fmt.Errorf("failed to generate resource management context with OrgAdmin and OrgName")
	}
	resmgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("failed to create resource manage client with given context: %v", err.Error())
	}
	info.OrgResMgmt = resmgmtClient
	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("failed to join in channel: %v", err.Error())
	}
	fmt.Println("successfully join in channel")
	return nil
}

func InstallCC(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	fmt.Println("start to install chaincode")
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("failed to create chaincode package: %v", err.Error())
	}
	installCCReq := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: DefaultChaincodeVersion,
		Package: ccPkg,
	}
	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("fail to install chaincode: %v", err.Error())
	}
	fmt.Println("successfully install chaincode")
	return nil
}

func InstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	fmt.Println("start to instantiate chaincode...")
	ccPolicy := policydsl.SignedByAnyMember([]string{"org1.example.com", "Org1MSP"})
	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: DefaultChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}
	_, err := info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("fail to instantiate chaincode: %v", err.Error())
	}
	fmt.Println("successfully instantiate chaincode")
	return nil
}

func GetChannelClient(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("fail to create channel client")
	}
	fmt.Println("successfully create channel client, use this client to invoke chaincode")
	return channelClient, nil
}
