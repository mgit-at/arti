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

package cmd

import (
	"log"
	"math"
	"strings"

	"github.com/mgit-at/arti/store"

	"github.com/spf13/viper"
)

const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

var (
	artifactName    string
	artifactVersion string
)

func selectStore(nameAndPath string) store.Store {
	stores := viper.Sub("stores")
	if stores == nil {
		log.Fatal("no stores specified!")
	}

	np := strings.SplitN(nameAndPath, "/", 2)
	if len(np) < 2 {
		np = append(np, "")
	}
	cfg := stores.Sub(np[0])
	if cfg == nil {
		log.Fatalln("no such store:", np[0])
	}
	s, err := store.NewStore(cfg, np[1])
	if err != nil {
		log.Fatalln("unable to initialize store:", err)
	}
	//	log.Printf("using store(type: %T): %+v", s, s)
	return s
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanizeBytes(size int64) (float64, string) {
	e := math.Floor(logn(float64(size), 1024))
	switch e {
	case 0:
		return float64(size), ""
	case 1:
		return float64(size) / KB, "k"
	case 2:
		return float64(size) / MB, "M"
	case 3:
		return float64(size) / GB, "G"
	case 4:
		return float64(size) / TB, "T"
	default:
		return float64(size) / PB, "P"
	}
}
