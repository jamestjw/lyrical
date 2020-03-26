package searcher

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

var searchService *youtube.SearchService

// InitialiseSearchService initialises a youtube search object
// with the given api key.
func InitialiseSearchService(apiKey string) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	searchService = youtubeService.Search
	log.Println("Sucessfully initialised Youtube Search Service.")
}

// GetVideoID returns the first youtubeID of a
// video that matches the query
func GetVideoID(query string) (youtubeID string, err error) {
	call := searchService.List("id, snippet").
		Type("video").
		Q(query).
		MaxResults(1).
		VideoDuration("short")

	res, err := call.Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(res.Items) == 0 {
		err = fmt.Errorf("no results could be found for your query: %s", query)
		return
	}

	result := res.Items[0]
	youtubeID = result.Id.VideoId
	return
}
