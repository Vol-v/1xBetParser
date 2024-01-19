package main

import (
	"context"
	"fmt"
	"log"
	"valentin-lvov/1x-parser/scrapper"
)

// func makeConnectionAndLoad(url string) (*context.Context, context.CancelFunc, error) {

// 	ctx, cancel := chromedp.NewContext(context.Background())

// 	err := chromedp.Run(ctx,
// 		chromedp.Navigate(url),
// 		// chromedp.Sleep(10*time.Second),
// 		chromedp.WaitVisible(`div#allBetsTable`, chromedp.ByQuery),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			var err error
// 			script_click_collapsed := `
//                 var elements = document.querySelectorAll('div.bet-title.bet-title_justify.min');
//                 elements.forEach(function(element) {
//                     element.click();
//                 });
//             `
// 			err = chromedp.Evaluate(script_click_collapsed, nil).Do(ctx)
// 			if err != nil {
// 				return err
// 			}

// 			return nil

// 		}),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			// selector := "div.bets_content.betsscroll > div.iScrollVerticalScrollbar.iScrollLoneScrollbar > div.iScrollIndicator" // CSS selector for the element
// 			selector := "div#allBetsTable"
// 			script_scroll := `
// 						var event = new WheelEvent('wheel', {
// 							deltaX: 0,
// 							deltaY: 10000 // Adjust the deltaY to simulate scroll amount
// 						});
// 						document.querySelector("div.bets_content.betsscroll").dispatchEvent(event);
// 						`

// 			var prevstyle string
// 			var currstyle string
// 			var betcount int

// 			for {
// 				// Get the current scroll position of the element
// 				err := chromedp.Evaluate(`document.querySelector('`+selector+`').style.transform`, &prevstyle).Do(ctx)
// 				if err != nil {
// 					return err
// 				}

// 				fmt.Println(prevstyle)

// 				err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
// 				if err != nil {
// 					return err
// 				}
// 				fmt.Printf("Current bet count:%d\n", betcount)

// 				// Scroll the element
// 				err = chromedp.Evaluate(script_scroll, nil).Do(ctx)
// 				if err != nil {
// 					return err
// 				}

// 				// Wait for dynamic content to load
// 				time.Sleep(3 * time.Second)

// 				// Check the scroll position again
// 				err = chromedp.Evaluate(`document.querySelector('`+selector+`').style.transform`, &currstyle).Do(ctx)
// 				if err != nil {
// 					return err
// 				}
// 				fmt.Println(currstyle)

// 				err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
// 				if err != nil {
// 					return err
// 				}
// 				fmt.Printf("bet count after wheel event:%d\n", betcount)

// 				// Break the loop if the scroll position hasn't changed
// 				if prevstyle == currstyle {
// 					break
// 				}
// 				// err = chromedp.Evaluate(`document.querySelectorAll('div.bet-inner').length`, &betcount).Do(ctx)
// 			}
// 			return nil
// 		}),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, nil, err
// 	}
// 	return &ctx, cancel, nil

// }

// func getContent(ctx *context.Context, selector string) (map[string]string, error) {
// 	var divs []*cdp.Node
// 	var span1Text, span2Text string

// 	err := chromedp.Run(*ctx,
// 		chromedp.Nodes(selector, &divs, chromedp.ByQueryAll),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	spanTextMap := make(map[string]string)

// 	for _, div := range divs {
// 		err = chromedp.Run(*ctx,
// 			chromedp.Text(`span.bet_type`, &span1Text, chromedp.ByQuery, chromedp.FromNode(div)),
// 			chromedp.Text(`span.koeff`, &span2Text, chromedp.ByQuery, chromedp.FromNode(div)),
// 		)
// 		if err != nil {
// 			log.Println("Error extracting text from spans:", err)
// 			continue
// 		}
// 		spanTextMap[span1Text] = span2Text
// 	}

// 	return spanTextMap, nil
// }

// func saveToFile(filename, content string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(content)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func main() {
	url := "https://1xbet.com/en/live/football/21317-africa-cup-of-nations/504058217-cape-verde-mozambique"

	var result map[string]string
	var ctx *context.Context
	var err error

	ctx, cancel, err := scrapper.MakeConnectionAndLoad(url)
	defer cancel()

	if err != nil {
		log.Fatal("Error creating ChromeDP context:", err)
		return
	}
	result, err = scrapper.GetContent(ctx, "div.bet-inner")
	if err != nil {
		log.Fatal("Error getting the content:", err)
		return
	}

	fmt.Println(len(result))

}
