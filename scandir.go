package scanfile
import (
	"fmt"
	"os"
	"time"
)

type (
	// WatchFileOpts ..
	WatchFileOpts struct {
		Listener   func(msg string, changed bool)
		Interval   int64
		WaitExists bool
	}
)

// WatchFile ...
func WatchFile(file string, opts *WatchFileOpts) error {
	var interval int64
	var lastMod int64
	firstTime := true

	if opts.Listener != nil {
		if opts.Interval > 0 {
			interval = opts.Interval
		}
	}

	for {
		f, err := os.Stat(file)
		if os.IsNotExist(err) {
			if opts != nil {
				if opts.WaitExists {
					if opts.Listener != nil {
						opts.Listener(fmt.Sprintf("File doesn't exists. Waiting file being created: %s", file), false)
					}
					firstTime = false
					time.Sleep(time.Duration(interval) * time.Second)
					continue
				} else {
					return fmt.Errorf("File doesn't exist: %s", file)
				}
			} else {
				return fmt.Errorf("File doesn't exist: %s", file)
			}
		} else {
			if err != nil {
				return fmt.Errorf("Error watching file: %s", file)
			}
		}

		if opts.Listener != nil && firstTime {
			opts.Listener(fmt.Sprintf("Watching %s ...", file), false)
		}

		if lastMod == 0 {
			lastMod = f.ModTime().Unix()
		} else {
			tempLastMod := f.ModTime().Unix()
			if tempLastMod != lastMod {
				if opts.Listener != nil {
					opts.Listener(fmt.Sprintf("%s has changed", file), true)
				}
			}
			lastMod = tempLastMod
		}

		firstTime = false
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
