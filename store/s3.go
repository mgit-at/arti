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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

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
		s.version = 4
		fallthrough
	case 4:
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
		if strings.HasSuffix(obj.Key, CSumExt) {
			continue
		}
		nv, f := path.Split(obj.Key)
		n, v := path.Split(strings.TrimSuffix(nv, "/"))
		n = strings.TrimSuffix(n, "/")
		if a, err := MakeArtifactListEntry(v, f, obj.Size); err == nil {
			list[n] = append(list[n], a)
		} else {
			log.Printf("ignoring invalid object: '%s'", obj.Key) // TODO: debug output
		}
	}

	return
}

func (s *S3Store) Has(artifact Artifact) (exists bool, filename string, err error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	hasHash := false
	p := path.Join(artifact.Name, artifact.Version.String())
	objCh := s.client.ListObjectsV2(s.bucket, p, true, doneCh)
	for obj := range objCh {
		if obj.Err != nil {
			err = fmt.Errorf("Error while listing objects: %v", obj.Err)
			return
		}
		f := path.Base(obj.Key)
		if strings.HasSuffix(f, CSumExt) {
			hasHash = true
		} else {
			if filename != "" {
				err = fmt.Errorf("found more then one file")
				return
			}
			filename = f
		}
	}
	if filename != "" || hasHash {
		exists = true
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

	if exists, _, err := s.Has(artifact); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("artifact already exists")
	}

	basename := filepath.Base(filename)
	p := path.Join(artifact.Name, artifact.Version.String(), basename)
	n, err := s.client.FPutObject(s.bucket, p, filename, "application/octet-stream")
	if err != nil {
		return fmt.Errorf("Error uploading file '%s': %v", basename, err)
	}

	hashSum, err := calcCSum(filename)
	if err != nil {
		return fmt.Errorf("Error calculating checksum of '%s': %v", basename, err)
	}
	_, err = s.client.PutObject(s.bucket, p+CSumExt, strings.NewReader(hashSum), "application/octet-stream")
	if err != nil {
		return fmt.Errorf("Error uploading hash: %v", err)
	}

	log.Printf("successfully uploaded %d Bytes to '%s:/%s'", n, s.bucket, p)
	log.Printf("%s", hashSum)

	return nil
}

func (s *S3Store) Get(artifact Artifact, filename string, keepCorrupted bool) (err error) {
	var exists bool
	var f string
	if exists, f, err = s.Has(artifact); err != nil {
		return
	}
	if !exists {
		err = fmt.Errorf("artifact not found")
		return
	}

	p := path.Join(artifact.Name, artifact.Version.String(), f)
	var hashObj *minio.Object
	if hashObj, err = s.client.GetObject(s.bucket, p+CSumExt); err != nil {
		return fmt.Errorf("Error while fetching hash: %v", err)
	}

	var hashValue bytes.Buffer
	if _, err = io.Copy(&hashValue, hashObj); err != nil {
		err = fmt.Errorf("Error while fetching hash: %v", err)
		return
	}

	target := filename
	if target == "" {
		target = f
	}
	if err = s.client.FGetObject(s.bucket, p, target); err != nil {
		err = fmt.Errorf("Error fetching file '%s': %v", f, err)
		return
	}
	var valid bool
	if valid, err = checkCSum(f, hashValue.String()); err != nil {
		return
	}
	if !valid {
		if !keepCorrupted {
			os.Remove(f)
		}
		return fmt.Errorf("hash-sum mismatch!")
	}
	return
}

func (s *S3Store) Del(artifact Artifact) (err error) {
	objectsCh := make(chan string, 10)
	errorCh := s.client.RemoveObjects(s.bucket, objectsCh)

	doneCh := make(chan struct{})
	defer close(doneCh)

	p := path.Join(artifact.Name, artifact.Version.String())
	objCh := s.client.ListObjectsV2(s.bucket, p, true, doneCh)
	for obj := range objCh {
		if obj.Err != nil {
			err = fmt.Errorf("Error while listing objects: %v", obj.Err)
			return
		}
		objectsCh <- obj.Key
	}
	close(objectsCh)

	errCnt := 0
	for e := range errorCh {
		err = e.Err
		errCnt++
	}
	if errCnt > 0 {
		err = fmt.Errorf("%d Errors during deletion, last: %v", errCnt, err)
	}
	return
}
