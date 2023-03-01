//use Chromedp to take screenshot of WebVfx animations
//add a frame_busy flag to sync frame
//use ffmpeg to convert captured picture to film, finally.

package main

import (
	"os"
	"os/exec"
)

func main() {

	_, _ = exec.Command(`sh`, `rmtmp.sh`).Output()

	pag := capturescreen()
	regenalpha(pag)
	_, _ = exec.Command(`sh`, `convert.sh`).Output()

	_, _ = exec.Command(`mv`, `t.webm`, os.Args[1]+`.webm`).Output()

}
