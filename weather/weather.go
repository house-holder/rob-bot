package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ATISResp struct {
	Airport string `json:"airport"`
	Code    string `json:"code"`
	ATIS    string `json:"datis"`
	Time    string `json:"time"`
}

func CmdMETAR(icao string) string {
	raw, err := fetchMETAR(icao)
	if err != nil {
		return "somethin ain't right"
	}
	return fmt.Sprintf("`%s`", raw)
}

func CmdTAF(icao string) string {
	raw, err := fetchTAF(icao)
	if err != nil {
		return "somethin ain't right"
	}
	return fmt.Sprintf("```%s```", raw)
}

func CmdWX(icao string) string {
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

func CmdATIS(icao string) (string, string, error) {
	atis, code, timeStr, err := fetchATIS(icao)
	caps := strings.ToUpper(icao)
	if err != nil {
		return "", "", fmt.Errorf("no ATIS for %s", caps)
	}

	atisTime, err := parseATISTime(timeStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse time: %w", err)
	}

	ageMinutes := int(time.Since(atisTime).Minutes())
	ageText := fmt.Sprintf("%d minute", ageMinutes)
	if ageMinutes != 1 {
		ageText += "s"
	}

	message := fmt.Sprintf("**(%s old)**\n>>> %s", ageText, atis)
	return message, code, nil
}

func parseATISTime(timeStr string) (time.Time, error) {
	if len(timeStr) != 4 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hours, err := strconv.Atoi(timeStr[0:2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hours: %w", err)
	}
	minutes, err := strconv.Atoi(timeStr[2:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minutes: %w", err)
	}

	now := time.Now().UTC()
	atisTime := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, time.UTC)

	if atisTime.After(now) {
		atisTime = atisTime.AddDate(0, 0, -1)
	}

	return atisTime, nil
}

func fetchATIS(icao string) (string, string, string, error) {
	station := strings.ToLower(icao)
	url := fmt.Sprintf("https://atis.info/api/%s", station)

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch ATIS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("API returns status %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("IO read failed: %w", err)
	}

	jsonObjs := []ATISResp{}
	err = json.Unmarshal(bytes, &jsonObjs)
	if err != nil {
		return "", "", "", fmt.Errorf("JSON parse failed: %w", err)
	}

	if len(jsonObjs) == 0 {
		return "", "", "", fmt.Errorf("no ATIS for %s", strings.ToUpper(icao))
	}
	tgt := jsonObjs[0]

	return tgt.ATIS, tgt.Code, tgt.Time, nil
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
