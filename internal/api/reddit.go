package api

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/random"
)

var imageURLPattern = regexp.MustCompile(`\.(jpg|png|gif|jpeg|webp)($|\?)`)

type RedditPost struct {
	Title    string
	ImageURL string
}

func NewRedditPost(title, imageURL string) *RedditPost {
	return &RedditPost{
		Title:    title,
		ImageURL: imageURL,
	}
}

// GetImageFromSubreddit fetches a random image post from the specified subreddit.
func GetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	return GetImageFromSubredditWithRetry(subreddit, 5)
}

// GetImageFromSubredditWithRetry fetches a random image post from the specified subreddit with retry functionality.
// It will retry up to maxRetries times with exponential backoff.
func GetImageFromSubredditWithRetry(subreddit string, maxRetries int) (*RedditPost, bool) {
	subreddit = strings.TrimSpace(subreddit)
	subreddit = strings.TrimPrefix(subreddit, "r/")

	if subreddit == "" {
		return nil, false
	}

	var post *RedditPost
	var success bool

	for attempt := range maxRetries {
		if attempt > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		post, success = attemptGetImageFromSubreddit(subreddit)
		if success {
			return post, true
		}
	}

	return nil, false
}

// attemptGetImageFromSubreddit makes a single attempt to fetch an image from a subreddit
func attemptGetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	score := random.Int(50, 1000)

	// before is a random time in the last 24 days
	before := time.Now().Add(-time.Duration(random.Int(0, 24)) * 24 * time.Hour).UTC().Unix()

	url := fmt.Sprintf("https://api.pullpush.io/reddit/submission/search?html_decode=True&subreddit=%s&before=%d&score=>%d&size=3", subreddit, before, score)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		return nil, false
	}

	if !resp.IsSuccess() {
		return nil, false
	}

	data := resp.Json().Get("data")
	if data.Len() == 0 {
		return nil, false
	}

	// Shuffle the data to randomize post selection
	data = jester.New(random.Shuffle(data.MustSlice()))

	return findImagePost(data)
}

func findImagePost(data *jester.Data) (*RedditPost, bool) {
	for post := range data.Iterator() {
		title := post.Get("title").MustString()
		postHint := post.Get("post_hint").MustString()
		url := post.Get("url").MustString()

		if postHint == "image" || imageURLPattern.MatchString(url) {
			return NewRedditPost(title, url), true
		}
	}

	return nil, false
}
