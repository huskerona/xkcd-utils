package xkcd_lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Retrieves the JSON of the latest XKCD comicl
func (xkcd *XKCD) GetLatestComic() error {
	defer trace("GetLatestComic")()
	url := XKCDURL + "/" + XKCDJSONURL

	var tempXkcd *XKCD
	var err error

	if tempXkcd, err = getComicFromURL(url); err != nil {
		return err
	}

	// &xkcd.Number has the same address as in SyncAll function when this method is called. (labeled as 3.)
	// fmt.Println("1. &xkcd: ", &xkcd.Number)

	// Warning is displayed when doing this assignment: xkcd = tempXkcd
	// Assignment to method receiver propagates on to callees but not to callers
	// Inspection info: Reports assignment to method receiver
	// When assigning a value to the method receiver it won't be reflected outside of the method itself.
	// Values will be reflected in subsequent calls from the same method.

	*xkcd = *tempXkcd

	// if copying tempXkcd to xkcd as in xkcd = tempXkcd, the &xkcd.Number will have a different address to the above
	// xkcd.Number (labeled as 1.). It would be expected to keep this address when it exits this method and
	// continues to work, but this sort of assignment is only visible here. A caller (eg. SyncAll) will have
	// the original address.
	// The new address is not applied.
	// An explanation for this is given in The Go Programming Language on page 161 (Section 6.3):
	//        Because url.Values is a map type and a map refers to its key/value pairs indirectly, any
	//        updates and deletions that url.Values.Add makes to the map elements are visible to the
	//        caller. However, as with ordinary functions, any changes a method makes to the reference
	//        itself, like setting it to nil or making it refer to a different map data structures, will not be
	//        reflected in the caller.

	//copyXkcd(xkcd, tempXkcd)

	fmt.Println("2. &xkcd: ", &xkcd.Number)

	return nil
}

// Retrieves the JSON of the specified XKCD comic
func (xkcd *XKCD) GetComic(num int) error {
	defer trace("GetComic")()
	url := fmt.Sprintf("%s/%d/%s", XKCDURL, num, XKCDJSONURL)

	var tempXkcd *XKCD
	var err error

	if tempXkcd, err = getComicFromURL(url); err != nil {
		return err
	}

	copyXkcd(xkcd, tempXkcd)

	return nil
}

func (xkcd *XKCD) DownloadImage() ([]byte, error) {
	defer trace("DownloadImage")()
	resp, err := http.Get(xkcd.ImageURL)

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

func copyXkcd(dest *XKCD, src *XKCD) {
	dest.Transcript = src.Transcript
	dest.Number = src.Number
	dest.ImageURL = src.ImageURL
	dest.Day = src.Day
	dest.Month = src.Month
	dest.Year = src.Year
	dest.Title = src.Title
	dest.SafeTitle = src.SafeTitle
	dest.ImageAlt = src.ImageAlt
	dest.Link = src.Link
	dest.News = src.News
}