package ui

import (
	"os"

	"golang.org/x/term"
)

func getTerminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24 //fallback to default
	}
	return w, h
}
