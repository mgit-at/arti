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
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

const HASH_FILENAME = "sha256sum"

var (
	bucketNameRe = regexp.MustCompile("^[-_.A-Za-z0-9]{3,}$")
)

type S3Store struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	version         int
	location        string
	bucket          string

	client *minio.Client
}

func NewS3Store(cfg *viper.Viper, path string) (Store, error) {
	s := &S3Store{}

	if !bucketNameRe.MatchString(path) {
		return nil, fmt.Errorf("'%s' is not a valid bucket name", path)
	}
	s.bucket = path

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

func (s *S3Store) List() (list ArtifactList, err error) {
	list = make(ArtifactList)

	doneCh := make(chan struct{})
	defer close(doneCh)

	objCh := s.client.ListObjectsV2(s.bucket, "", true, doneCh)
	for obj := range objCh {
		if obj.Err != nil {
			err = fmt.Errorf("Error while listing objects: %v", obj.Err)
			return
		}
		nv, f := path.Split(obj.Key)
		if f == HASH_FILENAME {
			continue
		}
		n, v := path.Split(strings.TrimSuffix(nv, "/"))
		n = strings.TrimSuffix(n, "/")
		if a, err := MakeArtifactListEntry(v, f); err == nil {
			list[n] = append(list[n], a)
		} else {
			log.Printf("ignoring invalid object: '%s'", obj.Key) // TODO: debug output
		}
	}

	return
}

func (s *S3Store) MakeBucket() (err error) {
	err = s.client.MakeBucket(s.bucket, s.location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		if exists, err := s.client.BucketExists(s.bucket); err == nil && exists {
			return nil
		}
		return err
	}
	return
}

func (s *S3Store) Put(artifact Artifact, filename string) error {
	if err := s.MakeBucket(); err != nil {
		return err
	}

	basename := filepath.Base(filename)
	p := fmt.Sprintf("%s/%s/", artifact.Name, artifact.Version)
	n, err := s.client.FPutObject(s.bucket, p+basename, filename, "application/octet-stream")
	if err != nil {
		return fmt.Errorf("Error uploading file '%s': %v", basename, err)
	}
	hashSum, err := calcSHA256(filename)
	if err != nil {
		return fmt.Errorf("Error calculating SHA256 hash of '%s': %v", basename, err)
	}
	_, err = s.client.PutObject(s.bucket, p+HASH_FILENAME, strings.NewReader(hashSum), "application/octet-stream")
	if err != nil {
		return fmt.Errorf("Error uploading file '%s': %v", basename, err)
	}

	log.Printf("successfully uploaded %d Bytes to '%s:/%s'", n, s.bucket, p)
	log.Printf("SHA256: %s", hashSum)

	return nil
}

func (s *S3Store) Get(artifact Artifact) (err error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	p := fmt.Sprintf("%s/%s/", artifact.Name, artifact.Version)
	objCh := s.client.ListObjectsV2(s.bucket, p, true, doneCh)
	for obj := range objCh {
		if obj.Err != nil {
			err = fmt.Errorf("Error while listing objects: %v", obj.Err)
			return
		}
		f := path.Base(obj.Key)
		if f == HASH_FILENAME {
			continue
		}
		err = s.client.FGetObject(s.bucket, p+f, f)
		return // TODO: check Hash
	}
	err = fmt.Errorf("artifact not found")
	return
}
