package scrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"valentin-lvov/1x-parser/cache"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/redis/go-redis/v9"
)

func MakeConnectionAndLoad(url string) (*context.Context, context.CancelFunc, error) {

	ctx, cancel := chromedp.NewContext(context.Background())

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// chromedp.Sleep(10*time.Second),
		chromedp.WaitVisible(`div#allBetsTable`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			/*click on all collapsed headers in the table to load the necessary data*/
			var err error
			script_click_collapsed := `
                var elements = document.querySelectorAll('div.bet-title.bet-title_justify.min');
                elements.forEach(function(element) {
                    element.click();
                });
            `
			err = chromedp.Evaluate(script_click_collapsed, nil).Do(ctx)
			if err != nil {
				return err
			}

			return nil

		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// selector := "div.bets_content.betsscroll > div.iScrollVerticalScrollbar.iScrollLoneScrollbar > div.iScrollIndicator" // scroll bar elelemt, might need later
			selector := "div#allBetsTable"
			/*call wheel event to scroll the table and load more content until there is no more content appears*/
			script_scroll := `
						var event = new WheelEvent('wheel', {
							deltaX: 0,
							deltaY: 10000 // Adjust the deltaY to simulate scroll amount
						});
						document.querySelector("div.bets_content.betsscroll").dispatchEvent(event);
						`

			var prevstyle string
			var currstyle string
			var betcount int

			for {
				// Data about table size is kept inside table.style.transform
				err := chromedp.Evaluate(`document.querySelector('`+selector+`').style.transform`, &prevstyle).Do(ctx)
				if err != nil {
					return err
				}

				fmt.Println(prevstyle)
				//this next part can be removed, it is just a log basically to see that scrolling bets table changed anything
				err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Current bet count:%d\n", betcount)

				// Scroll the element
				err = chromedp.Evaluate(script_scroll, nil).Do(ctx)
				if err != nil {
					return err
				}

				// Wait for dynamic content to load
				time.Sleep(3 * time.Second)

				// Check the scroll position again
				err = chromedp.Evaluate(`document.querySelector('`+selector+`').style.transform`, &currstyle).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Println(currstyle)

				err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("bet count after wheel event:%d\n", betcount)

				// Break the loop if the scroll position hasn't changed
				if prevstyle == currstyle {
					break
				}
				// err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
			}
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	return &ctx, cancel, nil

}

func GetContentFromSelector(ctx *context.Context) (map[string]string, error) {
	/* bet name and coefficient is currently kept inside div.bet-inner as text*/
	var divs []*cdp.Node
	var span1Text, span2Text string
	selector := "div.bet-inner"

	err := chromedp.Run(*ctx,
		chromedp.Nodes(selector, &divs, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}
	spanTextMap := make(map[string]string)

	for _, div := range divs {
		err = chromedp.Run(*ctx,
			chromedp.Text(`span.bet_type`, &span1Text, chromedp.ByQuery, chromedp.FromNode(div)),
			chromedp.Text(`span.koeff`, &span2Text, chromedp.ByQuery, chromedp.FromNode(div)),
		)
		if err != nil {
			log.Println("Error extracting text from spans:", err)
			continue
		}
		spanTextMap[span1Text] = span2Text
	}

	// fmt.Printf("Succesfully scrapped content!\n")
	return spanTextMap, nil
}
func SaveToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func AreMapsEqual(map_p1, map_p2 *map[string]string) bool {
	map1 := *map_p1
	map2 := *map_p2
	if len(map1) != len(map2) {
		return false
	}
	for key, value_1 := range map1 {
		value_2, exists := map2[key]
		if !exists || value_2 != value_1 {
			return false
		}
	}
	return true
}

func copyMap(map1, map2 *map[string]string) {
	mp1 := *map1
	mp2 := *map2
	for k, v := range mp1 {
		mp2[k] = v
	}
	return

}

func ScrapWebsite(url string) (map[string]string, error) {

	var result map[string]string
	var ctx *context.Context
	var err error

	ctx, cancel, err := MakeConnectionAndLoad(url)
	defer cancel()

	if err != nil {
		log.Fatal("Error creating ChromeDP context:", err)
		return nil, err
	}
	result, err = GetContentFromSelector(ctx)
	if err != nil {
		log.Fatal("Error getting the content:", err)
		return nil, err
	}
	return result, nil
}

func TrackWebsite(url string, duration time.Duration, interval time.Duration, rdb *redis.Client) error {
	/*track url for duration. Check the website every inetrval and store it in redis db rdb*/
	var ctx *context.Context
	var err error
	var currentContent *map[string]string
	// var previousContent *map[string]string
	// previousContent = &map[string]string{}

	ctx, cancel, err := MakeConnectionAndLoad(url) // load website inside headless chrome browser
	defer cancel()
	if err != nil {
		log.Fatal("Error creating ChromeDP context:", err)
		return err
	}

	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		res, err := GetContentFromSelector(ctx)
		currentContent = &res
		if err != nil {
			log.Fatal("Error getting the content:", err)
			return err
		}
		cache.StoreInRedis(rdb, url, *currentContent, duration)

		// if !AreMapsEqual(currentContent, previousContent) { // could replace this with just cache.StoreInRedis to not keep extra copy here
		// 	// but this will obviously increase the number of cache accesses
		// 	cache.StoreInRedis(rdb, url, *currentContent, duration)
		// 	previousContent = currentContent
		// }

	}
	return nil
}
