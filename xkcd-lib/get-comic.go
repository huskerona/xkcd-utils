package xkcd_lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Retrieves the JSON of the latest XKCD comicl
func GetLatestComic() (*XKCD, error) {
	defer trace("GetLatestComic")()
	url := XKCDURL + "/" + XKCDJSONURL

	return getComicFromURL(url)
}

// Retrieves the JSON of the specified XKCD comic
func GetComic(num int) (*XKCD, error) {
	defer trace("GetComic")()
	url := fmt.Sprintf("%s/%d/%s", XKCDURL, num, XKCDJSONURL)

	return getComicFromURL(url)
}

func DownloadImage(imageUrl string) ([]byte, error) {
	defer trace("DownloadImage")()
	resp, err := http.Get(imageUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image download status: %v", resp.Status)
	}

	imageByte, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("reading image from body: %v", err)
	}

	return imageByte, nil
}

// Does the actual communication with the XKCD
func getComicFromURL(url string) (*XKCD, error) {
	defer trace("getComicFromURL")()
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: (%d) %s", resp.StatusCode, resp.Status)
	}

	var xkcd XKCD

	if err := json.NewDecoder(resp.Body).Decode(&xkcd); err != nil {
		return nil, fmt.Errorf("comic JSON retrieve: %v", err)
	}

	return &xkcd, nil
}