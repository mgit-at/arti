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

	"github.com/mgit-at/arti/store"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload files to the store",
	Long:  `...tba...`,
	Run: func(cmd *cobra.Command, args []string) {

		stores := viper.Sub("stores")
		if stores == nil {
			log.Fatal("no stores specified!")
		}

		storeNames := []string{}
		for s := range stores.AllSettings() {
			storeNames = append(storeNames, s)
		}
		log.Printf("stores: %v\n", storeNames)

		s, err := store.NewStore(stores.Sub("minio"))
		if err != nil {
			log.Fatalln("unable to initialize store:", err)
		}
		log.Printf("using store(type: %T): %v", s, s)
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
