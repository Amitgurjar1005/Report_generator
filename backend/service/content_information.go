package service

import (
	"bytes"
	"fmt"
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
	seenInternal := map[string]bool{}
	seenExternal := map[string]bool{}
	result := &ContentInfo{}

	queue := []string{siteurl}
	visited := make(map[string]bool)

	maxPages := 100
	brokenLinkSeen := make(map[string]bool)
	brokenImageSeen := make(map[string]bool)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`(\+?\d[\d\s\-]{8,}\d)`)

	for len(queue) > 0 && len(visited) < maxPages {

		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

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

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			continue
		}

		text := doc.Text()
		if strings.Contains(strings.ToLower(text), "mod_security") ||
			strings.Contains(strings.ToLower(text), "modsecurity") ||
			strings.Contains(strings.ToLower(text), "not acceptable!") ||
			strings.Contains(strings.ToLower(text), "access denied") {

			continue
		}
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
				if !seenInternal[fullURL] {
					result.TotalInternalLinks++
					seenInternal[fullURL] = true
				}

				if !visited[fullURL] {
					queue = append(queue, fullURL)
				}
			} else {
				if !seenExternal[fullURL] {
					result.TotalExternalLinks++
					seenExternal[fullURL] = true
				}
			}

			if linkChecked < linkCheckLimit {

				if isBrokenLink(fullURL) {

					if !brokenLinkSeen[fullURL] {
						result.BrokenLinks++
						brokenLinkSeen[fullURL] = true
					}
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

			if isBrokenImage(imgURL.String()) {

				if !brokenImageSeen[imgURL.String()] {
					result.BrokenImages++
					fmt.Println("CHECKING IMAGE:", imgURL.String())
					brokenImageSeen[imgURL.String()] = true
				}
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

// func isBrokenImage(link string) bool {

// 	if v, ok := brokenCache[link]; ok {
// 		return v
// 	}

// 	req, err := http.NewRequest("HEAD", link, nil)
// 	if err != nil {
// 		brokenCache[link] = true
// 		return true
// 	}

// 	resp, err := client2.Do(req)
// 	if err != nil {
// 		brokenCache[link] = true
// 		return true
// 	}
// 	defer resp.Body.Close()

// 	broken := resp.StatusCode >= 400
// 	brokenCache[link] = broken
// 	return broken
// }

// func isBrokenLink(link string) bool {

// 	if v, ok := brokenCache[link]; ok {
// 		return v
// 	}

// 	req, err := http.NewRequest("HEAD", link, nil)
// 	if err != nil {
// 		brokenCache[link] = true
// 		return true
// 	}

// 	resp, err := client2.Do(req)
// 	if err != nil {
// 		brokenCache[link] = true
// 		return true
// 	}
// 	defer resp.Body.Close()

// 	broken := resp.StatusCode >= 400

// 	if broken {
// 		fmt.Println("BROKEN LINK:", link, "Status:", resp.StatusCode)
// 	}

// 	brokenCache[link] = broken
// 	return broken
// }

func isBrokenImage(link string) bool {

	if v, ok := brokenCache[link]; ok {
		return v
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		brokenCache[link] = true
		return true
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "image/*,*/*;q=0.8")
	req.Header.Set("Range", "bytes=0-2048")

	resp, err := client2.Do(req)
	if err != nil {
		brokenCache[link] = true
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		brokenCache[link] = true
		return true
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "image") {
		brokenCache[link] = true
		return true
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err == nil {
		if strings.Contains(strings.ToLower(string(body)), "<html") {
			brokenCache[link] = true
			return true
		}
	}

	brokenCache[link] = false
	return false
}

func isBrokenLink(link string) bool {

	if v, ok := brokenCache[link]; ok {
		return v
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return true
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Range", "bytes=0-1024")

	resp, err := client2.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	if status == 404 || status == 410 {
		brokenCache[link] = true
		fmt.Println("broken", link)
		return true
	}

	if status == 403 || status == 406 || status == 429 {
		brokenCache[link] = false
		return false
	}

	brokenCache[link] = false
	return false
}
