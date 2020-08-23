package main

import (
	"context"

	gcontext "golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/grpc"
)

type GCPSpeechConv struct {
	ctx  gcontext.Context
	conn *grpc.ClientConn

	client speech.SpeechClient
}

func NewGCPSpeechConv(accountFile string) (*GCPSpeechConv, error) {
	ctx := context.Background()
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
		option.WithServiceAccountFile(accountFile),
	)

	if err != nil {
		return nil, err
	}

	client := speech.NewSpeechClient(conn)

	return &GCPSpeechConv{ctx, conn, client}, nil
}

func (gcp *GCPSpeechConv) Convert(data []byte) (string, error) {
	resp, err := gcp.recognize(data)
	if err != nil {
		return "", err
	}

	var best *speech.SpeechRecognitionAlternative

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			if best == nil || alt.Confidence > best.Confidence {
				best = alt
			}
		}
	}

	if best == nil {
		return "", nil
	}

	return best.Transcript, nil
}

func (gcp *GCPSpeechConv) recognize(data []byte) (*speech.RecognizeResponse, error) {
	return gcp.client.Recognize(gcp.ctx, &speech.RecognizeRequest{
		Config: &speech.RecognitionConfig{
			Encoding:        speech.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speech.RecognitionAudio{
			AudioSource: &speech.RecognitionAudio_Content{Content: data},
		},
	})
}
