package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func menu(usage string) {
	fmt.Print("Zenon Address Generator\n" +
		"Usage: znn-address-generator [OPTIONS]\n\n" +
		"Options\n" +
		usage + "\n")
}

func confirmSettings() {
	fmt.Print("Do you want to proceed? (Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	if line != "" && line != "y" && line != "yes" {
		os.Exit(0)
	}
}

func displaySettings(s *Settings) {
	findMatch := s.findMatch()

	fmt.Println("---------------------------------------------------")
	fmt.Printf("Threads: %d\n", s.Threads)
	if !findMatch && !s.Infinite {
		fmt.Printf("Number of Mnemonics: %d\n", s.NumberOfMnemonics)
	} else if s.Infinite {
		fmt.Println("Number of Mnemonics: infinite")
	}
	fmt.Printf("Addresses per mnemonic: %d\n", s.AddressesPerMnemonic)
	fmt.Printf("Output file: %s\n", s.Filename)
	fmt.Printf("Verbose output: %t\n", s.Verbose)
	if findMatch {
		if s.Prefix != "" {
			fmt.Printf("Matching prefix: %s\n", s.Prefix)
		}
		if s.Suffix != "" {
			fmt.Printf("Matching suffix: %s\n", s.Suffix)
		}
		fmt.Printf("Save all mnemonics to separate file: %t\n", s.SaveAll)
	}
	fmt.Println("---------------------------------------------------")
}
