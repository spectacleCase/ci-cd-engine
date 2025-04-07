package config

type System struct {
	HttpPort   string `yaml:"httpPort"`
	SaltLength int    `yaml:"saltLength"`
}
