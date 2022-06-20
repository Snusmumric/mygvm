package main

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type TarInstaller struct {
	archivePath string
}

func (t *TarInstaller) Install(dst string) error {
	cmd := exec.Command(fmt.Sprintf("tar -C %s -xvf %s", dst, t.archivePath))
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to unpack archive")
	}

	return nil
}

func NewTarInstaller(archPath string) (*TarInstaller, error) {
	cmd := exec.Command("which tar")
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "tar command not found")
	}

	return &TarInstaller{
		archivePath: archPath,
	}, nil
}
