package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxyAddress    string
	ServerAddresses []string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	proxyAddress := strings.TrimSpace(os.Getenv("PROXY_ADDRESS"))
	if proxyAddress == "" {
		return nil, errors.New("no proxy address provided")
	}
	serverAddresses := strings.TrimSpace(os.Getenv("SERVER_ADDRESSES"))
	if serverAddresses == "" {
		return nil, errors.New("no server address provided")
	}
	serverAddressesList := strings.Split(serverAddresses, ",")
	for i := range serverAddressesList {
		serverAddressesList[i] = strings.TrimSpace(serverAddressesList[i])
	}

	return &Config{proxyAddress, serverAddressesList}, nil
}
