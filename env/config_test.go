package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Should use mongodb://127.0.0.1:27017 as the default MONGO_URI", func(t *testing.T) {
		err := os.Setenv("MONGO_URI", "")

		if err != nil {
			t.FailNow()
		}

		config := createConfig()

		require.Equal(t, "mongodb://127.0.0.1:27017", config.MongoURI)
	})

	t.Run("Should define the config by env params correctly", func(t *testing.T) {
		os.Setenv("MONGO_URI", "mongodb://mongo_host:27017")
		os.Setenv("SERVER_HOST", "http://unit.test.io")

		config := createConfig()

		require.Equal(t, "mongodb://mongo_host:27017", config.MongoURI)
		require.Equal(t, "http://unit.test.io", config.ServerHost)
	})
}
