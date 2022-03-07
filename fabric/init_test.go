package fabric

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const configPath = "../config.yaml"
const infoPath = "../info.yaml"

var info = &InitInfo{
	ChannelID:       "sdktestchannel",
	ChannelConfig:   os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric/multipeer/channel-artifacts/SDKchannel.tx",
	OrgAdmin:        "Admin",
	OrgName:         "Org1",
	OrdererOrgName:  "orderer.example.com",
	ChaincodeID:     "sdktestcc",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/hyperledger/fabric/multipeer/chaincode/go/log/",
	UserName:        "User1",
}

func TestInitSDK(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
}

func TestCreateChannel(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	err = CreateChannel(sdk, info)
	require.NoError(t, err)
}

func TestJoinChannel(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	err = JoinChannel(sdk, info)
	require.NoError(t, err)
}

func TestInstallCC(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	clientContext := sdk.Context(fabsdk.WithOrg(info.OrgName), fabsdk.WithUser(info.OrgAdmin))
	require.NotNil(t, clientContext)
	resmgmtClient, err := resmgmt.New(clientContext)
	require.NoError(t, err)
	info.OrgResMgmt = resmgmtClient
	err = InstallCC(sdk, info)
	require.NoError(t, err)
}

func TestInstantiateCC(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	clientContext := sdk.Context(fabsdk.WithOrg(info.OrgName), fabsdk.WithUser(info.OrgAdmin))
	require.NotNil(t, clientContext)
	resmgmtClient, err := resmgmt.New(clientContext)
	require.NoError(t, err)
	info.OrgResMgmt = resmgmtClient
	err = InstantiateCC(sdk, info)
	require.NoError(t, err)
}

func TestGetChannelClient(t *testing.T) {
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	channelClient, err := GetChannelClient(sdk, info)
	require.NoError(t, err)
	require.NotNil(t, channelClient)
	request := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn:         "query",
		Args:        [][]byte{[]byte("123")},
	}
	response, err := channelClient.Query(request, channel.WithRetry(retry.DefaultResMgmtOpts))
	require.NoError(t, err)
	require.NotNil(t, response)
	t.Log(string(response.Payload))
}

func TestConstructorFromYaml(t *testing.T) {
	info, err := ConstructorFromYaml(infoPath)
	require.NoError(t, err)
	require.NotNil(t, info)
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	channelClient, err := GetChannelClient(sdk, info)
	require.NoError(t, err)
	require.NotNil(t, channelClient)
	request := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn:         "query",
		Args:        [][]byte{[]byte("123")},
	}
	response, err := channelClient.Query(request, channel.WithRetry(retry.DefaultResMgmtOpts))
	require.NoError(t, err)
	require.NotNil(t, response)
	t.Log(string(response.Payload))
}

func TestQuery(t *testing.T) {
	info, err := ConstructorFromYaml(infoPath)
	require.NoError(t, err)
	require.NotNil(t, info)
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	query, err := Query(info, "123")
	require.NoError(t, err)
	require.NotNil(t, query)
	t.Log(query)
}

func TestCreate(t *testing.T) {
	info, err := ConstructorFromYaml(infoPath)
	require.NoError(t, err)
	require.NotNil(t, info)
	sdk, err := InitSDK(configPath, false, info)
	require.NoError(t, err)
	require.NotNil(t, sdk)
	defer sdk.Close()
	err = Create(info, "234", "hello")
	require.NoError(t, err)
	query, err := Query(info, "234")
	require.NoError(t, err)
	require.NotNil(t, query)
	t.Log(query)
}
