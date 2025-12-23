package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func cmdMETAR(icao string) string {
	raw, err := fetchMETAR(icao)
	if err != nil {
		return "somethin ain't right"
	}
	return fmt.Sprintf("`%s`", raw)
}

func cmdTAF(icao string) string {
	raw, err := fetchTAF(icao)
	if err != nil {
		return "somethin ain't right"
	}
	return fmt.Sprintf("```%s```", raw)
}

func cmdWX(icao string) string {
	icao = strings.ToUpper(icao)
	var metar, taf string
	var metarErr, tafErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		metar, metarErr = fetchMETAR(icao)
	}()
	go func() {
		defer wg.Done()
		taf, tafErr = fetchTAF(icao)
	}()
	wg.Wait()

	if metarErr != nil && tafErr != nil {
		return "somethin definitely ain't right"
	}
	if metarErr != nil {
		metar = "[METAR unavailable]"
	}
	if tafErr != nil {
		taf = "[TAF unavailable]"
	}

	reply := fmt.Sprintf("```\n%s\n\n%s```", metar, taf)
	return reply
}

func fetchMETAR(icao string) (string, error) {
	rootURL := "https://aviationweather.gov/api/data"
	url := fmt.Sprintf("%s/metar?ids=%s&format=json", rootURL, icao)

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var data []struct {
		RawOb string `json:"rawOb"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no METAR for %s", icao)
	}

	return data[0].RawOb, nil
}

func fetchTAF(icao string) (string, error) {
	rootURL := "https://aviationweather.gov/api/data"
	url := fmt.Sprintf("%s/taf?ids=%s&format=json", rootURL, icao)

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var data []struct {
		RawTAF string `json:"rawTAF"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no METAR for %s", icao)
	}

	raw := data[0].RawTAF
	raw = strings.ReplaceAll(raw, " FM", "\n  FM")
	raw = strings.ReplaceAll(raw, " PROB", "\n    PROB")
	raw = strings.ReplaceAll(raw, " TEMPO", "\n    TEMPO")
	raw = strings.ReplaceAll(raw, " BECMG", "\n    BECMG")
	return raw, nil
}
