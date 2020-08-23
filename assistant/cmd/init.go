package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func initialize() {

	t := time.Now()
	s := ""
	switch {
	case t.Hour() < 12:
		s = "Good morning!"
	case t.Hour() < 17:
		s = "Good afternoon."
	default:
		s = "Good evening."
	}

	greetings := fmt.Sprintf("%s Baba, how may i help you", s)
	cmd := exec.Command("espeak", greetings)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
