package actions

import (
	"github.com/eiannone/keyboard"
)

func isKeyPressed(key keyboard.Key) bool {

	_, keyboardKey, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}

	return keyboardKey == key

}
