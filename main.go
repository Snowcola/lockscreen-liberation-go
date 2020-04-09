package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/LockscreenLiberation/icon"
	"github.com/micmonay/keybd_event"
)

func main() {
	onExit := func() {
		fmt.Println("Starting onExit")
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		fmt.Println("Finished onExit")
	}
	// Should be called at the very beginning of main().
	systray.RunWithAppWindow("Lantern", 1024, 768, onReady, onExit)
}

func pressKey() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	key := 0x7E + 0xFFF
	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Select keys to be pressed
	kb.SetKeys(key)

	// Set shift to be pressed
	kb.HasSHIFT(true)

	// Press the selected keys

	err = kb.Launching()
	if err != nil {
		panic(err)
	}
	fmt.Println("pressed key")

}

func text(quit chan int) {
	for {
		select {
		case <-quit:
			fmt.Println("Disabled")
			return
		case <-time.After(3 * time.Minute):
			pressKey()
		}
	}
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Lockscreen Liberation")
	systray.SetTooltip("Lockscreen Liberation")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit")

	liberate := make(chan int)

	go func() {
		text(liberate)
	}()

	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Lockscreen Liberation")
		systray.SetTooltip("Lockscreen Liberation")

		systray.AddSeparator()
		mToggle := systray.AddMenuItem("Disable", "Toggle preventing the lockscreen")
		preventLockscreen := true
		toggle := func() {
			if preventLockscreen {
				mToggle.SetTitle("Enable")
				preventLockscreen = false
				liberate <- 0
			} else {
				mToggle.SetTitle("Disable")
				preventLockscreen = true
				go text(liberate)
			}
		}

		for {
			select {
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()

}
