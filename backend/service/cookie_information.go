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
	resp, err := client.Get(siteurl)
	if err != nil {
		return nil, err
	}
	for _, c := range resp.Cookies() {
		result.Name = c.Name
		result.Secure = c.Secure
		result.HttpOnly = c.HttpOnly
		result.Expires = c.Expires.String()
		result.Value = c.Value

	}
	return result, nil
}
