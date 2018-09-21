package env

import "os"

func InitMapFromEnv(e map[string]string) map[string]string {
	for key, _ := range e {
		v := os.Getenv(key)
		if v != "" {
			e[key] = v
		}
	}
	return e
}
