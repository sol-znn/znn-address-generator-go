package main

import "time"

// Settings holds all runtime configuration, mirroring the Dart Settings class.
type Settings struct {
	Threads              int
	NumberOfMnemonics    int
	AddressesPerMnemonic int

	Filename    string
	AllResults  string
	outputIsSet bool

	Prefix string
	Suffix string

	Verbose  bool
	SaveAll  bool
	Infinite bool
}

func newSettings() *Settings {
	timestamp := time.Now().Format("20060102_150405")
	return &Settings{
		Threads:              4,
		NumberOfMnemonics:    100,
		AddressesPerMnemonic: 5,
		Filename:             "./results/results-" + timestamp + ".txt",
		AllResults:           "./results/allresults-" + timestamp + ".txt",
	}
}

func (s *Settings) findMatch() bool {
	return s.Prefix != "" || s.Suffix != ""
}
