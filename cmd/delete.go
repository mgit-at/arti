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

	"github.com/mgit-at/arti/store"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <store>/<bucket>",
	Aliases: []string{"del"},
	Short:   "delete artifacts from the store",
	Long:    `...tba...`,
	Run:     deleteRun,
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&artifactName, "name", "n", "", "the name of the artifact")
	deleteCmd.MarkFlagRequired("name")
	deleteCmd.Flags().StringVarP(&artifactVersion, "version", "v", "", "the version of the artifact (must adhere to the semantic versioning scheme)")
	deleteCmd.MarkFlagRequired("version")
}

func deleteCheckFlagsAndArgs(cmd *cobra.Command, args []string) (string, store.Artifact) {
	if len(args) < 1 {
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

	a, err := store.MakeArtifact(artifactName, artifactVersion)
	if err != nil {
		log.Fatalln("invalid artifact specification:", err)
	}

	return args[0], a
}

func deleteRun(cmd *cobra.Command, args []string) {
	snp, a := deleteCheckFlagsAndArgs(cmd, args)

	s := selectStore(snp)

	if err := s.Del(a); err != nil {
		log.Fatalln("deletion failed:", err)
	}
}
