package web

type Configuration struct {
	Expose bool `json:"expose"`
	Port   int  `json:"port"`
}

func ConfigurationDefaults() (*Configuration, error) {
	return &Configuration{
		Expose: false,
		Port:   8080,
	}, nil
}
