package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func GeTTakeScreenshot(
	pageURL string,
	imageURL string,
	fileName string,
) error {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte

	js := fmt.Sprintf(`
		(function () {
			const target = %q;

			function normalize(url) {
				try {
					return new URL(url).pathname;
				} catch (e) {
					return url;
				}
			}

			const targetPath = normalize(target);

			const imgs = document.querySelectorAll('img');

			for (const img of imgs) {

				const src = img.currentSrc || img.src;

				if (normalize(src) === targetPath) {

					const isBroken =
						img.complete &&
						img.naturalWidth === 0;

					if (isBroken) {
						img.style.border = '5px solid red';
						img.style.boxSizing = 'border-box';
						img.title = 'BROKEN IMAGE';
					} else {
						img.style.border = '5px solid green';
						img.style.boxSizing = 'border-box';
						img.title = 'VALID IMAGE';
					}

					img.scrollIntoView({
						block: 'center',
						inline: 'center'
					});

					break;
				}
			}
		})();
	`, imageURL)

	err := chromedp.Run(
		ctx,

		// Desktop viewport
		chromedp.EmulateViewport(1920, 1080),

		chromedp.Navigate(pageURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),

		// Allow images/lazy loading
		chromedp.Sleep(3*time.Second),

		// Highlight image
		chromedp.Evaluate(js, nil),

		// Wait after scrolling
		chromedp.Sleep(1*time.Second),

		// Capture only visible desktop screen
		chromedp.CaptureScreenshot(&buf),
	)

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, buf, 0644)
}
