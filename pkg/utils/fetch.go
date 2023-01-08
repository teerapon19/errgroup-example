package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

func FetchJsonWithContext(ctx context.Context, dest any, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dest)
	if err != nil {
		return err
	}

	return nil
}

func FetchJson(dest any, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dest)
	if err != nil {
		return err
	}

	return nil
}
