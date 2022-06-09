package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// HOME/gobin -> HOME/gover/go1.16.15/go           | GobinSymlink
// HOME/go    -> HOME/gover/go1.16.15/gopath       | GopathSymlink

var (
	HomeDir       = os.Getenv("HOME")
	GopathSymlink = path.Join(HomeDir, "go")
	GobinSymlink  = path.Join(HomeDir, "gobin")
	GoBinaryDir   = path.Join(GobinSymlink, "bin")
	GoBinaryPath  = path.Join(GoBinaryDir, "go")

	VersionDir = path.Join(HomeDir, "gover")
	versionRE  = regexp.MustCompile(`^go([1-9.]+)$`)
	gopathBase = "gopath"
	gobinBase  = "go"
)

type goData struct {
	version string
	major   uint64
	minor   uint64
	patch   uint64

	isCurrent bool
}

func (d *goData) ValidateVersion(version string) error {
	matchedIdxs := versionRE.FindAllStringIndex(version, -1)
	if len(matchedIdxs) != 1 {
		return errors.Errorf("wrong format of go version string, doesn't match re %s", versionRE)
	}
	return nil
}

func (d *goData) SetVersion(version string) (err error) {
	err = d.ValidateVersion(version)
	if err != nil {
		return err
	}
	d.version = version

	parts := strings.Split(d.version[2:], ".")
	if len(parts) > 3 {
		return errors.Errorf("wrong format of string go version '%s' (too many parts of version ints)")
	}

	d.major, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse major go version")
	}
	if len(parts) == 1 {
		return nil
	}
	d.minor, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse minor go version")
	}
	if len(parts) == 2 {
		return nil
	}
	d.patch, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "failed to parse patch go version")
	}

	return nil
}

func validateGolink(linkname, linkvalue, expectedBase string) (resp goData, err error) {
	stat, err := os.Lstat(linkvalue)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to get %s stat", linkname)
	}
	if stat.Mode()&fs.ModeSymlink == 0 {
		return resp, errors.Errorf("%s file %s is not a symlink, (%s)", linkname, linkvalue, stat.Mode())
	}
	if stat.Mode()&fs.ModeSymlink == 0 {
		return resp, errors.Errorf("%s file %s is not a symlink, (%s)", linkname, linkvalue, stat.Mode())
	}

	// Unlinked stat
	linkTarget, err := os.Readlink(linkvalue)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to get %s unlinked stat", linkname)
	}
	if !strings.HasPrefix(linkTarget, VersionDir) {
		return resp, errors.Errorf(
			"%s (%s) symlink targets on a wrong file %s; %s prefix should present",
			linkvalue, linkname, linkTarget, VersionDir,
		)
	}
	prefix, base := path.Split(linkTarget)
	if base != expectedBase {
		return resp, errors.Errorf(
			"%s (%s) symlink targets on a wrong file %s; unexpected base: '%s' ('%s' should be)",
			linkvalue, linkname, linkTarget, base, expectedBase,
		)
	}

	_, baseWithVersion := path.Split(prefix[:len(prefix)-1])
	err = resp.SetVersion(baseWithVersion)
	if err != nil {
		return resp, errors.Wrapf(err,
			"%s (%s) symlink targets on a wrong file %s; invalid goVersion path part %s;",
			linkvalue, linkname, linkTarget, baseWithVersion,
		)

	}

	log.Debugf("unlinked name %s; extracted go version %s", linkTarget, resp.version)
	return resp, nil
}

func checkPathEnv() (gobinInPath bool) {
	PATH := os.Getenv("PATH")
	for _, path := range strings.Split(PATH, ":") {
		if path == GoBinaryDir {
			return true
		}
	}
	return false
}

func checkVersionOfBin(binpath string) (ver string, err error) {
	cmd := exec.Command(binpath, "version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ver, errors.Wrapf(err, "failed to get stdout Cmd (%s) pipe", cmd)
	}
	err = cmd.Start()
	if err != nil {
		return ver, errors.Wrapf(err, "failed to start Cmd (%s)", cmd)
	}
	resultBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return ver, errors.Wrapf(err, "failed to read cmd stdout %s", cmd)
	}
	cmd.Wait()
	parts := bytes.Split(resultBytes, []byte(" "))
	return string(parts[2]), nil
}

func checkConsistency() goData {
	log.Debug("Start consistency check...\n")

	gobinGo, err := validateGolink("Gobin", GobinSymlink, gobinBase)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debugf("[OK] goBin (%s) is a symlink pointing on possible repo with golang bins", GobinSymlink)

	//---
	gopathGo, err := validateGolink("Gopath", GopathSymlink, gopathBase)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Debugf("[OK] goPath (%s) is a symlink pointing on possible goPath repo", GopathSymlink)

	//---
	if gobinGo.version != gopathGo.version {
		log.Fatalf(
			"broken gobin (%s) - gopath (%s) symlink pair; they target on different go versions %s - %s",
			GobinSymlink, GopathSymlink, gobinGo.version, gopathGo.version,
		)
	}

	//---
	gobinInPath := checkPathEnv()
	if !gobinInPath {
		log.Fatalf("[ERROR] gobin (%s) not found in PATH enviroment variable", GobinSymlink)
	}
	goBinVer, err := checkVersionOfBin(GoBinaryPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debugf("go bin version %s", goBinVer)

	gobinGo.isCurrent = true
	return gobinGo
}
