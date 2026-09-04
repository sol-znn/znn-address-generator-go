package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var firstChars = map[byte]bool{'p': true, 'q': true, 'r': true, 'z': true}

type flagSpec struct {
	names  []string // e.g. []string{"t", "threads"}
	isBool bool
	help   string
}

var flagSpecs = []flagSpec{
	{[]string{"t", "threads"}, false, "number of cpu threads (default 4)"},
	{[]string{"n", "mnemonics"}, false, "number of mnemonics to generate (default 100)"},
	{[]string{"a", "addresses"}, false, "number of addresses per mnemonic (default 5)"},
	{[]string{"o", "output"}, false, "output file name (default ./results/results-<time>.txt)"},
	{[]string{"p", "prefix"}, false, "save addresses that match a specific prefix"},
	{[]string{"s", "suffix"}, false, "save addresses that match a specific suffix"},
	{[]string{"i", "infinite"}, true, "generate infinite mnemonics"},
	{[]string{"all"}, true, "save all generated mnemonics, even if they don't match"},
	{[]string{"v", "verbose"}, true, "display mnemonics as they're being generated"},
	{[]string{"h", "help"}, true, "Displays help information"},
}

// specByName maps every flag name (short and long) to its spec.
func specByName() map[string]flagSpec {
	m := map[string]flagSpec{}
	for _, spec := range flagSpecs {
		for _, name := range spec.names {
			m[name] = spec
		}
	}
	return m
}

// canonicalName returns the long name for a flag spec (or its only name).
func (f flagSpec) canonicalName() string {
	return f.names[len(f.names)-1]
}

func parseArgs(args []string) *Settings {
	s := newSettings()

	if len(args) == 0 {
		menu(usageText())
		os.Exit(0)
	}

	byName := specByName()
	values := map[string]string{}
	seen := map[string]bool{}

	for i := 0; i < len(args); i++ {
		tok := args[i]

		if !strings.HasPrefix(tok, "-") {
			fmt.Printf("Error! Unrecognized argument: %s\n", tok)
			fmt.Println("(Did you forget a leading '-' or '--' before a flag name?)")
			os.Exit(1)
		}

		name := strings.TrimLeft(tok, "-")
		var inlineValue string
		hasInline := false
		if idx := strings.Index(name, "="); idx != -1 {
			inlineValue = name[idx+1:]
			name = name[:idx]
			hasInline = true
		}

		spec, ok := byName[name]
		if !ok {
			fmt.Printf("Error! Unrecognized flag: %s\n", tok)
			os.Exit(1)
		}

		canonical := spec.canonicalName()
		seen[canonical] = true

		if spec.isBool {
			if hasInline {
				values[canonical] = inlineValue
			} else {
				values[canonical] = "true"
			}
			continue
		}

		if hasInline {
			values[canonical] = inlineValue
			continue
		}

		if i+1 >= len(args) {
			fmt.Printf("Error! Flag %s requires a value\n", tok)
			os.Exit(1)
		}
		i++
		values[canonical] = args[i]
	}

	if seen["help"] {
		menu(usageText())
		os.Exit(0)
	}

	if seen["threads"] {
		n, err := strconv.Atoi(values["threads"])
		if err != nil {
			fmt.Println("Error! --threads must be a number")
			os.Exit(1)
		}
		s.Threads = n
		if s.Threads < 1 {
			fmt.Println("Error! You must use at least one thread")
			os.Exit(1)
		}
	}

	if seen["mnemonics"] {
		n, err := strconv.Atoi(values["mnemonics"])
		if err != nil {
			fmt.Println("Error! --mnemonics must be a number")
			os.Exit(1)
		}
		s.NumberOfMnemonics = n
		if s.NumberOfMnemonics < 1 {
			fmt.Println("Error! You must generate at least one mnemonic")
			os.Exit(1)
		}
	}

	if seen["addresses"] {
		n, err := strconv.Atoi(values["addresses"])
		if err != nil {
			fmt.Println("Error! --addresses must be a number")
			os.Exit(1)
		}
		s.AddressesPerMnemonic = n
		if s.AddressesPerMnemonic < 1 {
			fmt.Println("Error! You must generate at least one address per mnemonic")
			os.Exit(1)
		}
	}

	if seen["prefix"] {
		s.Prefix = toLower(values["prefix"])
		isValid(s, s.Prefix)
		if len(s.Prefix) == 0 || !firstChars[s.Prefix[0]] {
			fmt.Println("Error! The prefix should start with one of these letters: p, q, r, z")
			os.Exit(1)
		}
	}

	if seen["suffix"] {
		s.Suffix = toLower(values["suffix"])
		isValid(s, s.Suffix)
	}

	if seen["output"] {
		s.Filename = values["output"]
		s.outputIsSet = true
	}

	// The all-results file always lives under ./results, and a custom
	// -o/--output path may point at a directory that doesn't exist yet,
	// so make sure both output directories are present.
	for _, path := range []string{s.Filename, s.AllResults} {
		if dir := filepath.Dir(path); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error! could not create %s: %v\n", dir, err)
				os.Exit(1)
			}
		}
	}

	if seen["verbose"] && values["verbose"] == "true" {
		s.Verbose = true
	}
	if seen["all"] && values["all"] == "true" {
		s.SaveAll = true
	}
	if seen["infinite"] && values["infinite"] == "true" {
		s.Infinite = true
	}

	return s
}

func usageText() string {
	var b strings.Builder
	for _, spec := range flagSpecs {
		var names []string
		for _, n := range spec.names {
			if len(n) == 1 {
				names = append(names, "-"+n)
			} else {
				names = append(names, "--"+n)
			}
		}
		fmt.Fprintf(&b, "  %-20s %s\n", strings.Join(names, ", "), spec.help)
	}
	return b.String()
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
