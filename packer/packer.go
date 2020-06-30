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
package packer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type TPackerParam struct {
	StartCommitID string
	EndCommitID   string
	Format        string
	OutputPath    string
	ArchivePrefix string
	DiffFilter    string
	Verbose       bool
	NoPrefix      bool
}

type TChangeItem struct {
	ChangeType string `json:"type" bson:"type"`
	OldPath    string `json:"old_path" bson:"old_path"`
	NewPath    string `json:"new_path" bson:"new_path"`
}

type TOtaInfo struct {
	ProjectName    string        `json:"project_name" bson:"project_name"`
	LastOtaVersion string        `json:"last_ota_version" bson:"ota_version"`
	OtaVersion     string        `json:"ota_version" bson:"ota_version"`
	FullUpdate     bool          `json:"is_full_update" bson:"is_full_update"`
	Changes        []TChangeItem `json:"changes" bson:"changes"`
}

func execCommand(cmd_name string, params []string) (string, error) {
	cmd := exec.Command(cmd_name, params...)

	//fmt.Println(cmd.Path)
	//fmt.Println(cmd.Dir)
	//fmt.Println(cmd.Args)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	return out.String(), err
}

func execBashCommand(cmd_str string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmd_str)

	//fmt.Println(cmd.Path)
	//fmt.Println(cmd.Dir)
	//fmt.Println(cmd.Args)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	return out.String(), err
}

