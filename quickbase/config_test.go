package quickbase

import (
	"testing"

	"github.com/spf13/viper"
)

func newTestConfig() StandardConfig {
	return StandardConfig{viper: viper.New()}
}

func TestConfig(t *testing.T) {
	//	t.Fatal("Failing test to remind us to write them.")
}
