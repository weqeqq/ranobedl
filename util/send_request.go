package util

import (
	"fmt"
	"net/http"
	"time"
)

func SendRequest(url string) (*http.Response, error) {
	if response, err := http.Get(url); err != nil {
		return nil, err
	} else {
		if response.StatusCode == http.StatusTooManyRequests {
			defer response.Body.Close()

			time.Sleep(time.Second)
			return SendRequest(url)
		}
		if response.StatusCode == http.StatusInternalServerError {
			defer response.Body.Close()

			fmt.Println("Ранобэлиб гандон сука")
			time.Sleep(time.Second)
			return SendRequest(url)
		}
		if response.StatusCode != http.StatusOK {
			defer response.Body.Close()

			return nil, fmt.Errorf("Status code not 200, %s", response.Status)
		}
		return response, nil
	}
}
