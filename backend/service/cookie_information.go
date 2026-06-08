package service

import (
	"net/http"
	"time"
)

type CookieInfo struct {
	Name     string `json:"name"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
	Expires  string `json:"expires"`
	Value    string `json:"value"`
}

func GetCookieInformation(siteurl string) (*CookieInfo, error) {
	result := &CookieInfo{}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", siteurl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		result.Name = c.Name
		result.Secure = c.Secure
		result.HttpOnly = c.HttpOnly
		result.Expires = c.Expires.String()
		result.Value = c.Value
	}

	return result, nil
}
