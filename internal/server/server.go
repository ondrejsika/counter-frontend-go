package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ondrejsika/counter-frontend-go/version"
)

//go:embed favicon.ico
var favicon []byte

func Server() {
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	failOnError := false
	if os.Getenv("FAIL_ON_ERROR") == "1" {
		failOnError = true
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

	hostname, _ := os.Hostname()

	err := checkApiStatus(apiOrigin)
	if err != nil {
		if failOnError {
			log.Fatalln("Quitting due to error and FAIL_ON_ERROR=1:", err)
		}
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(favicon)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter, backendHostname, _, extraText, err := api(apiOrigin)
		if err != nil {
			if failOnError {
				log.Fatalln("Quitting due to error and FAIL_ON_ERROR=1:", err)
			}
		}
		counterStr := fmt.Sprintf("%d", counter)
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
				<h2>`+hostname+`</h2>
				<h2>`+backendHostname+`</h2>
			</div>
		</section>
		</body></html>
		`)
	})
	http.HandleFunc("/api/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"live": true}`)
	})
	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"version": "`+version.Version+`"}`)
	})
	fmt.Println("Listen on 0.0.0.0:" + port + ", see http://127.0.0.1:" + port)
	http.ListenAndServe(":"+port, nil)
}

func api(origin string) (int, string, string, string, error) {
	type CounterResponse struct {
		Counter   int    `json:"counter"`
		Hostname  string `json:"hostname"`
		Version   string `json:"version"`
		ExtraText string `json:"extra_text"`
	}

	resp, err := http.Get(origin + "/api/counter")
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
