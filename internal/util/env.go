package util

import "os"

func GetEnv(keys ...string) string {
	for _, k := range keys {
		v := os.Getenv(k)
		if v != "" {
			return v
		}
	}
	return ""
}
