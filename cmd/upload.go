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
	"os"
	"strings"

	"github.com/mgit-at/arti/store"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <store>/<path or bucket> <file>",
	Short: "upload files to the store",
	Long:  `...tba...`,
	Run:   uploadRun,
}

var (
	artifactName    string
	artifactVersion string
)

func init() {
	RootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&artifactName, "name", "n", "", "the name of the artifact")
	uploadCmd.MarkFlagRequired("name")
	uploadCmd.Flags().StringVarP(&artifactVersion, "version", "v", "", "the version of the artifact (must adhere to the semantic versioning scheme)")
	uploadCmd.MarkFlagRequired("version")
}

func checkFlagsAndArgs(cmd *cobra.Command, args []string) (string, string, store.Artifact) {
	if len(args) < 2 {
		cmd.Help()
		os.Exit(1)
	}
	if artifactName == "" {
		log.Println("please specifiy the artefact name")
		cmd.Help()
		os.Exit(1)
	}
	if artifactVersion == "" {
		log.Println("please specifiy the artefact version")
		cmd.Help()
		os.Exit(1)
	}

	a := store.Artifact{Name: artifactName}
	v, err := semver.ParseTolerant(artifactVersion)
	if err != nil {
		log.Fatalln("invalid version:", err)
	}
	a.Version = v

	return args[0], args[1], a
}

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
	log.Printf("using store(type: %T): %+v", s, s)
	return s
}

func uploadRun(cmd *cobra.Command, args []string) {
	snp, fn, a := checkFlagsAndArgs(cmd, args)

	s := selectStore(snp)

	if err := s.Put(a, fn); err != nil {
		log.Fatalln("upload failed:", err)
	}
}
