package text

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const ApiKey = "UKHKHD37C43H4ZUU65WBZNMXVJKOO6RO"

type WitAiContact struct {
	Confidence float32 `json:"confidence"`
	Suggested  bool    `json:"suggested"`
	Type       string  `json:"type"`
	Value      string  `json:"value"`
}

type WitAiIntent struct {
	Confidence float32 `json:"confidence"`
	Value      string  `json:"value"`
}

type WitAiEntities struct {
	Contact []WitAiContact    `json:"contact"`
	Intent  []WitAiIntent     `json:"intent"`
	SongTitle  []WitAiContact `json:"song_title"`
}

type WitAiOutcome struct {
	Text     string        `json:"text"`
	Entities WitAiEntities `json:"entities"`
	Intent   string        `json:"intent"`
}

type WitAiResponse struct {
	Text     string         `json:"text"`
	Outcomes []WitAiOutcome `json:"outcomes"`
}


func ConvertAudioToWitAiResponse(speechByte *bytes.Buffer) *WitAiResponse {

	var response *WitAiResponse

	stringResponse := sendWitBuff(speechByte)

	fmt.Println(stringResponse)

	rawResponseByte := []byte(stringResponse)
	err := json.Unmarshal(rawResponseByte, &response)

	if err != nil {
		log.Println(err)
	}

	return response
}

func sendWitBuff(buffer *bytes.Buffer) string {
	url := "https://api.wit.ai/speech?v=20141022"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, buffer)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+ApiKey)
	req.Header.Set("Content-Type", "audio/raw;encoding=signed-integer;bits=16;rate=20k;endian=little")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
