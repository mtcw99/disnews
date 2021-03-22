package core

import (
	"testing"
)

func TestLinkFixNoPrefix(t *testing.T) {
	link := LinkFix("github.com")
	if link != "https://github.com" {
		t.Fatal("Does not have https:// added in")
	}
}

func TestLinkFixWithPrefix(t *testing.T) {
	link := LinkFix("https://github.com")
	if link != "https://github.com" {
		t.Fatal("Does have https:// added in")
	}
}
