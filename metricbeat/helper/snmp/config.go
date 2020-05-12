package snmp

import (
    "time"
)

// Config for an HTTP helper
type Config struct {
    Host      string             `config:"host"`
    Port      int                `config:"port"`
    Community string             `config:"community"`
    Timeout   time.Duration      `config:"timeout"`
    Version   string             `config:"version"`
    MaxOids   int                `config:"max_oids"`
}

func defaultConfig() Config {
    return Config{
        Host: "localhost",
        Port: 161,
        Community: "public",
        Version: "v2c",
        Timeout: time.Duration(2) * time.Second,
        MaxOids: 10,
    }
}
