package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	Port        int    `env:"PORT" default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" required:"true"`
	JWTSecret   string `env:"JWT_SECRET" required:"true"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	val := reflect.ValueOf(cfg).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		key := field.Tag.Get("env")
		if key == "" {
			continue
		}

		raw := os.Getenv(key)

		// Required check
		if field.Tag.Get("required") == "true" && raw == "" {
			return nil, fmt.Errorf("%s is required but not set", key)
		}

		// Default value
		if raw == "" {
			raw = field.Tag.Get("default")
		}

		// Set field value
		if raw != "" {
			switch field.Type.Kind() {
			case reflect.String:
				val.Field(i).SetString(raw)
			case reflect.Int:
				n, err := strconv.Atoi(raw)
				if err != nil {
					return nil, fmt.Errorf("invalid value for %s: %v", key, err)
				}
				val.Field(i).SetInt(int64(n))
			}
		}
	}

	return cfg, nil
}
