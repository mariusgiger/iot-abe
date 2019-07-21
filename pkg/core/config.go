package core

// Config for iot-abe
type Config struct {
	EtherscanAPIKey string `yaml:"etherscanAPIKey"`
	EtherscanURL    string `yaml:"etherscanURL"`
	UseTestnet      bool   `yaml:"useTestnet"`
	ETHNodeURL      string `yaml:"ethNodeUrl"`
	ETHWSSNodeURL   string `yaml:"ethWssNodeUrl"`
	ClientCert      string `yaml:"clientCert"`
	ClientCertKey   string `yaml:"clientCertKey"`
	EthKeystoreDir  string `yaml:"ethKeystoreDir"`
	DataDir         string `yaml:"dataDir"`

	//Server Server `yaml:"server"`
	BuildTime string
	GitHash   string
	Version   string
}

//Server config for iot-abe
type Server struct {
}
