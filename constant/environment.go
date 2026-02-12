package constant

import "github.com/spf13/viper"

const (
	EnvironmentLocal       = "local"
	EnvironmentDevelopment = "development"
	EnvironmentStaging     = "staging"
	EnvironmentProduction  = "production"
)

func GetEnvironment() string {
	return viper.GetString("ENVIRONMENT")
}
