package web

type Configuration struct {
	Expose           bool `json:"expose"`
	Port             int  `json:"port"`
	RequestTimeoutMs int  `json:"request_timeout_ms"`
}

func ConfigurationDefaults() (*Configuration, error) {
	return &Configuration{
		Expose:           false,
		Port:             8080,
		RequestTimeoutMs: 30_000,
	}, nil
}
