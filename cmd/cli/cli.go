package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/spf13/cobra"
)

const (
	bucket  = "saferwall-samples"
	region  = "us-east-1"
	timeout = 0
)

var (
	forceRescan bool
	outputDir   string
	username    string
	password    string
)

// scanFile scans an individual file or a directory.
func scanFile(cmd *cobra.Command, args []string) {

	pathToScan := args[0]
	_, err := os.Stat(pathToScan)
	if os.IsNotExist(err) {
		log.Fatalf("%s does not exist", pathToScan)
	}

	// Walk over directory.
	fileList := []string{}
	filepath.Walk(pathToScan, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Obtain a token.
	token, err := login(username, password)
	check(err)

	// Upload files
	for _, filename := range fileList {

		// Get sha256
		data, err := ioutil.ReadFile(filename)
		check(err)
		sha256 := getSha256(data)

		// Do we have the file in S3.
		if !isFileFoundInObjStorage(sha256) {
			uploadObject(bucket, region, sha256, filename)
		}

		// Check if we the file exists in the DB.
		if !isFileFoundInDB(sha256, token) {
			scan(sha256, token)
			time.Sleep(timeout * time.Second)
			continue
		}

		// Force rescan the file?.
		if forceRescan {
			rescan(sha256, token)
			time.Sleep(timeout * time.Second)
		}
	}
}

func s3upload(cmd *cobra.Command, args []string) {

	filePath := args[0]
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		exitErrorf("%s does not exist", filePath)
	}

	objKeys := listObject(bucket, region, false)

	// Walk over directory.
	fileList := []string{}
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Upload files
	for _, filename := range fileList {
		// Check if we have the file already in our database.
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("failed to read file %s", filename)
			continue
		}
		key := crypto.GetSha256(dat)
		found := stringInSlice(key, objKeys)
		if !found {
			uploadObject(bucket, region, key, filename)
		} else {
			fmt.Printf("file name %s already in s3 bucket", filename)
		}
	}

}

// rescanFile reads a list of sha256 from the clipboard and trigger a rescan.
func rescanFile(cmd *cobra.Command, args []string) {

	// Obtain a token.
	token, err := login(username, password)
	check(err)

	clipContent, err := clipboard.ReadAll()
	check(err)

	shaList := strings.Split(clipContent, "\r\n")
	for _, sha256 := range shaList {
		rescan(sha256, token)

		// Wait for file to be scanned.
		time.Sleep(timeout * time.Second)
	}

}

// processAuthTokens processes a given username and password as env variables.
func processAuthTokens() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv(DefaultAuthUsername)
	password := os.Getenv(DefaultAuthPassword)
	err = SetAuthentificationData(username, password)
	if err != nil {
		check(err)
	}
}

// downloadFile a list of sha256 from the clipboard and download them.
func downloadFile(cmd *cobra.Command, args []string) {

	clipContent, err := clipboard.ReadAll()
	check(err)

	shaList := strings.Split(clipContent, "\r\n")
	for _, sha256 := range shaList {
		// Create a new file.
		filePath := path.Join(outputDir, sha256)
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			continue
		}
		err = downloadObject(bucket, region, sha256, f)
		if err != nil {
			log.Println(err)
			continue
		}

	}
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "sfwcli",
		Short: "A cli tool for saferwall.com",
		Long:  "A cli tool to interfact with saferwall APIs (scan, rescan, upload, ...)",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Vesion number",
		Long:  "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("You are using version 0.0.1")
		},
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan file",
		Long:  "Scan a file or directory",
		Args:  cobra.MinimumNArgs(1),
		Run:   scanFile,
	}

	var rescanCmd = &cobra.Command{
		Use:   "rescan",
		Short: "Resccan file",
		Long:  "Rescan a file or directory",
		Run:   rescanFile,
	}

	var s3UploadCmd = &cobra.Command{
		Use:   "s3upload",
		Short: "S3 upload",
		Long:  "Batch upload to S3",
		Args:  cobra.MinimumNArgs(1),
		Run:   s3upload,
	}

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download file",
		Long:  "Download a file or directory",
		Run:   downloadFile,
	}

	// Init root command.
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(rescanCmd)
	rootCmd.AddCommand(s3UploadCmd)
	rootCmd.AddCommand(downloadCmd)

	// Init flags
	scanCmd.Flags().BoolVarP(&forceRescan, "forcerescan", "f", false, "Force rescan the file.")
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory to download the files (required")
	downloadCmd.MarkFlagRequired("output")

	// load config
	processAuthTokens()
	// Get credentials.
	username = os.Getenv(DefaultAuthUsername)
	password = os.Getenv(DefaultAuthPassword)
	if username == "" || password == "" {
		log.Fatal("SAFERWALL_AUTH_USERNAME or SAFERWALL_AUTH_PASSWORD env variable are not set.")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
