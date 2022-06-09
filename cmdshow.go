package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type showList []goData

func (s showList) Len() int {
	return len(s)
}

func (s showList) Less(i, j int) bool {
	if s[i].major == s[j].major {
		if s[i].minor == s[j].minor {
			return s[i].patch < s[j].patch
		}
		return s[i].minor < s[j].minor
	}
	return s[i].major < s[j].major
}

func (s showList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

const (
	emptyPref   = "   "
	currentPref = " * "
)

func (s showList) String() string {
	if s == nil {
		return "no versions found"
	}
	var b strings.Builder
	b.WriteString(emptyPref)
	b.WriteString("Available go versions:\n\n")
	for _, v := range s {
		if v.isCurrent {
			b.WriteString(currentPref)
		} else {
			b.WriteString(emptyPref)
		}
		b.WriteString(v.version)
		b.WriteString("\n")
	}

	return b.String()
}

func commandShow() {
	currentGo := checkConsistency()

	goBin, err := os.Open(VersionDir)
	if err != nil {
		log.Fatalf("Failed to open goBin symlink %s", err)
	}
	binLs, err := goBin.ReadDir(0)
	if err != nil {
		log.Fatalf("Failed to read gobin dir %s", err)
	}
	showL := make(showList, 0, len(binLs))

	var garbage []string

	for _, item := range binLs {
		name := item.Name()
		goVer := goData{}
		err := goVer.SetVersion(name)
		if err != nil {
			garbage = append(garbage, name)
			continue
		}

		if name != currentGo.version {
			showL = append(showL, goVer)
			continue
		}

		showL = append(showL, currentGo)
	}

	if len(garbage) != 0 {
		log.Warnf("found some garbage in versionDir (%s): %s", VersionDir, strings.Join(garbage, ","))
	}

	sort.Sort(showL)
	fmt.Println(showL)
}
