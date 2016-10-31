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
	"strings"

	"github.com/mgit-at/arti/store"

	"github.com/spf13/viper"
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
	s, err := store.NewStore(stores.Sub(np[0]), np[1])
	if err != nil {
		log.Fatalln("unable to initialize store:", err)
	}
	log.Printf("using store(type: %T): %+v", s, s) // TODO: debug output
	return s
}
