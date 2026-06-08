package service

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func GeTTakeScreenshot(pageURL, fileName string) error {

	// 👇 IMPORTANT: create folder first
	os.MkdirAll("screenshots", 0755)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/snap/bin/chromium"),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, buf, 0644)
}
