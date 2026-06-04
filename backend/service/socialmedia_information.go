package service

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SocialLinks struct {
	FacebookURL  string `json:"facebook_url"`
	InstagramURL string `json:"instagram_url"`
	LinkedInURL  string `json:"linkedin_url"`
	TwitterURL   string `json:"twitter_url"`
	YouTubeURL   string `json:"youtube_url"`
	GitHubURL    string `json:"github_url"`
	TelegramURL  string `json:"telegram_url"`
}

func GetSocialMediaInformation(siteURL string) (*SocialLinks, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	result := &SocialLinks{}

	queue := []string{siteURL}
	visited := map[string]bool{}

	maxPages := 10

	baseURL, _ := url.Parse(siteURL)

	for len(queue) > 0 && len(visited) < maxPages {

		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		resp, err := client.Get(current)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			continue
		}

		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {

			href, ok := s.Attr("href")
			if !ok || href == "" {
				return
			}

			href = strings.TrimSpace(href)

			if result.FacebookURL == "" && strings.Contains(href, "facebook.com") {
				result.FacebookURL = href
			}
			if result.InstagramURL == "" && strings.Contains(href, "instagram.com") {
				result.InstagramURL = href
			}
			if result.LinkedInURL == "" && strings.Contains(href, "linkedin.com") {
				result.LinkedInURL = href
			}
			if result.TwitterURL == "" &&
				(strings.Contains(href, "twitter.com") || strings.Contains(href, "x.com")) {
				result.TwitterURL = href
			}
			if result.YouTubeURL == "" &&
				(strings.Contains(href, "youtube.com") || strings.Contains(href, "youtu.be")) {
				result.YouTubeURL = href
			}
			if result.GitHubURL == "" && strings.Contains(href, "github.com") {
				result.GitHubURL = href
			}
			if result.TelegramURL == "" &&
				(strings.Contains(href, "t.me") || strings.Contains(href, "telegram.me")) {
				result.TelegramURL = href
			}

			u, err := url.Parse(href)
			if err != nil {
				return
			}

			fullURL := baseURL.ResolveReference(u)

			if fullURL.Host == baseURL.Host {
				if !visited[fullURL.String()] {
					queue = append(queue, fullURL.String())
				}
			}
		})
	}

	return result, nil
}
