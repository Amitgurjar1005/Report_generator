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

type PagePerformance struct {
	ResponseTime   string `json:"response_time"`
	TTFB           string `json:"ttfb"`
	PageSize       int64  `json:"page_size"`
	HTMLSize       int64  `json:"html_size"`
	CSSSize        int64  `json:"css_size"`
	JavaScriptSize int64  `json:"javascript_size"`
	ImageSize      int64  `json:"image_size"`

	TotalRequests      int    `json:"total_requests"`
	CompressionEnabled bool   `json:"compression_enabled"`
	CompressionType    string `json:"compression_type"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func GetPagePerformance(siteURL string) (*PagePerformance, error) {

	result := &PagePerformance{}

	base, err := url.Parse(siteURL)
	if err != nil {
		return nil, err
	}

	queue := []string{siteURL}
	visited := make(map[string]bool)

	maxPages := 20

	for len(queue) > 0 && len(visited) < maxPages {

		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true

		start := time.Now()
		req, err := http.NewRequest("GET", current, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

		req.Header.Set("Accept",
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

		req.Header.Set("Accept-Language",
			"en-US,en;q=0.9")

		req.Header.Set("Referer",
			"https://www.google.com/")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			continue
		}

		if len(visited) == 1 {
			result.ResponseTime = time.Since(start).String()
			result.TTFB = result.ResponseTime
		}

		result.HTMLSize += int64(len(body))
		result.TotalRequests++

		encoding := resp.Header.Get("Content-Encoding")
		if encoding != "" {
			result.CompressionEnabled = true
			result.CompressionType = encoding
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			continue
		}

		pageBase, err := url.Parse(current)
		if err != nil {
			continue
		}

		doc.Find(`link[rel="stylesheet"]`).Each(func(i int, s *goquery.Selection) {

			href, ok := s.Attr("href")
			if !ok {
				return
			}

			result.CSSSize += getResourceSize(pageBase, href)
			result.TotalRequests++
		})

		doc.Find(`script[src]`).Each(func(i int, s *goquery.Selection) {

			src, ok := s.Attr("src")
			if !ok {
				return
			}

			result.JavaScriptSize += getResourceSize(pageBase, src)
			result.TotalRequests++
		})

		imgCount := 0

		doc.Find(`img[src]`).Each(func(i int, s *goquery.Selection) {

			if imgCount >= 20 {
				return
			}

			src, ok := s.Attr("src")
			if !ok {
				return
			}

			result.ImageSize += getResourceSize(pageBase, src)

			result.TotalRequests++

			imgCount++
		})

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

			fullURL := pageBase.ResolveReference(u).String()

			parsed, err := url.Parse(fullURL)
			if err != nil {
				return
			}

			if parsed.Host == base.Host &&
				!visited[fullURL] {

				queue = append(queue, fullURL)
			}
		})
	}

	result.PageSize =
		result.HTMLSize +
			result.CSSSize +
			result.JavaScriptSize +
			result.ImageSize

	return result, nil
}

func getResourceSize(baseURL *url.URL, resource string) int64 {

	resourceURL, err := baseURL.Parse(resource)
	if err != nil {
		return 0
	}

	resp, err := client.Head(resourceURL.String())
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.ContentLength > 0 {
		return resp.ContentLength
	}

	return 0
}
