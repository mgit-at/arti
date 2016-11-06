// Copyright © 2016 mgIT GmbH <office@mgit.at>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/spf13/viper"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type Artifact struct {
	Name    string
	Version semver.Version
}

func MakeArtifact(name, version string) (Artifact, error) {
	a := Artifact{Name: name}
	v, err := semver.ParseTolerant(version)
	if err != nil {
		return Artifact{}, fmt.Errorf("invalid version string: %v", err)
	}
	a.Version = v
	return a, nil
}

type ArtifactVersion struct {
	Version  semver.Version
	Filename string
	Filesize int64
}

func MakeArtifactVersion(version, filename string, filesize int64) (ArtifactVersion, error) {
	a := ArtifactVersion{Filename: filename, Filesize: filesize}
	v, err := semver.ParseTolerant(version)
	if err != nil {
		return ArtifactVersion{}, fmt.Errorf("invalid version string: %v", err)
	}
	a.Version = v
	return a, nil
}

type ArtifactVersions []ArtifactVersion

func (a ArtifactVersions) Len() int {
	return len(a)
}

func (a ArtifactVersions) Less(i, j int) bool {
	return a[i].Version.LT(a[j].Version)
}

func (a ArtifactVersions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ArtifactList map[string]ArtifactVersions

type Store interface {
	List() (ArtifactList, error)
	Has(artifact Artifact) (bool, string, error)
	Put(artifact Artifact, filename string) error
	Get(artifact Artifact, filename string, keepCorrupted bool) error
	Del(artifact Artifact) error
}

func NewStore(cfg *viper.Viper, path string) (store Store, err error) {
	t := cfg.GetString("type")
	switch strings.ToLower(t) {
	case "s3":
		store, err = NewS3Store(cfg, path)
	default:
		err = fmt.Errorf("unknown store type: %s", t)
	}
	return
}
