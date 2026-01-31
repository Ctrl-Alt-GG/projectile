package utils

import (
	"slices"
	"strings"

	"github.com/mitchellh/go-ps"
)

func PSGrep(exeName string) (bool, error) {
	processes, err := ps.Processes()
	if err != nil {
		return false, err
	}
	found := slices.ContainsFunc(processes, func(proc ps.Process) bool {
		return strings.Contains(proc.Executable(), exeName)
	})
	return found, nil
}
