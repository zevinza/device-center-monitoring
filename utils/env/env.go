package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func Load(defaultValues map[string]interface{}) {
	LoadEnvironmentPrivate(defaultValues)
	LoadEnvironmentLocal()
	LoadEnvironmentSystem(defaultValues)
	LoadEnvironmentParameter(defaultValues)
	MergeAllEnvironment()
}

func LoadEnvironmentSystem(defaultValues map[string]interface{}) {
	systemEnv := viper.New()
	systemEnv.AutomaticEnv()
	for k := range defaultValues {
		keyUpper := strings.ToUpper(k)
		if value := systemEnv.Get(keyUpper); value != nil && value != "" {
			viper.Set(keyUpper, value)
		} else if value := systemEnv.Get(k); value != nil && value != "" {
			viper.Set(keyUpper, value)
		}
	}
}

func LoadEnvironmentLocal() {
	envFiles := []string{".env"}

	root := findProjectRoot()
	if root != "" {
		rootEnv := filepath.Join(root, ".env")
		if rootEnv != ".env" {
			envFiles = append(envFiles, rootEnv)
		}
	}

	for _, envFile := range envFiles {
		envMap, err := gotenv.Read(envFile)
		if err != nil {
			continue
		}

		_ = gotenv.Load(envFile)
		for k, v := range envMap {
			key := strings.ToUpper(k)
			viper.Set(key, v)
		}
		// Found and loaded a .env file, we're done
		return
	}

}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return ""
}

func LoadEnvironmentParameter(defaultValues map[string]interface{}) {
	paramEnv := viper.New()
	paramEnv.AllowEmptyEnv(false)
	for k := range defaultValues {
		if flagKey := strcase.ToKebab(k); nil == pflag.Lookup(flagKey) {
			pflag.String(flagKey, "", k)
		}
	}

	if os.Getenv("ENVIRONMENT_SIMULATION") != "" {
		for k := range defaultValues {
			pflag.CommandLine.Set(strcase.ToKebab(k), viper.GetString(k))
		}
	}

	pflag.Parse()
	if err := paramEnv.BindPFlags(pflag.CommandLine); nil == err {
		paramEnvKeys := paramEnv.AllKeys()
		for i := range paramEnvKeys {
			if stringValue := paramEnv.GetString(paramEnvKeys[i]); stringValue != "" {
				viper.Set(strcase.ToSnake(paramEnvKeys[i]), stringValue)
			}
		}
	}
}

func LoadEnvironmentPrivate(defaultValues map[string]interface{}) {
	for k, v := range defaultValues {
		key := strings.ToUpper(k)
		viper.SetDefault(key, v)
	}
}

func MergeAllEnvironment() {
	keys := viper.AllKeys()
	for i := range keys {
		stringValue := viper.GetString(keys[i])
		if stringValue == "" {
			if value := viper.Get(keys[i]); nil != value {
				stringValue = fmt.Sprintf("%v", value)
			}
		}
		os.Setenv(strings.ToUpper(keys[i]), stringValue)
	}
}
