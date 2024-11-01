package config

type Configurations struct {
    Server LLMServerConfig
    Verify VerificationServer
}

type LLMServerConfig struct {
    Port int
}

type VerificationServer struct {
    UuidFile string `mapstructure:"uuidfile"`
    VerifyUrl string
}
