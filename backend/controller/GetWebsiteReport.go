package controller

import (
	"backend/backend/service"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type request struct {
	URL string `json:"url" binding:"required"`
}

func GetWebsiteReport(c *gin.Context) {
	var r request

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var wg sync.WaitGroup

	var (
		pagePerformance any
		seoPerformance  any
		cookieInfo      any
		contentInfo     any
		networkInfo     any
		socialLinksInfo any
		riskInformation any
	)

	wg.Add(7)

	go func() {
		defer wg.Done()
		pagePerformance, _ = service.GetPagePerformance(r.URL)
	}()

	go func() {
		defer wg.Done()
		seoPerformance, _ = service.GetSeoInformation(r.URL)
	}()

	go func() {
		defer wg.Done()
		cookieInfo, _ = service.GetCookieInformation(r.URL)
	}()

	go func() {
		defer wg.Done()
		contentInfo, _ = service.GetContentInformation(r.URL)
	}()

	go func() {
		defer wg.Done()
		networkInfo, _ = service.GetNetworkInformation(r.URL)
	}()

	go func() {
		defer wg.Done()
		socialLinksInfo, _ = service.GetSocialMediaInformation(r.URL)
	}()
	go func() {
		defer wg.Done()
		riskInformation, _ = service.GetRiskInformation(r.URL)
	}()

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"page_performance": pagePerformance,
		"seo_performance":  seoPerformance,
		"cookie_info":      cookieInfo,
		"content_info":     contentInfo,
		"network_info":     networkInfo,
		"sociallink_info":  socialLinksInfo,
		"risk_infor":       riskInformation,
	})
}
