package xkcd_lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var OfflineIndexPath = "/.xkcd"
const JSONFileName = "info.0.json"

// Sets the location of the offline indexed storage to user's home folder with .xkcd as subfolder
func SetOfflineIndexPath(homeDir string) {
	OfflineIndexPath = fmt.Sprintf("%s%s", homeDir, OfflineIndexPath)
}

// Writes a info.0.json file into the indexed folder for the given comic
func WriteJSON(xkcd *XKCD) error {
	var result []byte
	var err error
	ShowVerbose("Writing JSON file...\n")

	if xkcd == nil {
		return fmt.Errorf("xkcd is nil")
	}

	path := fmt.Sprintf("%s/%d", OfflineIndexPath, xkcd.Number)
	jsonOutput := fmt.Sprintf("%s/%s", path, JSONFileName)


	if result, err = json.MarshalIndent(xkcd, "", "  "); err != nil {
		return fmt.Errorf("writing JSON output: %v", err)
	}

	ShowVerbose(fmt.Sprintf("Writing JSON to: %s\n", path))

	if err = ioutil.WriteFile(jsonOutput, result, 0755); err != nil {
		return err
	}

	return nil
}

// Writes an image to a destination.
// The destination is made of OfflineIndexPath + Comic Number + original image file name
func WriteImage(filename string, imageByte []byte) (int, error) {
	if len(imageByte) == 0 {
		return 0, fmt.Errorf("empty image")
	}

	if err := ioutil.WriteFile(filename, imageByte, os.ModePerm); err != nil {
		return 0, err
	}

	return len(imageByte), nil
}

// Check if the Offline index directory exists together with a subdirectory that is the comic ID.
func IsOfflineIndexAvailable(path string) error {

	ShowVerbose(fmt.Sprintf("Checking folder '%s'...", path))

	if _, err := os.Stat(path); os.IsExist(err) {
		ShowVerbose("FOUND\n")
		return nil
	}

	ShowVerbose("Not found. Trying to create...")
	if err := os.MkdirAll(path, 0777); err != nil {
		ShowVerbose(fmt.Sprintf("FAILED - %v\n", err))
		return err
	}

	ShowVerbose("CREATED\n")
	return nil
}

func CreateComicFolder(comicNum int) error {
	destination := fmt.Sprintf("%s/%d", OfflineIndexPath, comicNum)

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		if err := os.MkdirAll(destination, 0777); err != nil {
			ShowVerbose(fmt.Sprintf("FAILED - %v\n", err))
			return err
		}
	}

	return nil
}

// Extracts an image file name from ImageURL field
func GetImageFileName(imageURL string) (string, error) {
	lastIndex := strings.LastIndex(imageURL, "/")

	if lastIndex == -1 {
		return "", fmt.Errorf("missing image name: %s", imageURL)
	}

	result := imageURL[lastIndex+1:]

	return result, nil
}

// Reads the JSON file and returns XKCD object
func readJSON(jsonFile string) (*XKCD, error) {
	file, err := ioutil.ReadFile(jsonFile)

	if err != nil {
		return nil, fmt.Errorf("error reading JSON '%s': %v", jsonFile, err)
	}

	var xkcd XKCD
	err = json.Unmarshal(file, &xkcd)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON '%s': %v", jsonFile, err)
	}

	return &xkcd, nil
}