package main

/*
	1. Imgur Integration (Done)
	2. Wallpaper Setting
	3. Debugging
*/

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/reujab/wallpaper"
	"github.com/turnage/graw/reddit"
)

var loc = "/home/pro/Dropbox/Code/golang/snoowall/Wallpapers/"
var datafile = "data"
var name = "info.agent"
var path = fmt.Sprintf("%s%s", loc, name)
var subreddit = "gonenatural"

func saveWall(b []byte) error {
	timestamp := time.Now()
	filename := fmt.Sprintf("%s%s_%s.jpg", loc, subreddit, timestamp.Format("2006-01-02_15-04-05"))
	err := ioutil.WriteFile(filename, b, 0600)
	if err == nil {
		fmt.Println("Saved")
	}
	return err
}

func setWall(file string) error {
	background, err := wallpaper.Get()
	if err != nil {
		fmt.Println("[DEBUG] Can't find previous wallpaper:", err)
	}
	fmt.Println("Current wallpaper:", background)
	err = wallpaper.SetFromURL(file)
	if err == nil {
		fmt.Println("Wallpaper set!")
	}
	return err
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	rate := 5 * time.Second
	var after string
	script, err := reddit.NewScript("graw:snoowall:0.3.1 by /u/psychemerchant", rate)
	if err != nil {
		fmt.Println("[DEBUG] Failed to create script handle: ", err)
		return
	}
	harvest, err := script.Listing(fmt.Sprintf("/r/%s", subreddit), after)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to fetch /r/%s: %s", subreddit, err)
		return
	}
	post := harvest.Posts[rand.Intn(20)]
	// str := fmt.Sprintf("Harvest:\n %#v", harvest.Posts[1])
	// ioutil.WriteFile("harvest", []byte(str), 0600)
	// bin, _ := ioutil.ReadFile(datafile)
	// after = string(bin)
	// fmt.Println("After:", after)
	// ioutil.WriteFile(datafile, []byte(post.Name), 0600)
	fmt.Printf("[Title]: %s\n[URL]: %s\n", post.Title, post.URL)
	// fmt.Printf("[Type]: %s - %s - %s\n", post.Media.OEmbed.Type, post.Media.OEmbed.ProviderName, post.Media.OEmbed.ProviderURL)
	// fmt.Printf("%+v", post)
	err = setWall(post.URL)
	if err != nil {
		fmt.Println("[DEBUG] Wallpaper setting error:", err)
	}
	resp, err := http.Get(post.URL)
	if err != nil || post.IsRedditMediaDomain == false {
		fmt.Println("[DEBUG]: Couldn't fetch resource:", post.URL, err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	saveWall(body)

}
