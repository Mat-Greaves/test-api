package config

// interface describing what the config looks like
type Configer interface {
	App() Apper
}

// config instance
type Config struct {
	AppSetting App
}

// method on config to return app config
func (c Config) App() Apper {
	return c.AppSetting
}

// return config instance
func GetConfig() Configer {
	return &Config{}
}

// Interface for application config to satisfy
type Apper interface {
	GetPort() string
}

// Struct to hang application settings methods off of
type App struct {
}

// method on App struct to fetch some config
func (a App) GetPort() string {
	return ":8080"
}
