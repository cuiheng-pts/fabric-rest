package fabric

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

const infoConfigPath = "../info.yaml"

func TestViperGet(t *testing.T) {
	viper.SetConfigFile("YAML")
	viper.SetConfigName("info")
	viper.SetConfigFile(infoConfigPath)
	err := viper.ReadInConfig()
	require.NoError(t, err)
	org := viper.GetStringMapString("Org")
	require.NotNil(t, org)
	t.Log(org)
	t.Log(reflect.TypeOf(org))
}
