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
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
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

func NewArtifact(name, version string) (*Artifact, error) {
	a := &Artifact{Name: name}
	v, err := semver.ParseTolerant(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version string: %v", err)
	}
	a.Version = v
	return a, nil
}

type Store interface {
	List() error
	Put(artifact Artifact, filename string) error
}

func calcSHA256(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	var hashSum []byte
	return base64.StdEncoding.EncodeToString(hash.Sum(hashSum)), nil
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
