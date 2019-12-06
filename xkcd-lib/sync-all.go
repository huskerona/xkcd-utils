package xkcd_lib

import (
	"fmt"
	"log"
	"os"
)

// Downloads all the comics from XKCD and writes them into offline index storage.
// Any file that already exists locally is skipped.
func SyncAll() error {
	defer trace("SyncAll")()
    xkcd := XKCD{}

	err := xkcd.GetLatestComic()

	fmt.Println("3. &xkcd: ", &xkcd.Number)


	if err != nil {
		return fmt.Errorf("xkcd: sync all: %v\n", err)
	}

	fmt.Println("Downloading... (this may take some time)")

	doDownload(xkcd.Number)

	fmt.Printf("\nXKCD synced with offline storage")

	return err
}

// Indexes the files to the offline index storage
func IndexFiles(xkcd *XKCD, ls *LocalStorage) error {
	defer trace("IndexFiles")()

	var err error
	var imageName string

	if !ls.JSONFileExists {
		if err = WriteJSON(xkcd); err != nil {
			ShowVerbose(fmt.Sprintf("%v\n", err))
		}
	}

	if !ls.ImageFileExists {
		if imageName, err = GetImageFileName(xkcd.ImageURL); err != nil {
			ShowVerbose(fmt.Sprintf("%v\n", err))
			return err
		}

		path := fmt.Sprintf("%s/%d", OfflineIndexPath, xkcd.Number)
		imageOutput := fmt.Sprintf("%s/%s", path, imageName)

		var imageContent []byte

		if imageContent, err = xkcd.DownloadImage(); err != nil {
			ShowVerbose(fmt.Sprintf("%v\n", err))
			return err
		}

		if _, err = WriteImage(imageOutput, imageContent); err != nil {
			ShowVerbose(fmt.Sprintf("%v\n", err))
			return err
		}
	}

	return nil
}

// Retrieves the status of the given comic by enumeration through the files in the directory.
// The only known name is the JSON file name. The other files are assumed to be the actual comic image.
func GetIndexStatus(comicNum int) *LocalStorage {
	defer trace("GetIndexStatus")()

	ls := LocalStorage{JSONFileExists: false, ImageFileExists: false}

	path := fmt.Sprintf("%s/%d", OfflineIndexPath, comicNum)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ls
	}

	f, err := os.Open(path)

	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)

	for _, file := range files {
		if file.Name() == JSONFileName {
			ls.JSONFileExists = true
		} else {
			ls.ImageFileExists = true
		}
	}

	return &ls
}

func doDownload(latestComicNum int) {
	defer trace("doDownload")()

	for i := 1; i <= latestComicNum; i++ {
		ls := GetIndexStatus(i)

		if ls.ImageFileExists && ls.JSONFileExists {
			continue
		}

		fmt.Printf("\nDownloading #%d of %d...", i, latestComicNum)

		xkcd := XKCD{}
		var err error

		if err = xkcd.GetComic(i); err != nil {
			ShowVerbose(fmt.Sprintf("Error getting commic number %d. Skipping!", i))
			continue
		}

		if !ls.JSONFileExists || !ls.ImageFileExists {
			if err := CreateComicFolder(i); err != nil {
				log.Fatal(err)
			}
		}

		if err := IndexFiles(&xkcd, ls); err != nil {
			fmt.Printf("error: %v", err)
			continue
		}

		fmt.Print("OK")
	}
}