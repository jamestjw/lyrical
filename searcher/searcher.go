package searcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

var searchService *youtube.SearchService

type SearchService struct {
	searchService *youtube.SearchService
}

// InitialiseSearchService initialises a youtube search object
// with the given api key.
func InitialiseSearchService(apiKey string) *SearchService {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	searchService = youtubeService.Search
	log.Println("Sucessfully initialised Youtube Search Service.")
	return &SearchService{searchService: searchService}
}

// GetVideoID returns the first youtubeID of a
// video that matches the query
func (s *SearchService) GetVideoID(query string) (youtubeID string, err error) {
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
