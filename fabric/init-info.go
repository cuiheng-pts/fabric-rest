package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/viper"
)

type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgAdmin       string
	OrgName        string
	OrdererOrgName string
	OrgResMgmt     *resmgmt.Client
	ChannelExist   bool
	ChannelClient  *channel.Client

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
	ChaincodeExist  bool
}

func ConstructorFromYaml(yamlPath string) (*InitInfo, error) {
	viper.SetConfigFile("YAML")
	viper.SetConfigName("info")
	viper.SetConfigFile(yamlPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("can not fild config YNAL file")
		} else {
			return nil, fmt.Errorf("fail to read yaml file: %v", err.Error())
		}
	}
	channel := viper.GetStringMapString("Channel")
	org := viper.GetStringMapString("Org")
	chaincode := viper.GetStringMapString("Chaincode")
	return &InitInfo{
		ChannelID:       channel["id"],
		ChannelConfig:   channel["config"],
		OrgAdmin:        org["admin"],
		OrgName:         org["name"],
		ChannelExist:    channel["Exist"] == "true",
		OrdererOrgName:  org["orderername"],
		ChaincodeID:     chaincode["id"],
		ChaincodeGoPath: chaincode["gopath"],
		ChaincodePath:   chaincode["path"],
		UserName:        org["user"],
		ChaincodeExist:  chaincode["Exist"] == "true",
	}, nil
}
