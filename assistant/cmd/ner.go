package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"gopkg.in/jdkato/prose.v2"
)

// ProdigyOutput represents a single entry of Prodigy's JSON Lines output.
//
// `LabeledEntity` is a structure defined by prose that specifies where the
// entities are within the given `Text`.
type ProdigyOutput struct {
	Text   string
	Spans  []prose.LabeledEntity
	Answer string
}

// ReadProdigy reads our JSON Lines file line-by-line, populating a
// slice of `ProdigyOutput` structures.
func ReadProdigy(jsonLines []byte) []ProdigyOutput {
	dec := json.NewDecoder(bytes.NewReader(jsonLines))
	entries := []ProdigyOutput{}
	for {
		ent := ProdigyOutput{}
		err := dec.Decode(&ent)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		entries = append(entries, ent)
	}
	return entries
}

// Split divides our human-annotated data set into two groups: one for training
// our model and one for testing it.
//
// We're using an 80-20 split here, although you may want to use a different
// split.
func Split(data []ProdigyOutput) ([]prose.EntityContext, []ProdigyOutput) {
	cutoff := int(float64(len(data)) * 0.8)

	train, test := []prose.EntityContext{}, []ProdigyOutput{}
	for i, entry := range data {
		if i < cutoff {
			train = append(train, prose.EntityContext{
				Text:   entry.Text,
				Spans:  entry.Spans,
				Accept: entry.Answer == "accept"})
		} else {
			test = append(test, entry)
		}
	}

	return train, test
}

func train() {
	data, err := ioutil.ReadFile("../reddit-product.jsonl")
	if err != nil {
		panic(err)
	}
	train, test := Split(ReadProdigy(data))

	// Here, we're training a new model named PRODUCT with the training portion
	// of our annotated data.
	//
	// Depending on your hardware, this should take around 1 - 3 minutes.
	model := prose.ModelFromData("PRODUCT", prose.UsingEntities(train))

	// Now, let's test our model:
	correct := 0.0
	for _, entry := range test {
		// Create a document without segmentation, which isn't required for NER.
		doc, err := prose.NewDocument(
			entry.Text,
			prose.WithSegmentation(false),
			prose.UsingModel(model))

		if err != nil {
			panic(err)
		}
		ents := doc.Entities()

		if entry.Answer != "accept" && len(ents) == 0 {
			// If we rejected this entity during annotation, prose shouldn't
			// have labeled it.
			correct++
		} else {
			// Otherwise, we need to verify that we found the correct entities.
			expected := []string{}
			for _, span := range entry.Spans {
				expected = append(expected, entry.Text[span.Start:span.End])
			}
			if reflect.DeepEqual(expected, ents) {
				correct++
			}
		}
	}
	fmt.Printf("Correct (%%): %f\n", correct/float64(len(test)))
	model.Write("PRODUCT") // Save the model to disk.
}

func recognize(data string) []prose.Entity {

	// model := prose.ModelFromDisk("PRODUCT")

	doc, err := prose.NewDocument(data,
		prose.WithSegmentation(false))

	if err != nil {
		panic(err)
	}

	// Iterate over the doc's named-entities:
	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
		// Windows 10 PRODUCT
	}

	return doc.Entities()
}
