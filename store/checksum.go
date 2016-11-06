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
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	CSumExt           = ".checksum"
	CSumAlgoSeperator = ":"

	CSumAlgoSHA256 = "sha256"

	CSumAlgoDefault = CSumAlgoSHA256
)

type CSumAlgo struct {
	calc  func(filename string) (string, error)
	check func(filename, toCompare string) (bool, error)
}

var CSumAlgos map[string]CSumAlgo

func init() {
	CSumAlgos = make(map[string]CSumAlgo)

	CSumAlgos[CSumAlgoSHA256] = CSumAlgo{calcSHA256, checkSHA256}
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
	return hex.EncodeToString(hash.Sum(hashSum)), nil
}

func checkSHA256(filename, toCompare string) (bool, error) {
	hashSum, err := calcSHA256(filename)
	return hashSum == toCompare, err
}

func calcCSum(filename string) (string, error) {
	cs, err := CSumAlgos[CSumAlgoDefault].calc(filename)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{CSumAlgoDefault, cs}, CSumAlgoSeperator), nil
}

func checkCSum(filename, toCompare string) (bool, error) {
	tmp := strings.SplitN(toCompare, CSumAlgoSeperator, 2)
	if len(tmp) != 2 {
		return false, fmt.Errorf("invalid checksum file")
	}
	algo, found := CSumAlgos[tmp[0]]
	if !found {
		return false, fmt.Errorf("unkown checksum algorithm: %s", tmp[0])
	}
	return algo.check(filename, tmp[1])
}
