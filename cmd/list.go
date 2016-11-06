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
	"sort"

	"github.com/blang/semver"
	"github.com/mgit-at/arti/store"
	"github.com/spf13/cobra"
)

type ByteSize int64

var listCmd = &cobra.Command{
	Use:     "list <store>/<bucket>",
	Aliases: []string{"ls"},
	Short:   "list all artifacts in the store",
	Run:     listRun,
}

var (
	numericSize bool
)

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&numericSize, "numeric-size", "N", false, "print file sizes in bytes rather than in human readable format")
	listCmd.Flags().StringVarP(&artifactName, "name", "n", "", "list all version of this artifact")
	listCmd.Flags().StringVarP(&artifactVersion, "version", "v", "", "range of version to list")
}

func listCheckFlagsAndArgs(cmd *cobra.Command, args []string) (snp string, versions semver.Range) {
	if len(args) < 1 {
		cmd.Help()
		os.Exit(1)
	}
	snp = args[0]

	versions = nil
	if artifactVersion != "" {
		var err error
		if versions, err = semver.ParseRange(artifactVersion); err != nil {
			log.Fatalln("invalid version range:", err)
		}
	}

	return
}

func listRun(cmd *cobra.Command, args []string) {
	snp, versions := listCheckFlagsAndArgs(cmd, args)

	s := selectStore(snp)

	artifacts, err := s.List(artifactName, versions)
	if err != nil {
		log.Fatalln("listing artifacts failed:", err)
	}

	if artifactName == "" {
		listNames(artifacts)
	} else {
		if av, found := artifacts[artifactName]; found {
			listVersions(av)
		}
	}
}

func listNames(a store.ArtifactList) {
	names := []string{}
	for name := range a {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		sizeTotal := int64(0)
		for _, v := range a[name] {
			sizeTotal += v.Filesize
		}

		if numericSize {
			log.Printf("%s\t%d %12d", name, len(a[name]), sizeTotal)
		} else {
			size, mult := humanizeBytes(sizeTotal)
			log.Printf("%s\t%d %6.1f%sB", name, len(a[name]), size, mult)
		}
	}
}

func listVersions(av store.ArtifactVersions) {
	sort.Sort(sort.Reverse(av))
	for _, v := range av {
		if numericSize {
			log.Printf("%v\t%12d %s", v.Version, v.Filesize, v.Filename)
		} else {
			size, mult := humanizeBytes(v.Filesize)
			log.Printf("%v\t%6.1f%sB %s", v.Version, size, mult, v.Filename)
		}
	}
}
