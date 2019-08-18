package env

import "os"

type config struct {
	MongoURI   string
	ServerHost string
}

// Config represents the environment variables this project uses
var Config config

func init() {
	Config = createConfig()
}

func createConfig() config {
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		mongoURI = "mongodb://127.0.0.1:27017"
	}

	return config{
		MongoURI:   mongoURI,
		ServerHost: os.Getenv("SERVER_HOST"),
	}
}
