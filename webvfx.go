//use Chromedp to take screenshot of WebVfx animations
//add a frame_busy flag to sync frame
//use ffmpeg to convert captured picture to film, finally.

package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/exec"
)

func main() {
	var buf []byte
	var executed *runtime.RemoteObject
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	_, _ = exec.Command(`sh`, `rmtmp.sh`).Output()

	gotExp := make(chan bool, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			frame := ev.Args[0]
			total := ev.Args[1]
			fn := fmt.Sprintf("tmp/f%04s.png", string(frame.Value))
			fmt.Println(string(frame.Value), string(total.Value))
			go func() {
				chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))

				if err := os.WriteFile(fn, buf, 0o644); err != nil {
					log.Fatal(err)
				}
				if string(frame.Value) == string(total.Value) {
					_, _ = exec.Command(`sh`, `convert.sh`).Output()

					_, _ = exec.Command(`mv`, `t.webm`, os.Args[1]+`.webm`).Output()
					os.Exit(0)
				}
				chromedp.Run(ctx, chromedp.Evaluate(`Elusien.frame_busy=0;`, &executed))
			}()
		}
	})
	path, _ := os.Getwd()
	var filename string
	_ = path
	filename = fmt.Sprintf(`file://%s/%s`, path, os.Args[1])
	fmt.Println(filename)
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(filename),
	); err != nil {
		log.Fatal(err)
	}
	<-gotExp
}
