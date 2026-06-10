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

type RiskAnalysis struct {
	ExposedGit       bool `json:"exposed_git"`
	ExposedEnv       bool `json:"exposed_env"`
	BackupFiles      bool `json:"backup_files"`
	DirectoryListing bool `json:"directory_listing"`
	DebugInformation bool `json:"debug_information"`
	SensitiveFiles   bool `json:"sensitive_files"`
	PublicAdminPanel bool `json:"public_admin_panel"`
	RiskScore        int  `json:"risk_score"`
}

var client3 = &http.Client{
	Timeout: 10 * time.Second,
}

func GetRiskInformation(siteurl string) (*RiskAnalysis, error) {

	result := &RiskAnalysis{}
	queue := []string{siteurl}
	visited := map[string]bool{}
	max := 100
	for len(queue) > 0 && len(visited) < max {
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

		reader, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(reader))
		if err != nil {
			continue
		}
		pageText := strings.ToLower(doc.Text())
		if strings.Contains(pageText, "index of") {
			result.DirectoryListing = true
		}
		debugkeyword := []string{
			"stack trace",
			"exception",
			"panic:",
			"runtime error",
			"debug",
		}
		for _, v := range debugkeyword {
			if strings.Contains(pageText, v) {
				result.DebugInformation = true
				break
			}
		}
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if !ok || href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "#") {
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

			pasrsed, err := url.Parse(fullurl)
			if err != nil {
				return
			}
			if pasrsed.Host != base.Host {
				return
			}
			if !visited[fullurl] {
				queue = append(queue, fullurl)
			}
		})
	}
	baseurl, err := url.Parse(siteurl)
	if err != nil {
		return nil, err
	}
	root := baseurl.Scheme + "://" + baseurl.Host
	if geturlexit(root + "/.git/") {
		result.ExposedGit = true
	}
	if geturlexit(root + "/.env") {
		result.ExposedEnv = true
	}

	backupFiles := []string{
		"/backup.zip",
		"/backup.sql",
		"/site.zip",
		"/db.sql",
		"/website.zip",
	}
	for _, v := range backupFiles {
		if geturlexit(root + v) {
			result.BackupFiles = true
			break
		}
	}
	sensitiveFiles := []string{
		"/config.json",
		"/database.yml",
		"/settings.php",
		"/config.php",
	}
	for _, s := range sensitiveFiles {
		if geturlexit(root + s) {
			result.SensitiveFiles = true
			break
		}
	}
	adminPaths := []string{
		"/admin",
		"/login",
		"/dashboard",
		"/admin/login",
	}
	for _, a := range adminPaths {
		if geturlexit(root + a) {
			result.PublicAdminPanel = true
			break
		}
	}

	score := 0

	if result.ExposedGit {
		score += 30
	}
	if result.ExposedEnv {
		score += 30
	}
	if result.BackupFiles {
		score += 15
	}

	if result.DirectoryListing {
		score += 10
	}

	if result.DebugInformation {
		score += 5
	}

	if result.SensitiveFiles {
		score += 10
	}

	if result.PublicAdminPanel {
		score += 10
	}

	percentage := int(float64(score) / 110.0 * 100)

	result.RiskScore = percentage

	return result, nil
}

func geturlexit(siteurl string) bool {

	resp, err := client3.Get(siteurl)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
