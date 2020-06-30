package main

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func main() {

	// Find home directory.
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Error Getting Home Directory", err)
		return
	}

	viper.SetDefault("project-name", "OTA")
	viper.AddConfigPath(homeDir)
	viper.SetConfigName(".ota_packer_test")
	viper.SetConfigType("yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		viper.Set("project-name", "xxxx")
		viper.WriteConfig()
	} else {
		viper.Set("project-name", "yy")
		// Everything marked with safe won't overwrite any file,
		// but just create if not existent
		viper.SafeWriteConfig()
	}
}
