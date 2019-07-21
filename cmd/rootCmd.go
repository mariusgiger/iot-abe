package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	log        = logrus.New()
	configPath string
	cfg        = &core.Config{}
	keysPath   string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "iot-abe",
	Short: "iot-abe - attribute-based access control for the IoT",
	Long:  `iot-abe - attribute-based access control for the IoT.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf("Something somewhere went terribly wrong:\n %v\n\n", err)
		os.Exit(-1)
	}
}

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	RootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", "./config.yml", "config path")

	configPath = NormalizePath(configPath)
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(configFile), cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cfg.DataDir = NormalizePath(cfg.DataDir)

	if _, err := os.Stat(cfg.DataDir); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.DataDir, os.ModePerm) //TODO change perms
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else if err != nil {
		log.Fatalf("error: %v", err)
	}

	keysPath = path.Join(cfg.DataDir, "keys")
	err = os.MkdirAll(keysPath, os.ModePerm) //TODO change perms
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cfg.EthKeystoreDir = NormalizePath(cfg.EthKeystoreDir)
	cfg.BuildTime = BuildTime
	cfg.Version = Version
	cfg.GitHash = GitHash
}

// NormalizePath expands a path to its absolute representation
func NormalizePath(path string) string {
	expandedPath, err := homedir.Expand(path)
	if err != nil {
		log.Fatalf("could not expand path, %v", err)
	}

	path, err = filepath.Abs(expandedPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return path
}
