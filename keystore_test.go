package main

import (
	"crypto/ed25519"
	"testing"

	"github.com/tyler-smith/go-bip39"
)

// TestKnownVector checks derivation against addresses produced by the
// original Dart implementation (znn_sdk_dart) for a fixed mnemonic, to
// confirm the Go port is bit-for-bit compatible.
func TestKnownVector(t *testing.T) {
	mnemonic := "neither fiction member destroy become hour artefact aspect brisk hello subway worry sell side space secret desk draw sting shoe scheme vital hedgehog track"
	want := []string{
		"z1qq9t4c0hk80xwwunvwncnh7nz5yjw4l0mder7f",
		"z1qqh55r4f85z9e3rxalhmsc7aczaqv4e6y265r5",
		"z1qrr9lpsg2knz45qgu6gngt9jhrfy3echylss8v",
		"z1qp4cc9u4dsfxgjehdnlu22hn0pkshu505juz64",
		"z1qzsgh49gmlpu5c7ajlrg77sacdv3uts2kkly6z",
	}

	seed := bip39.NewSeed(mnemonic, "")
	for i, wantAddr := range want {
		accountSeed := deriveAccountKey(seed, i)
		priv := ed25519.NewKeyFromSeed(accountSeed[:])
		pub := priv.Public().(ed25519.PublicKey)
		got := addressFromPublicKey(pub)
		if got != wantAddr {
			t.Errorf("address %d: got %s, want %s", i, got, wantAddr)
		}
	}
}
