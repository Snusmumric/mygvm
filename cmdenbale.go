package main

import "os/exec"

func cmdEnable(version string) {
	versionList := readBinDir()

	current := &goData{}
	versionToSet := &goData{}
	for _, goDataItem := range versionList {
		if goDataItem.version == version {
			*versionToSet = goDataItem
		}
		if goDataItem.isCurrent {
			*current = goDataItem
		}
	}
	if versionToSet.version == "" {
		log.Fatalf("desired version (%s) doesn't exist. %s",
			version, versionList)
	}
	if versionToSet.version == current.version {
		log.Infof("desired version is already active")
		return
	}

	// switch gobin symlink
	// ln -sfn HOME/gover/new_version/go HOME/gobin
	cmd := exec.Command("ln", "-sfn", versionToSet.GobinPath(), GobinSymlink)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to switch gobin symlink via '%s', error occured: %s; cmdout: %s",
			cmd.String(), err, output,
		)
	}

	// switch gopath symlink
	cmd = exec.Command("ln", "-sfn", versionToSet.GopathPath(), GopathSymlink)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to switch gopath symlink via '%s', error occured: %s; cmdout: %s",
			cmd.String(), err, output,
		)
	}

	// FIXME Think, can this made cheaper than checkConsistency
	*current = checkConsistency()

	log.Infof("Version successfully set to %s", current.version)
}
