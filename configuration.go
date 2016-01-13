package main

import (
	fmtlog "log"
	"time"

	//	"github.com/BurntSushi/toml"
	"github.com/emilevauge/traefik/provider"
	"github.com/emilevauge/traefik/types"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"os"
)

// GlobalConfiguration holds global configuration (with providers, etc.).
// It's populated from the traefik configuration file passed as an argument to the binary.
type GlobalConfiguration struct {
	Port                      string              `short:"p" long:"port" description:"Reverse proxy port" required:"false" default:":80"`
	GraceTimeOut              int64               `short:"g" long:"graceTimeOut" description:"Timeout in seconds. Duration to give active requests a chance to finish during hot-reloads" required:"false" default:"10"`
	AccessLogsFile            string              `long:"accessLogsFile" description:"Access logs file" required:"false"`
	TraefikLogsFile           string              `long:"traefikLogsFile" description:"Traefik logs file. If not defined, logs to stdout" required:"false"`
	Certificates              []Certificate       `group:"Certificates" namespace:"cert"`
	LogLevel                  string              `short:"l" long:"logLevel" description:"Log level" required:"false"  default:"ERROR"`
	ProvidersThrottleDuration time.Duration       `long:"providersThrottleDuration" description:"Backends throttle duration: minimum duration between 2 events from providers before applying a new configuration. It avoids unnecessary reloads if multiples events are sent in a short amount of time." required:"false"  default:"2s"`
	Docker                    *provider.Docker    `group:"Docker" namespace:"docker"`
	File                      *provider.File      `group:"File" namespace:"file"`
	Web                       *WebProvider        `group:"Web" namespace:"web"`
	Marathon                  *provider.Marathon  `group:"Marathon" namespace:"marathon"`
	Consul                    *provider.Consul    `group:"Consul" namespace:"consul"`
	Etcd                      *provider.Etcd      `group:"Etcd" namespace:"etcd"`
	Zookeeper                 *provider.Zookepper `group:"Zookeeper" namespace:"zookeeper"`
	Boltdb                    *provider.BoltDb    `group:"Boltdb" namespace:"boltdb"`
}

// Certificate holds a SSL cert/key pair
type Certificate struct {
	CertFile string `long:"certFile"`
	KeyFile  string `long:"keyFile"`
}

// NewGlobalConfiguration returns a GlobalConfiguration with default values.
func NewGlobalConfiguration() *GlobalConfiguration {
	globalConfiguration := new(GlobalConfiguration)
	// default values
	globalConfiguration.Port = ":80"
	globalConfiguration.GraceTimeOut = 10
	globalConfiguration.LogLevel = "ERROR"
	globalConfiguration.ProvidersThrottleDuration = time.Duration(2 * time.Second)

	return globalConfiguration
}

// LoadFileConfig returns a GlobalConfiguration from reading the specified file (a toml file).
func LoadFileConfig(file string) *GlobalConfiguration {
	configuration := NewGlobalConfiguration()
	//	if _, err := toml.DecodeFile(file, configuration); err != nil {
	//		fmtlog.Fatalf("Error reading file: %s", err)
	//	}
	viper.SetEnvPrefix("traefik")
	viper.AutomaticEnv()
	viper.SetConfigName("traefik")        // name of config file (without extension)
	viper.AddConfigPath("/etc/traefik/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.traefik") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		fmtlog.Fatalf("Error reading file: %s", err)
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmtlog.Fatalf("Error reading file: %s", err)
	}

	argsConfiguration := NewGlobalConfiguration()
	parser := flags.NewParser(argsConfiguration, flags.Default)
	parser.WriteHelp(os.Stdout)
	return configuration
}

type configs map[string]*types.Configuration
