package main

import (
	"context"
	"log"
	"os"

	language "cloud.google.com/go/language/apiv1"
	"github.com/golang/protobuf/proto"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func ner(text string) {

	ctx := context.Background()

	client, err := language.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	printResp(analyzeEntities(ctx, client, text))
}

func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}
