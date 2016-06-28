// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"path"
	"trans/analysis"

	"github.com/spf13/cobra"
)

var getstring_dbname string
var getstring_srcpath string

// getstringCmd represents the getstring command
var getstringCmd = &cobra.Command{
	Use:   "getstring",
	Short: "Extract chinese characters",
	Long:  `Extract Chinese characters from a file or directory and save it to a text file`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if len(getstring_srcpath) == 0 {
			cmd.Help()
			return
		}
		analysis.GetInstance().GetString(path.Clean(getstring_dbname), path.Clean(getstring_srcpath))
	},
}

func init() {
	RootCmd.AddCommand(getstringCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getstringCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getstringCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getstringCmd.Flags().StringVarP(&getstring_dbname, "db", "d", "dictionary.txt", "Translation data dictionary")
	getstringCmd.Flags().StringVarP(&getstring_srcpath, "src", "s", "", "The extracted file or directory path")
}
