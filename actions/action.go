package actions

import (
	"fmt"
)


func Greet() {
	SpeakText("Hi, how can I help?")
}

func CommandUnknown() {
	SpeakText("Sorry, I don't understand. Please take that again")
}

func playMusic(songTitle string, artist string ) {
	if artist == ""{
		SpeakText(fmt.Sprintf("Playing %s",songTitle))
		return
	}
	SpeakText(fmt.Sprintf("Playing %s by %s",songTitle,artist))
}

func PlayMusic(songTitle string, artist string ) {

	if songTitle == ""{
		SpeakText("Playing requested song")
		return
	}

	playMusic(songTitle, artist)


}