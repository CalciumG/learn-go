package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func makeRequest(apiURL string, params map[string]string) string {
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        log.Fatalln(err)
    }

    q := req.URL.Query()
    for key, value := range params {
        q.Add(key, value)
    }
    req.URL.RawQuery = q.Encode()

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalln(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }

    return string(body)
}

func main() {
    http.HandleFunc("/api/get-player-rank", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        name := query.Get("name")
        if name == "" {
            http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
            return
        }

        params := map[string]string{
            "name": name,
        }

        response := makeRequest("https://api.the-finals-leaderboard.com/v1/leaderboard/live/crossplay", params)

        fmt.Fprint(w, "Response from API: "+response)
    })

    http.ListenAndServe(":8080", nil)
}