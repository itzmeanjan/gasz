package config

import "github.com/spf13/viper"

// Load - Given path to config file, reads file content into memory
func Load(file string) error {
	viper.SetConfigFile(file)

	return viper.ReadInConfig()
}

// Get - Given key present in config file, returns associated value
// given, it has been loaded in to memory, by calling 👆 function
func Get(key string) string {
	return viper.GetString(key)
}
