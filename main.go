package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

// APIResponse maps to the JSON response from the Open Weather Map API
type APIResponse struct {
	Summary []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Weather struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
		Pressure int     `json:"pressure"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Location string `json:"name"`
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	http.HandleFunc("/receive", slashCommandHandler)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(port, nil)
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/spin":
		params := &slack.Msg{Text: s.Text}
		opts := strings.Split(params.Text, ",")
		n := rand.Int() % len(opts)
		winner := opts[n]

		response := fmt.Sprintf("Hold on tight, we're spinning the wheel!\n\n And the winner is\n *%v*", winner)
		w.Write([]byte(response))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
