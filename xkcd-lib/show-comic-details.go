package xkcd_lib

import "fmt"

func ShowComicDetails(XKCD *XKCD) {
	fmt.Printf("Title: %s\n", XKCD.Title)
	fmt.Printf("Number %d released on %s-%s-%s\n",
		XKCD.Number, XKCD.Year, XKCD.Month, XKCD.Day)
	fmt.Printf("Image: %s\n", XKCD.ImageAlt)
	fmt.Printf("Transcript: %s\n", XKCD.Transcript)
}
