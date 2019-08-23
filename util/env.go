package util

import (
	"os"
	"strings"
)

// GetEnvMap returns a map of the Environment variables.
func GetEnvMap() map[string]string {
	var values = make(map[string]string)

	for _, env := range os.Environ() {
		kv := strings.Split(env, "=")

		values[kv[0]] = kv[1]
	}

	return values
}
