package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	fp.UserAgent = "MyCustomAgent 1.0"
	feed, _ := fp.ParseURL("https://www.eldiarioar.com/rss")
	fmt.Println(feed.Items[0].Title)
	fmt.Println(feed.Items[0].Content)

	fmt.Println("hello news!")
}
