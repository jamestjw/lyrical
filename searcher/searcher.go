package searcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

var searchService *youtube.SearchService

type YoutubeService struct {
	searchService        *youtube.SearchService
	playlistItemsService *youtube.PlaylistItemsService
}

// InitialiseYoutubeService initialises a youtube search object
// with the given api key.
func InitialiseYoutubeService(apiKey string) *YoutubeService {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	searchService := youtubeService.Search
	playlistItemsService := youtubeService.PlaylistItems
	log.Println("Sucessfully initialised Youtube Search Service.")
	return &YoutubeService{searchService: searchService, playlistItemsService: playlistItemsService}
}

// GetVideoID returns the first youtubeID of a
// video that matches the query
func (s *YoutubeService) GetVideoID(query string) (youtubeID string, err error) {
	call := s.searchService.List("id, snippet").
		Type("video").
		Q(query).
		MaxResults(1)

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

func (s *YoutubeService) GetVideoIDs(playlistID string, maxResults int64) (youtubeIDs []string, err error) {
	call := s.playlistItemsService.List("contentDetails").
		PlaylistId(playlistID).
		MaxResults(maxResults)

	res, err := call.Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(res.Items) == 0 {
		err = fmt.Errorf("No playlist could be found with the ID: %s", playlistID)
		return
	}

	for _, item := range res.Items {
		if item.ContentDetails != nil {
			youtubeIDs = append(youtubeIDs, item.ContentDetails.VideoId)
		}
	}
	return
}
