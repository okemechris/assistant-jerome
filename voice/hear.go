package voice

import (
	"assistant-jerome/actions"
	"assistant-jerome/text"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"log"
	"time"
)

var audioRunning bool
var handling bool

func InitAudio() {
	if !audioRunning {
		portaudio.Initialize()
		audioRunning = true
	}
}

func FreeAudio() {
	if audioRunning {
		portaudio.Terminate()
	}
}

const DefaultQuietTime = time.Second

type State int

const (
	Waiting State = iota
	Listening
)

type ListenOpts struct {
	State            func(State)
	QuietDuration    time.Duration
	AlreadyListening bool
}

func ListenIntoBuffer(opts ListenOpts) (*bytes.Buffer, error) {
	in := make([]int16, 8196)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return nil, err
	}

	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return nil, err
	}

	var (
		buf            bytes.Buffer
		heardSomething = opts.AlreadyListening
		quiet          bool
		quietTime      = opts.QuietDuration
		quietStart     time.Time
		lastFlux       float64
	)

	vad := NewVAD(len(in))

	if quietTime == 0 {
		quietTime = DefaultQuietTime
	}

	if opts.State != nil {
		if heardSomething {
			opts.State(Listening)
		} else {
			opts.State(Waiting)
		}
	}

reader:
	for {
		err = stream.Read()
		if err != nil {
			return nil, err
		}

		err = binary.Write(&buf, binary.LittleEndian, in)
		if err != nil {
			return nil, err
		}

		flux := vad.Flux(in)

		if lastFlux == 0 {
			lastFlux = flux
			continue
		}

		if heardSomething {
			if flux*1.75 <= lastFlux {
				if !quiet {
					quietStart = time.Now()
				} else {
					diff := time.Since(quietStart)

					if diff > quietTime {
						break reader
					}
				}

				quiet = true
			} else {
				quiet = false
				lastFlux = flux
			}
		} else {
			if flux >= lastFlux*1.75 {
				heardSomething = true
				if opts.State != nil {
					opts.State(Listening)
				}
			}

			lastFlux = flux
		}
	}

	err = stream.Stop()
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func Listen() {

	InitAudio()
	defer FreeAudio()

	opts := ListenOpts{
		QuietDuration:    1 * time.Second,
		AlreadyListening: true,
	}

	fmt.Printf("speak now\n")

	buf, err := ListenIntoBuffer(opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("recognizing...\n")
	handling = true

	convertedResponse := text.ConvertAudioToWitAiResponse(buf)

	determineAction(convertedResponse)

	Listen()
}

func determineAction(witAiResponse *text.WitAiResponse) {

	if len(witAiResponse.Outcomes) < 1 {
		return
	}

	if len(witAiResponse.Outcomes[0].Entities.Intent) < 1 {
		return
	}

	switch witAiResponse.Outcomes[0].Entities.Intent[0].Value {

	case "greetings":
		actions.Greet()
	case "play_music":
		//actions.PlayMusic(witAiResponse.Outcomes[0].Entities.SongTitle[0].Value,witAiResponse.Outcomes[0].Entities.Contact[0].Value)
		actions.PlayMusic("","")
	default:
		actions.CommandUnknown()
	}
}
