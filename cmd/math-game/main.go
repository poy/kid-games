package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/gorilla/mux"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/say-problem/problem.ogg", sayProblemHandler)
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

func parseValue(name string, values url.Values) (int, error) {
	v := values.Get(name)
	if v == "" {
		return 0, fmt.Errorf("You don't have %q value", name)
	}

	x, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q: %v", name, err)
	}

	return int(x), nil
}

func sayProblemHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	a, err := parseValue("a", values)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := parseValue("b", values)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

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
				Text: fmt.Sprintf("Ryder, what is %d plus %d equal", a, b),
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
		var ih indexHTML
		ih.NumberA = rand.Intn(10)
		ih.NumberB = rand.Intn(10)

		if err := indexHTMLTemplate.Execute(w, ih); err != nil {
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
	NumberA, NumberB int
}

var indexHTMLTemplate = template.Must(template.New("index").Parse(`
<!doctype html>
<html>
<head>
  <title>Math Game</title>
</head>
<body>
    <div>
    <div>
        <audio controls>
            <source src="/say-problem/problem.ogg?a={{.NumberA}}&b={{.NumberB}}" type="audio/ogg">
            Your browser does not support the audio element.
        </audio>
    </div>
    <div>
        <input type="hidden" id="numberA" value="{{.NumberA}}"><input type="hidden" id="numberB" value="{{.NumberB}}"><input type="text" id="wordInput"><button id="enterButton" onClick="window.location.href=window.location.href">Enter Word</button>
    </div>
    </div>
<script src="main.js"></script>
</body>
</html>
`))
