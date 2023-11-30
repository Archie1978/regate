package version

import (
	"runtime/debug"
	"time"
)

var GitTag = ""

const (
	govcsTimeLayout = "2006-01-02T15:04:05Z"
	ourTimeLayout   = "20060102"
)

// buildInfoVCS returns VCS information of the build.
func buildInfoVCS() (s VCSInfo, ok bool) {
	info, ok := debug.ReadBuildInfo()
	if ok {
		s.Deps = info.Deps
		for _, v := range info.Settings {
			switch v.Key {
			case "vcs.revision":
				s.Commit = v.Value
			case "vcs.modified":
				if v.Value == "true" {
					s.Dirty = true
				}
			case "vcs.time":
				t, err := time.Parse(govcsTimeLayout, v.Value)
				if err == nil {
					s.Date = t
				}
			}
		}
		if s.Commit != "" && s.Date.Unix() != 0 {
			ok = true
		}
	}
	return
}

// VCSInfo represents the git repository state.
type VCSInfo struct {
	Deps   []*debug.Module
	Commit string // head commit hash
	Date   time.Time
	Dirty  bool
}

func Version() (version string) {
	if GitTag != "" {
		version = GitTag
	}

	vcsinfo, ok := buildInfoVCS()
	if !ok {
		return "0.0.0-unknown"
	}

	if version == "" {
		return "0.0.0-unknown-" + vcsinfo.Commit
	}
	if vcsinfo.Dirty {
		return version + "-unstable-" + vcsinfo.Commit
	}
	return version + "-stable"
}

func Date() (binaryDate time.Time) {
	vcsinfo, ok := buildInfoVCS()
	if !ok {
		return time.Unix(0, 0)
	}
	return vcsinfo.Date
}

func PackageDepance() (version string) {
	if GitTag != "" {
		version = GitTag
	}

	vcsinfo, ok := buildInfoVCS()
	if !ok {
		return "0.0.0-unknown"
	}

	if version == "" {
		return "0.0.0-unknown-" + vcsinfo.Commit
	}
	if vcsinfo.Dirty {
		return version + "-unstable-" + vcsinfo.Commit
	}
	return version + "-stable"
}
