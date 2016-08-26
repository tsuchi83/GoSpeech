package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	speech "google.golang.org/api/speech/v1beta1"
)

const usage = `Usage: gospeech <service account keyfile>
Service account keyfile is JSON File.
Check site for details. https://cloud.google.com/speech/docs/common/auth
`

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read service account keyfile: %v", err)
	}
	config, err := google.JWTConfigFromJSON(b, speech.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Unable to read service account keyfile: %v", err)
	}
	client := config.Client(ctx)

	srv, err := speech.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// for Syncrecognize
	res, err := srv.Speech.Syncrecognize(&speech.SyncRecognizeRequest{
		Config: &speech.RecognitionConfig{
			Encoding:     "FLAC",
			SampleRate:   16000,
			LanguageCode: "ja-JP",
		},
		Audio: &speech.RecognitionAudio{
			Uri: "gs://<your bucket name>/demo/sample.flac",
		},
	}).Do()
	if err != nil {
		log.Fatal(err)
	}

	// Print the results.
	for _, result := range res.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}
