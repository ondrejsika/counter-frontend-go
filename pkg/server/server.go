package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//go:embed favicon.ico
var favicon []byte

func Server() {
	apiOrigin := "http://127.0.0.1"
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

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(favicon)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter, backendHostname, _, extraText, _ := api(apiOrigin)
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
	fmt.Println("Listen on 0.0.0.0:3000, see http://127.0.0.1:3000")
	http.ListenAndServe(":3000", nil)
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
		return 0, "", "", "", err
	}
	defer resp.Body.Close()

	var data CounterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", "", "", err
	}

	return data.Counter, data.Hostname, data.Version, data.ExtraText, nil
}
