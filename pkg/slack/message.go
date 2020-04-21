package slack

import (
	"bytes"
	"log"
	"net/http"
)

type Data struct {
	Text string `json:"text"`
}

func PostMessage(d Data, endpoint string) error {
	jsonStr := `{"text":"` + d.Text + `"}`

	log.Printf("[INFO] Json Payload: %#v", jsonStr)
	log.Printf("[INFO] endpoint: %s", endpoint)
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
