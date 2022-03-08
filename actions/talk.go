package actions

import (
	"log"
	"os/exec"
)

func SpeakText(textToSpeak string) {
	cmd := exec.Command("espeak", textToSpeak)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
