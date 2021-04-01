package core

import (
	"testing"
)

func TestLinkFixNoPrefixSmall(t *testing.T) {
	link := LinkFix("a.aa")
	if link != "https://a.aa" {
		t.Fatal("Does not have https:// added in")
	}
}

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
