package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strconv"
)

func capturescreen() int {
	var buf []byte
	var fn string
	var frames int
	var executed *runtime.RemoteObject

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	state := 1
	capWhite := make(chan bool, 1)
	capBlack := make(chan bool, 1)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			frame := ev.Args[0]
			total := ev.Args[1]
			if state == 1 {
				fn = fmt.Sprintf("tmp/W%04s.png", string(frame.Value))
			} else {
				fn = fmt.Sprintf("tmp/B%04s.png", string(frame.Value))

			}
			fmt.Printf("%s/%s\n", string(frame.Value), string(total.Value))
			go func() {
				chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))

				if err := os.WriteFile(fn, buf, 0o644); err != nil {
					log.Fatal(err)
				}
				if string(frame.Value) == string(total.Value) {
					if state == 1 {
						state = 2
						close(capWhite)
					} else {
						frames, _ = strconv.Atoi(string(total.Value))
						close(capBlack)

					}
				}
				chromedp.Run(ctx, chromedp.Evaluate(`Elusien.frame_busy=0;`, &executed))
			}()
		}
	})

	path, _ := os.Getwd()
	var filename string
	_ = path
	filename = fmt.Sprintf(`file://%s/%s`, path, os.Args[1])

	fmt.Println("White:")
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(filename),
		chromedp.SetAttributeValue("/html/body", "style", "background-color:#FFF;"),
	); err != nil {
		log.Fatal(err)
	}
	<-capWhite

	fmt.Println("Black:")
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Reload(),
		chromedp.SetAttributeValue("/html/body", "style", "background-color:#000;"),
	); err != nil {
		log.Fatal(err)
	}
	<-capBlack
	return (frames)
}
