package apiserver

type Config struct {
	Debug       bool   `yaml:"debug"`
	ListenAddr  string `yaml:"listenAddr"`
	LogInfoPath string `yaml:"logInfoPath"`
	LogErrPath  string `yaml:"logErrPath"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Error string      `json:"error"`
}
