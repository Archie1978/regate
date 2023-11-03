// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

// In go 1.18 and beyond, the go tool embeds VCS information into the build.

const (
	govcsTimeLayout = "2006-01-02T15:04:05Z"
	ourTimeLayout   = "20060102"
)

// buildInfoVCS returns VCS information of the build.
func buildInfoVCS(info *debug.BuildInfo) (s VCSInfo, ok bool) {
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
				s.Date = t.Format(ourTimeLayout)
			}
		}
	}
	if s.Commit != "" && s.Date != "" {
		ok = true
	}
	return
}

const (
	VersionMajor = 1          // Major version component of the current release
	VersionMinor = 13         // Minor version component of the current release
	VersionPatch = 5          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
)

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// ArchiveVersion holds the textual version string used for Geth archives. e.g.
// "1.8.11-dea1ce05" for stable releases, or "1.8.13-unstable-21c059b6" for unstable
// releases.
func ArchiveVersion(gitCommit string) string {
	vsn := Version
	if VersionMeta != "stable" {
		vsn += "-" + VersionMeta
	}
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

func VersionWithCommit(gitCommit, gitDate string) string {
	vsn := VersionWithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	if (VersionMeta != "stable") && (gitDate != "") {
		vsn += "-" + gitDate
	}
	return vsn
}
