package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// version is set via -ldflags "-X main.version=x.y.z" when building a release;
// it defaults to the current development version otherwise.
var version = "1.0.0"

func main() {
	settings := parseArgs(os.Args[1:])
	displaySettings(settings)
	confirmSettings()

	findMatch := settings.findMatch()

	var bar *progressBar
	if !settings.Verbose {
		if findMatch {
			bar = newProgressBar("Finding matches", estimateTotal(settings), true)
		} else {
			bar = newProgressBar("Generating", int64(settings.NumberOfMnemonics), true)
		}
	}

	results := make(chan string)
	var fileLock sync.Mutex

	spawn := func() {
		go func() {
			results <- generateAddresses(settings)
		}()
	}

	for i := 0; i < settings.Threads; i++ {
		spawn()
	}

	complete := 0
	var resultCache strings.Builder

	flush := func() {
		fileLock.Lock()
		defer fileLock.Unlock()
		if err := saveToFile(settings, resultCache.String(), false); err != nil {
			fmt.Fprintf(os.Stderr, "\nError writing to file: %v\n", err)
		}
		resultCache.Reset()
	}

	for result := range results {
		if settings.Verbose {
			fmt.Print(result)
		} else {
			bar.increment()
		}

		if findMatch && isMatch(settings, result) {
			if err := saveToFile(settings, result, true); err != nil {
				fmt.Fprintf(os.Stderr, "\nError writing to file: %v\n", err)
			}
			if !settings.Infinite {
				fmt.Print("\n\nFound a match!\n")
				fmt.Print(result)
				os.Exit(0)
			}
		}

		resultCache.WriteString(result)
		complete++

		if complete%10 == 0 {
			flush()
		}

		if complete >= settings.NumberOfMnemonics && !settings.Infinite && !findMatch {
			if complete%10 != 0 {
				flush()
			}
			return
		}

		spawn()
	}
}

// generateAddresses creates one random mnemonic and derives
// addressesPerMnemonic addresses from it, formatted as a text block.
func generateAddresses(settings *Settings) string {
	store, err := newRandomKeyStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError generating mnemonic: %v\n", err)
		os.Exit(1)
	}

	var b strings.Builder
	b.WriteString(store.mnemonic)
	b.WriteString("\n")
	for _, addr := range store.deriveAddressesByRange(0, settings.AddressesPerMnemonic) {
		b.WriteString(addr)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}
