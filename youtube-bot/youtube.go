package main

import (
	"fmt"
	"strings"

	youtubescraper "github.com/PChaparro/go-youtube-scraper"
)

type VideoData struct {
	Title string
	Url   string
}

func init() {
	youtubescraper.GetVideosData("", "youtube", 1, 1, false)
}

func searchVideos(query string) (*[]VideoData, string, error) {
	res, err := youtubescraper.GetVideosData("", query, 10, 10, false)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	data := make([]VideoData, 0, 10)
	videoList := ""

	for i, video := range res.Videos {
		data = append(data, VideoData{
			Title: video.Title,
			Url:   video.Url,
		})

		videoList += fmt.Sprintf("%d. %s\n\n", i+1, video.Title)
	}

	return &data, videoList, nil
}

func embedVideoPlayer(url string) string {
	embedLink := strings.Split(strings.Replace(url, "/watch?v=", "/embed/", 1), "&pp=")[0]

	return `<iframe width="100%" height="350px" src="` + embedLink + `" title="YouTube video player" style=
	"border-radius:0.3rem;" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>`
}
