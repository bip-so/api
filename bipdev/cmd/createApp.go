/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func writeFile(FolderName string, FileName string, FileExtention string) {
	FileNameWithGo := FileName + FileExtention
	go_file_content := "package " + FolderName + "\n"
	go_file_data := []byte(go_file_content)
	err := ioutil.WriteFile(FolderName+"/"+FileNameWithGo, go_file_data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createFolder(FolderName string) error {
	_, err := os.Stat(FolderName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(FolderName, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	return nil
}

// createAppCmd represents the createApp command
var createAppCmd = &cobra.Command{
	Use:   "createApp",
	Short: "Will scaffold sample app for BIP Platfrom",
	Run: func(cmd *cobra.Command, args []string) {
		// check is folder name is present
		if len(args) == 0 {
			fmt.Println("Error: Please provide the app name")
			os.Exit(1)
		}
		folderName := strings.ToLower(args[0])

		folderExistFlg, _ := folderExists(folderName)
		if folderExistFlg == true {
			fmt.Println("Error: Folder exists")
			os.Exit(1)
		}

		tryCreateFolder := createFolder(folderName)
		if tryCreateFolder != nil {

			fmt.Println(tryCreateFolder)
			os.Exit(1)
		}

		filer_list := []string{
			"controller",
			"dto",
			"errors",
			"route_handler",
			"routes",
			"serializers",
			"services",
			"validators",
		}

		for _, fileName := range filer_list {
			fmt.Printf("Writing %v in %v \n", fileName, folderName)
			// get current working directory
			// cwd, _ := os.Getwd()
			writeFile(folderName, fileName, ".go")
		}

		// Module with the same file name.
		//writeFile(folderName, folderName, ".md")
		fmt.Println("Module is created.")
	},
}

func init() {
	rootCmd.AddCommand(createAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createAppCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
