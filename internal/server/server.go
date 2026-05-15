package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ondrejsika/counter-frontend-go/version"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

//go:embed favicon.ico
var favicon []byte

func Server(versionOverride string) {
	if versionOverride != "" {
		version.Version = versionOverride
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	failOnError := false
	if os.Getenv("FAIL_ON_ERROR") == "1" {
		failOnError = true
	}

	readOnly := false
	if os.Getenv("READ_ONLY") == "1" {
		readOnly = true
	}

	apiOrigin := "http://127.0.0.1:8000"
	if os.Getenv("API_ORIGIN") != "" {
		apiOrigin = os.Getenv("API_ORIGIN")
	}

	fontColor := "#000000"
	if os.Getenv("FONT_COLOR") != "" {
		fontColor = os.Getenv("FONT_COLOR")
	}

	backgroundColor := "#ffffff"
	if os.Getenv("BACKGROUND_COLOR") != "" {
		backgroundColor = os.Getenv("BACKGROUND_COLOR")
	}

	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	hostname, _ := os.Hostname()

	err := checkApiStatus(apiOrigin)
	if err != nil {
		if failOnError {
			Logger.Fatal().Str("hostname", hostname).Msg("Quitting due to error and FAIL_ON_ERROR=1: " + err.Error())
		}
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(favicon)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter, backendHostname, backendVersion, extraText, err := api(apiOrigin, readOnly)
		if err != nil {
			if failOnError {
				Logger.Fatal().Str("hostname", hostname).Msg("Quitting due to error and FAIL_ON_ERROR=1: " + err.Error())
			}
		}
		Logger.Info().
			Str("hostname", hostname).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("counter", counter).
			Msg(r.Method + " " + r.URL.Path)
		counterStr := fmt.Sprintf("%d", counter)

		// Check if User-Agent header exists
		if userAgentList, ok := r.Header["User-Agent"]; ok {
			// Check if User-Agent header has some data
			if len(userAgentList) > 0 {
				// If User-Agent starts with curl, use plain text
				if strings.HasPrefix(userAgentList[0], "curl") {
					indexPlainText(w, hostname, backendHostname, backendVersion, extraText, counterStr)
				} else {
					// If User-Agent header presents and not starts with curl
					// use HTML (Chrome, Safari, Firefox, ...)
					indexHtml(
						w, hostname, backendHostname, backendVersion, extraText,
						fontColor, backgroundColor, counterStr)
				}
			}
		} else {
			// If User-Agent header doesn't exists, use plain text
			indexPlainText(
				w, hostname, backendHostname, backendVersion, extraText, counterStr)
		}
	})
	http.HandleFunc("/api/livez", func(w http.ResponseWriter, r *http.Request) {
		Logger.Info().
			Str("hostname", hostname).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg(r.Method + " " + r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"live": true}`)
	})
	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		Logger.Info().
			Str("hostname", hostname).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg(r.Method + " " + r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"version": "`+version.Version+`"}`)
	})
	Logger.Info().Str("hostname", hostname).Msg("Listen on 0.0.0.0:" + port + ", see http://127.0.0.1:" + port)
	http.ListenAndServe(":"+port, nil)
}

func api(origin string, readOnly bool) (int, string, string, string, error) {
	type CounterResponse struct {
		Counter   int    `json:"counter"`
		Hostname  string `json:"hostname"`
		Version   string `json:"version"`
		ExtraText string `json:"extra_text"`
	}

	apiPath := origin + "/api/counter"
	if readOnly {
		apiPath = origin + "/api/read-counter"
	}

	resp, err := http.Get(apiPath)
	if err != nil {
		return -1, "", "", "", err
	}
	defer resp.Body.Close()

	var data CounterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return -1, "", "", "", err
	}

	return data.Counter, data.Hostname, data.Version, data.ExtraText, nil
}

func checkApiStatus(origin string) error {
	resp, err := http.Get(origin + "/api/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API status is not OK: %d", resp.StatusCode)
	}

	return nil
}

func indexHtml(
	w http.ResponseWriter, hostname, backendHostname, backendVersion,
	extraText, fontColor, backgroundColor, counterStr string,
) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html>
		<html lang="en"><head>
		<meta charset="UTF-8">
		<title>`+extraText+`</title>
		<style>
		html, body {
			height: 100%;
			color: `+fontColor+`;
			background-color: `+backgroundColor+`
		}
		.center-parent {
			width: 100%;
			height: 100%;
			display: table;
			text-align: center;
		}
		.center-parent > .center-child {
			display: table-cell;
			vertical-align: middle;
		}
		</style>
		<style>
		h1 {
			font-family: Arial;
			font-size: 5em;
		}
		h2 {
			font-family: Arial;
			font-size: 2em;
		}
		</style>
		<link rel="icon" href="/favicon.ico">
		</head>
		<body>
		<section class="center-parent">
			<div class="center-child">
				<h1>👋</h1>
				<h1>`+extraText+`</h1>
				<h1>`+counterStr+`</h1>
				<h2>`+hostname+` `+version.Version+`</h2>
				<h2>`+backendHostname+` `+backendVersion+`</h2>
			</div>
		</section>
		</body></html>
		`)
}

func indexPlainText(
	w http.ResponseWriter, hostname, backendHostname, backendVersion,
	extraText, counterStr string,
) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, `👋 `+extraText+` `+counterStr+` `+hostname+` `+version.Version+
		` `+backendHostname+` `+backendVersion)
}
