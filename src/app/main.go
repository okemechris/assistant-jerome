package main

import (
	"assistant-jerome/src/app/text"
	"assistant-jerome/src/app/voice"
)
func main() {
	// initialize()
	// initAudio()
	text.InitSpeechToText()
	voice.CommandHandler()
	// go forever()
	// select {}
}

// func forever() {
// 	for {

// 		if isKeyPressed(keyboard.KeyEsc) {
// 			fmt.Println("esc pressed")
// 		}

// 		if !handling {
// 			CommandHandler()
// 		}

// 		time.Sleep(time.Second)
// 	}

// }
