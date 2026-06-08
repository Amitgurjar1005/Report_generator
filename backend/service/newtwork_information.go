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

type NetworkInfo struct {
	APIEndpoints    []string          `json:"api_endpoints"`
	GraphQLEndpoint string            `json:"graphql_endpoint"`
	WebSocketUsage  bool              `json:"websocket_usage"`
	CORSHeaders     map[string]string `json:"cors_headers"`
	AllowedMethods  string            `json:"allowed_methods"`
}

func GetNetworkInformation(siteURL string) (*NetworkInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	result := &NetworkInfo{
		CORSHeaders: make(map[string]string),
	}
	visited := make(map[string]bool)

	queue := []string{siteURL}
	max := 20
	for len(queue) > 0 && len(visited) <= max {
		current := queue[0]
		queue = queue[1:]
		visited[current] = true
		resp, err := client.Get(current)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		resp.Body.Close()
		html := string(body)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		apiRegex := regexp.MustCompile(`(?i)/api/[a-zA-Z0-9/_-]+`)
		matches := apiRegex.FindAllString(html, -1)
		exist := make(map[string]bool)
		for _, m := range matches {
			if !exist[m] {
				exist[m] = true
				result.APIEndpoints = append(result.APIEndpoints, m)
			}
		}
		if strings.Contains(strings.ToLower(html), "graphql") {
			result.GraphQLEndpoint = "/graphql"
		}
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
			parsed, err := url.Parse(fullurl)
			if err != nil {
				return
			}
			if parsed.Host != base.Host {
				return
			}
			if !visited[fullurl] {
				queue = append(queue, fullurl)
			}

		})

		corsKeys := []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Credentials",
		}
		for _, k := range corsKeys {
			if v := resp.Header.Get(k); v != "" {
				result.CORSHeaders[k] = v
			}
		}
		req, err := http.NewRequest("OPTIONS", current, nil)
		if err == nil {
			req.Header.Set("Origin", siteURL)
			req.Header.Set("Access-Control-Request-Method", "GET")

			res, err := client.Do(req)
			if err == nil {
				cors := res.Header.Get("Access-Control-Allow-Methods")
				allow := res.Header.Get("Allow")

				if cors != "" {
					result.AllowedMethods = cors
				} else {
					result.AllowedMethods = allow
				}

				res.Body.Close()
			}
		}
	}
	return result, nil
}
