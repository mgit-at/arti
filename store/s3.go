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
	"fmt"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

type S3Store struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	version         int
	location        string

	client *minio.Client
}

func NewS3Store(cfg *viper.Viper) (Store, error) {
	s := &S3Store{}
	s.endpoint = cfg.GetString("endpoint")
	s.accessKeyID = cfg.GetString("access-key-id")
	s.secretAccessKey = cfg.GetString("secret-access-key")
	s.useSSL = !cfg.GetBool("nossl")
	s.version = cfg.GetInt("version")
	s.location = cfg.GetString("location")

	var err error
	switch s.version {
	case 2:
		s.client, err = minio.NewV2(s.endpoint, s.accessKeyID, s.secretAccessKey, s.useSSL)
	case 0:
		fallthrough
	case 4:
		s.version = 4
		s.client, err = minio.NewV4(s.endpoint, s.accessKeyID, s.secretAccessKey, s.useSSL)
	default:
		err = fmt.Errorf("invalid S3 protocol version %d (only 2 and 4 are supported)", s.version)
	}
	if err != nil {
		return nil, err
	}

	return Store(s), nil
}

func (s *S3Store) List() error {
	return ErrNotImplemented
}

func (s *S3Store) Put(artifact Artifact, file string) error {
	return ErrNotImplemented
}
