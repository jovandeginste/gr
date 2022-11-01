package app

type Configuration struct {
	Retry         bool   `mapstructure:"retry"`
	RootDirectory string `mapstructure:"root_directory"`
	Shell         string `mapstructure:"shell"`
}
