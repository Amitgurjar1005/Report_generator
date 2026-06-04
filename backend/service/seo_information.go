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

type SEOInfo struct {
	Title                 string `json:"title"`
	TitleLength           int    `json:"title_length"`
	MetaDescription       string `json:"meta_description"`
	MetaDescriptionLength int    `json:"meta_description_length"`
	CanonicalURL          string `json:"canonical_url"`
	H1Count               int    `json:"h1_count"`
	H2Count               int    `json:"h2_count"`
	H3Count               int    `json:"h3_count"`
}

func GetSeoInformation(siteurl string) (*SEOInfo, error) {
	result := &SEOInfo{}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	queue := []string{siteurl}
	visited := make(map[string]bool)
	max := 20

	for len(queue) > 0 && len(visited) <= max {
		current := queue[0]
		queue = queue[1:]
		visited[current] = true

		resp, err := client.Get(current)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		reader, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(reader))
		if err != nil {
			return nil, err
		}
		result.Title = doc.Find("title").Text()

		result.TitleLength = len(result.Title)

		result.MetaDescription, _ =
			doc.Find(`meta[name="description"]`).Attr("content")

		result.MetaDescriptionLength = len(result.MetaDescription)

		result.H1Count += doc.Find("h1").Length()

		result.H2Count += doc.Find("h2").Length()

		result.H3Count += doc.Find("h3").Length()

		result.CanonicalURL, _ = doc.Find(`link[rel="canonical"]`).Attr("href")

		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if !ok || href == "" {
				return
			}
			if strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "#") {
				return
			}
			base, err := url.Parse(current)
			if err != nil {
				return
			}
			u, err := url.Parse(href)
			if err != nil {
				return
			}
			fullurl := base.ResolveReference(u).String()
			parsed, _ := url.Parse(fullurl)
			if parsed.Host != base.Host {
				return
			}

			if !visited[fullurl] {
				queue = append(queue, fullurl)
			}
		})
	}
	return result, nil
}
