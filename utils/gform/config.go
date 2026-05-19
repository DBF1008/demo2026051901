package gform

// Config ...
type Config struct {
	Driver string `json:"driver"`
	Dsn             string `json:"dsn"`
	SetMaxOpenConns int    `json:"setMaxOpenConns"`
	SetMaxIdleConns int    `json:"setMaxIdleConns"`
	Prefix          string `json:"prefix"`
}

// ConfigCluster ...
type ConfigCluster struct {
	Master []Config
	Slave  []Config
	Driver string
	Prefix string
}
