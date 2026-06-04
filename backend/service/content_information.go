package service

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ContentInfo struct {
	TotalWords         int      `json:"total_words"`
	TotalParagraphs    int      `json:"total_paragraphs"`
	TotalImages        int      `json:"total_images"`
	TotalSpan          int      `json:"total_span"`
	TotalVideos        int      `json:"total_videos"`
	TotalInternalLinks int      `json:"total_internal_links"`
	TotalExternalLinks int      `json:"total_external_links"`
	BrokenLinks        int      `json:"broken_links"`
	BrokenImages       int      `json:"broken_images"`
	EmailsFound        []string `json:"emails_found"`
	PhonesFound        string   `json:"phones_found"`
}

var client2 = &http.Client{
	Timeout: 5 * time.Second,
}

var brokenCache = map[string]bool{}

func GetContentInformation(siteurl string) (*ContentInfo, error) {

	result := &ContentInfo{}

	queue := []string{siteurl}
	visited := make(map[string]bool)

	maxPages := 20

	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`(\+?\d[\d\s\-]{8,}\d)`)

	for len(queue) > 0 && len(visited) < maxPages {

		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		resp, err := client2.Get(current)
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

		text := doc.Text()

		result.TotalWords += len(strings.Fields(text))
		result.TotalParagraphs += doc.Find("p").Length()
		result.TotalImages += doc.Find("img").Length()
		result.TotalSpan += doc.Find("span").Length()
		result.TotalVideos += doc.Find("video").Length()

		base, _ := url.Parse(current)

		linkCheckLimit := 10
		imgCheckLimit := 10

		linkChecked := 0
		imgChecked := 0

		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {

			href, _ := s.Attr("href")

			if href == "" ||
				strings.HasPrefix(href, "#") ||
				strings.HasPrefix(href, "mailto:") ||
				strings.HasPrefix(href, "javascript:") {
				return
			}

			u, err := url.Parse(href)
			if err != nil {
				return
			}

			fullURL := base.ResolveReference(u).String()

			parsed, err := url.Parse(fullURL)
			if err != nil {
				return
			}

			if parsed.Host == base.Host {
				result.TotalInternalLinks++

				if !visited[fullURL] {
					queue = append(queue, fullURL)
				}
			} else {
				result.TotalExternalLinks++
			}

			if linkChecked < linkCheckLimit {
				if isBroken(fullURL) {
					result.BrokenLinks++
				}
				linkChecked++
			}
		})

		doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {

			if imgChecked >= imgCheckLimit {
				return
			}

			src, ok := s.Attr("src")
			if !ok || src == "" {
				return
			}

			imgURL, err := base.Parse(src)
			if err != nil {
				return
			}

			if isBroken(imgURL.String()) {
				result.BrokenImages++
			}

			imgChecked++
		})

		result.EmailsFound = append(
			result.EmailsFound,
			emailRegex.FindAllString(text, -1)...,
		)

		if result.PhonesFound == "" {
			result.PhonesFound = phoneRegex.FindString(text)
		}
	}

	return result, nil
}

func isBroken(link string) bool {

	if v, ok := brokenCache[link]; ok {
		return v
	}

	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		brokenCache[link] = true
		return true
	}

	resp, err := client2.Do(req)
	if err != nil {
		brokenCache[link] = true
		return true
	}
	defer resp.Body.Close()

	broken := resp.StatusCode >= 400
	brokenCache[link] = broken

	return broken
}
