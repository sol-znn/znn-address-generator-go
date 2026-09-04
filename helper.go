package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

// isValid checks that input only contains Bech32 charset characters, and
// that the combined prefix+suffix length does not exceed 37.
func isValid(s *Settings, input string) {
	for _, c := range input {
		if !strings.ContainsRune(bech32Charset, c) {
			fmt.Printf("Error! %s is invalid\n", input)
			fmt.Println("Valid characters are [a-z0-9], excluding 1, b, i, o.")
			os.Exit(1)
		}
	}

	if len(s.Prefix)+len(s.Suffix) > 37 {
		fmt.Println("Error! combined prefix and suffix length cannot exceed 37")
		os.Exit(1)
	}
}

// saveToFile appends data to the appropriate output file.
func saveToFile(s *Settings, data string, match bool) error {
	findMatch := s.findMatch()

	if match {
		return appendToFile(s.Filename, data)
	}

	if findMatch && s.SaveAll {
		return appendToFile(s.AllResults, data)
	} else if !findMatch {
		return appendToFile(s.Filename, data)
	}
	return nil
}

func appendToFile(path, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	return err
}

// isMatch reports whether any line of data matches ^z1q{prefix}.*{suffix}$.
func isMatch(s *Settings, data string) bool {
	re := regexp.MustCompile("^z1q" + s.Prefix + ".*" + s.Suffix + "$")
	for _, line := range strings.Split(data, "\n") {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// estimateTotal estimates the number of mnemonics needed to find a match,
// for progress bar purposes.
func estimateTotal(s *Settings) int64 {
	exp := len(s.Prefix) + len(s.Suffix)
	total := math.Pow(32, float64(exp)) / float64(s.AddressesPerMnemonic)
	if total > float64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(total)
}
