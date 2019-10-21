package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/gorilla/mux"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var spellingWords = []string{
	"king",
	"range",
	"flung",
	"bring",
	"fang",
	"long",
	"hung",
	"strong",
	"saw",
	"but",
	"math",
	"this",
	"chip",
	"shop",
	"shed",
	"brush",
	"when",
	"flash",
	"only",
	"shoe",
	"duck",
	"fudge",
	"track",
	"bridge",
	"sick",
	"pledge",
	"clock",
	"dodge",
	"new",
	"whose",
	"grabbed",
	"flipped",
	"popped",
	"hummed",
	"grinned",
	"budded",
	"begged",
	"dotted",
	"give",
	"live",
	"jumped",
	"stamped",
	"rested",
	"filled",
	"yelled",
	"pinched",
	"planted",
	"ended",
	"gone",
	"little",
	"catch",
	"fetch",
	"itch",
	"scratch",
	"stretch",
	"switch",
	"clutch",
	"blotch",
	"any",
	"many",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/say-word/{word}.ogg", sayWordHandler)
	r.HandleFunc("/{file}", serveContent)
	r.HandleFunc("/", serveContent)

	log.Fatal(http.ListenAndServe(getPort(), r))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8080"
	}
	return fmt.Sprintf(":" + p)
}

func sayWordHandler(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]
	ctx := r.Context()

	c, err := texttospeech.NewClient(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Ask Dad: failed to create client: %s", err)
		fmt.Fprintf(w, "Ask Dad: failed to create client: %s", err)
		return
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
			Pitch:         0,
			SpeakingRate:  0.6,
		},
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: word,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Wavenet-D",
		},
	}
	resp, err := c.SynthesizeSpeech(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Ask Dad: failed to make request: %s", err)
		fmt.Fprintf(w, "Ask Dad: failed to make request: %s", err)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(resp.AudioContent)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Ask Dad: failed to make write data: %s", err)
		fmt.Fprintf(w, "Ask Dad: failed to make write data: %s", err)
		return
	}
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["file"]
	if filename == "" || filename == "index.html" || filename == "index.htm" {
		word := spellingWords[rand.Intn(len(spellingWords))]

		if err := indexHTMLTemplate.Execute(w, indexHTML{Word: word}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed: %s", err)
			return
		}
		return
	}

	f, err := os.Open(filepath.Join("content", filename))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed: %s", err)
		return
	}

	if _, err := io.Copy(w, f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed: %s", err)
		return
	}
}

type indexHTML struct {
	Word string
}

var indexHTMLTemplate = template.Must(template.New("index").Parse(`
<!doctype html>
<html>
<head>
  <title>GopherJS DOM example - Create and Append Element</title>
</head>
<body>
    <div>
    <div>
        <audio controls>
            <source src="/say-word/{{.Word}}.ogg" type="audio/ogg">
            Your browser does not support the audio element.
        </audio>
    </div>
    <div>
        <input type="hidden" id="answerInput" value="{{.Word}}"><input type="text" id="wordInput"><button id="enterButton" onClick="window.location.href=window.location.href">Enter Word</button>
    </div>
    </div>
<script src="main.js"></script>
</body>
</html>
`))
