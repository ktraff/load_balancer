package load_balancer

import (
	"fmt"
	"net/url"
	"os"
)

type Configuration struct {
	backends *[]url.URL
}

func Config() Configuration {
	return Configuration{
		backends: GetBackends(),
	}
}

// Returns a list of available backend API endpoints.
func GetBackends() *[]url.URL {
	var backends []url.URL
	for i := 1; i > 0; i++ {
		backend := os.Getenv(fmt.Sprintf("BACKEND_%v", i))
		url, err := url.Parse(backend)
		if backend == "" || err != nil {
			break
		}
		fmt.Println(fmt.Sprintf("Configuring backend %q", backend))
		backends = append(backends, *url)
	}
	return &backends
}

// Checks for an existing environment variable.  If none exists,
// return a default.
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
