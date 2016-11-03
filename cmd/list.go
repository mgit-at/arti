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
}

func listCheckFlagsAndArgs(cmd *cobra.Command, args []string) string {
	if len(args) < 1 {
		cmd.Help()
		os.Exit(1)
	}

	return args[0]
}

func listRun(cmd *cobra.Command, args []string) {
	snp := listCheckFlagsAndArgs(cmd, args)

	s := selectStore(snp)

	artifacts, err := s.List()
	if err != nil {
		log.Fatalln("listing artifacts failed:", err)
	}

	for name, versions := range artifacts {
		log.Printf("%s:", name)
		for _, v := range versions {
			if numericSize {
				log.Printf("  %v: %12d  %s", v.Version, v.Filesize, v.Filename)
			} else {
				size, mult := humanizeBytes(v.Filesize)
				log.Printf("  %v: %6.1f %sB  %s", v.Version, size, mult, v.Filename)
			}
		}
	}
}
