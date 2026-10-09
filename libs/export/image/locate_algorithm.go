package image

import (
	"os"
	"os/exec"
)

// Locate finds the freeze binary: $CUPPA_FREEZE first, then the PATH.
func Locate() (string, error) {
	return locate(os.Getenv, exec.LookPath)
}

func locate(getenv func(string) string, lookPath func(string) (string, error)) (string, error) {
	if p := getenv(EnvFreeze); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		return "", ErrFreezeMissing
	}
	p, err := lookPath("freeze")
	if err != nil {
		return "", ErrFreezeMissing
	}
	return p, nil
}
