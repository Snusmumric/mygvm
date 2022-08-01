package main

import (
	"os"
)

type goInstaller interface {
	Install(installDir string) error
}

func cmdInstall(version string, installer goInstaller) {
	goToInstall := goData{}
	err := goToInstall.SetVersion(version)
	if err != nil {
		log.Fatalf("installation error: invalid version string furnished: (%s), %s", version, err)
	}
	installPath := goToInstall.VersionPath()
	err = os.MkdirAll(installPath, 0755)
	if err != nil {
		log.Fatalf("installation error: failed to create directory %s", err)
	}

	err = installer.Install(installPath)
	if err != nil {
		log.Fatalf("installation error: go version %s; %s", version, err)
	}

	ver, err := checkVersionOfBin(goToInstall.BinaryPath())
	if err != nil {
		log.Fatalf("installation error: failed to check binary version (path: %s); %s",
			goToInstall.BinaryPath(), err,
		)
	}
	if ver != goToInstall.version {
		//TODO remove installed data; write to reuse some part of code in ordinary uninstallation
		log.Fatalf("installation error: installed version %s not equal to desired %s", ver, goToInstall.version)
	}

	log.Infof("version '%s' successfully installed into %s", version, goToInstall.GobinPath())
}
