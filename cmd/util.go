package cmd

import (
	"fmt"
	"github.com/ajmyyra/go-epp-fi/pkg/registry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

func getRegistryClient(cmd *cobra.Command) (*registry.Client, error) {
	keyPath, err := getConfigString("CLIENT_KEY")
	if err != nil {
		return nil, err
	}
	clientKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load client key from " + keyPath)
	}

	certPath, err := getConfigString("CLIENT_CERT")
	if err != nil {
		return nil, err
	}
	clientCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load client certificate from " + certPath)
	}

	server, err := getConfigString("SERVER")
	if err != nil {
		return nil, err
	}

	username, err := getConfigString("USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := getConfigString("PASSWORD")
	if err != nil {
		return nil, err
	}

	client, err := registry.NewRegistryClient(username, password, server, 700, clientKey, clientCert)

	return client, nil
}

func getConfigString(name string) (string, error) {
	envName := fmt.Sprintf("%s_%s", envPrefix, name)

	value := viper.GetString(name)
	if value == "" {
		errMsg := fmt.Sprintf("Config variable %s (env: %s) not defined.", name, envName)
		return "", errors.New(errMsg)
	}

	return value, nil
}
