/*
Copyright © 2020 Ethan <ycm_hy@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cli

import (
	_ "fmt"

	"github.com/yicm/OtaPackageTool/packer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genCmdInfo = &packer.TPackerParam{
	StartCommitID: "",
	EndCommitID:   "",
	Format:        "tar",
	OutputPath:    "",
	ArchivePrefix: "",
	DiffFilter:    "ACMRT",
	AppName:       "",
	AppRootPath:   "",
	Verbose:       false,
	NoPrefix:      false,
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate package file",
	Long:  `Generate a specific version package by entering different configuration parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(genCmdInfo.Verbose)
		//fmt.Println(genCmdInfo.StartCommitID)
		//fmt.Println(genCmdInfo.EndCommitID)
		//fmt.Println(genCmdInfo.Format)
		//fmt.Println(genCmdInfo.OutputPath)
		//fmt.Println(genCmdInfo.DiffFilter)
		//fmt.Println(genCmdInfo.ArchivePrefix)

		projectName := viper.GetString("project-name")
		appName := viper.GetString("app-name")
		appRootPath := viper.GetString("app-root-path")
		if projectName != "" {
			viper.Set("project-name", projectName)
			viper.WriteConfig()
			genCmdInfo.ArchivePrefix = projectName
		}

		if appName != "" {
			viper.Set("app-name", appName)
			viper.WriteConfig()
			genCmdInfo.AppName = appName
		}

		if appRootPath != "" {
			viper.Set("app-root-path", appRootPath)
			viper.WriteConfig()
			genCmdInfo.AppRootPath = appRootPath
		}

		// Start packing
		packer.Pack(genCmdInfo)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	genCmd.Flags().StringVarP(&genCmdInfo.StartCommitID, "start-commit-id", "s", "HEAD~1", "Start revision")
	genCmd.Flags().StringVarP(&genCmdInfo.EndCommitID, "end-commit-id", "e", "HEAD", "End revision")
	genCmd.Flags().StringVarP(&genCmdInfo.Format, "format", "f", "tar", "The format of the archive, supporting zip and tar")
	genCmd.Flags().StringVarP(&genCmdInfo.OutputPath, "output", "o", "", "Output destination path of the archive")
	genCmd.Flags().StringVarP(&genCmdInfo.ArchivePrefix, "prefix", "p", "ota_packer", "Prefixed to the filename in the archive while project name is not set.")
	genCmd.Flags().StringVarP(&genCmdInfo.DiffFilter, "diff-filter", "F", "ACMRT", "git diff --diff-filter and a similar designation")
	genCmd.Flags().BoolVarP(&genCmdInfo.Verbose, "verbose", "v", false, "Show packaging process statistics")
}
