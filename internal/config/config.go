package config

import (
	"os"
	"strconv"
	"time"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := env(key, "")
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envDur(key string, def time.Duration) time.Duration {
	v := env(key, "")
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func CheckScopes() bool {
	return envBool("TAPE_CHECK_SCOPES", true)
}

func CheckAssertionsCount() bool {
	return envBool("TAPE_CHECK_ASSERTIONS_COUNT", true)
}

func CheckEnd() bool {
	return envBool("TAPE_CHECK_END", true)
}

func Timeout() time.Duration {
	return envDur("TAPE_TIMEOUT", 3*time.Second)
}

func StrictTransitions() bool {
	return envBool("TAPE_STRICT_TRANSITIONS", true)
}
