package actions

import (
	"assistant-jerome/src/app/text"
	"fmt"
	"log"
	"os/exec"
)


func Greet() {
	text.SpeakText("Hi, how can I help?")
}

func CommandUnknown() {
	text.SpeakText("Sorry, I don't understand. Please take that again")
}

func PlayMusic(songTitle string, artist string ) {
	text.SpeakText(fmt.Sprintf("Playing %s by %s",songTitle,artist))
	ExecCommand("rhythmbox-client","--play")
}

func PlayMusics() {
	//text.SpeakText(fmt.Sprintf("Playing %s by %s",songTitle,artist))
	ExecCommand("rhythmbox-client","--play")
}




func ExecCommand(name string, arg ...string){
	cmd := exec.Command(name, arg...)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}