func handleDiffString(diff string, otaInfo *TOtaInfo) error {
	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		itemInfo := strings.Fields(scanner.Text())
		if len(itemInfo) == 2 {
			change_item := TChangeItem{
				ChangeType: itemInfo[0],
				OldPath:    itemInfo[1],
				NewPath:    itemInfo[1],
			}
			otaInfo.Changes = append(otaInfo.Changes, change_item)
		}
		if len(itemInfo) == 3 {
			if strings.HasPrefix(itemInfo[0], "R") {
				change_item := TChangeItem{
					ChangeType: itemInfo[0],
					OldPath:    itemInfo[1],
					NewPath:    itemInfo[2],
				}
				otaInfo.Changes = append(otaInfo.Changes, change_item)

			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func getDateTimeString() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func Pack(params *TPackerParam) error {
	workDirectory, _ := execCommand("pwd", []string{"-L"})
	rootDirectory, _ := execCommand("git", []string{"rev-parse", "--show-toplevel"})

	if workDirectory != rootDirectory {
		fmt.Println("Error: You must operate at the root of the project!!!")
		return errors.New("You must operate at the root of the project!!!")
	}

	lastRevision, _ := execCommand("git", []string{"rev-parse",
		"--short", params.StartCommitID})
	currentRevision, _ := execCommand("git", []string{"rev-parse",
		"--short", params.EndCommitID})

	// Remove '\n'
	rootDirectory = strings.TrimSuffix(rootDirectory, "\n") + "/"
	lastRevision = strings.TrimSuffix(lastRevision, "\n")
	currentRevision = strings.TrimSuffix(currentRevision, "\n")

	if params.OutputPath == "" {
		params.OutputPath = "./"
	} else {
		params.OutputPath += "/"
	}

	ret, _ := execCommand("git", []string{"diff", "-r", "--name-status",
		"--diff-filter=ACMRD", lastRevision, currentRevision})
	// Check if it is full updates
	fullUpdate := false
	if strings.Contains(lastRevision, currentRevision) {
		fullUpdate = true
	} else {
		fullUpdate = false
	}

	// Generate ota_info.json by diff files
	otaInfo := TOtaInfo{
		ProjectName:    params.ArchivePrefix,
		LastOtaVersion: lastRevision,
		OtaVersion:     currentRevision,
		FullUpdate:     fullUpdate,
		Changes:        nil,
	}

	err := handleDiffString(ret, &otaInfo)
	if err != nil {
		fmt.Println("Pack error: ", err)
		return err
	}

	data, err := json.MarshalIndent(&otaInfo, "", "    ")
	if err != nil {
		fmt.Println("Pack error: json.marshal failed with ", err)
		return err
	}

	// Set archive name
	incrementalUpdatePackageName := params.ArchivePrefix +
		"-" + getDateTimeString() +
		"-" + lastRevision +
		"-to-" + currentRevision +
		"." + params.Format
	fmt.Printf("-----------------------------------------------------------------\n")
	fmt.Printf("%16s  |  %s\n", "Project Name", params.ArchivePrefix)
	fmt.Printf("------------------+----------------------------------------------\n")
	fmt.Printf("%16s  |  %s\n", "Output Path", params.OutputPath)
	fmt.Printf("------------------+----------------------------------------------\n")
	fmt.Printf("%16s  |  %s\n", "Output", incrementalUpdatePackageName)
	fmt.Printf("------------------+----------------------------------------------\n")
	fmt.Printf("%16s  |  %s\n", "Changelog", "ota_info.json")
	fmt.Printf("------------------+----------------------------------------------\n")
	fmt.Printf("%s\n", string(data))
	fmt.Printf("------------------+----------------------------------------------\n")
	// While the result of `git diff --diff-filter=ACMR` is null
	ret, _ = execCommand("git", []string{"diff", "-r", "--name-status",
		"--diff-filter=ACMR", lastRevision, currentRevision})
	if ret == "" {
		// Write OTA changelog to file with json format
		_ = ioutil.WriteFile("ota_info.json", data, 0644)
		//fmt.Printf("%s\n", string(data))
		// Package
		_, err = execBashCommand("tar -cvf " + params.OutputPath +
			incrementalUpdatePackageName + " ota_info.json && " +
			"rm ota_info.json")
		if err != nil {
			fmt.Println("Pack error: failed to execute bash command ", err)
			return err
		}
		return nil
	} else {
		// Write OTA changelog to file with json format
		_ = ioutil.WriteFile(params.OutputPath+"ota_info.json", data, 0644)
		//fmt.Printf("%s\n", string(data))
	}

	// Export a specific commit with git-archive
	if fullUpdate {
		fmt.Println("Export a specific commit : " + currentRevision)
		_, err = execBashCommand("git archive --format " + params.Format +
			" --output " + params.OutputPath + incrementalUpdatePackageName +
			" " + currentRevision)
	} else {
		// Generage incremental update package
		_, err = execBashCommand("git diff -r --no-commit-id --name-only --diff-filter=" +
			params.DiffFilter +
			" --line-prefix=" +
			rootDirectory + " " +
			lastRevision + " " +
			currentRevision +
			" | xargs git archive --format " + params.Format +
			" -o " + params.OutputPath +
			incrementalUpdatePackageName +
			" " + currentRevision)
		if err != nil {
			fmt.Println("Pack error: failed to execute bash command ", err)
			return err
		}
	}

	// Package ota_info.json into upgrade package
	mvTargetDirCmd := ""
	if params.OutputPath != "./" {
		mvTargetDirCmd = " && mv " + incrementalUpdatePackageName +
			" " + params.OutputPath
	}
	_, err = execBashCommand("mkdir -p ~/.tmp/ota_packer && tar -xvf " +
		params.OutputPath + incrementalUpdatePackageName +
		" -C ~/.tmp/ota_packer && mv " +
		params.OutputPath + "ota_info.json ~/.tmp/ota_packer && " +
		"rm " + params.OutputPath + incrementalUpdatePackageName + " && " +
		"cd ~/.tmp/ota_packer && tar -cvf " +
		incrementalUpdatePackageName + " * && " +
		"mv ~/.tmp/ota_packer/" + incrementalUpdatePackageName +
		" " + rootDirectory + " && " +
		"cd " + rootDirectory + " && " +
		"rm -rf ~/.tmp/ota_packer " +
		mvTargetDirCmd)
	if err != nil {
		fmt.Println("Pack error: failed to execute bash command ", err)
		return err
	}

	return nil
}
