package main

import (
	"os/exec"

	"github.com/pkg/errors"
)

type TarInstaller struct {
	archivePath string
}

func (t *TarInstaller) Install(dst string) error {
	cmd := exec.Command("tar", "-C", dst, "-xvf", t.archivePath)
	log.Info("unpacking archive via tar: ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to unpack archive %s", output)
	}

	return nil
}

func NewTarInstaller(archPath string) (*TarInstaller, error) {
	cmd := exec.Command("which", "tar")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "tar command not found %s", output)
	}

	return &TarInstaller{
		archivePath: archPath,
	}, nil
}
