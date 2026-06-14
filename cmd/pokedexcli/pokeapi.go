package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Outer function convert HTTP Response Body into Go Struct, once receive from fetchJSON() inner function
func callAPI[T any](url string, apiConf *ApiConfig) (T, error) {
	var result T

	// A. If cache has this URL, get cached []byte (cacheEntry exists or not timed-out)
	if cachedData, ok := apiConf.Cache.Get(url); ok {
		err := json.Unmarshal(cachedData, &result)
		return result, err
	}

	// B. If no cache for this URL, call API.
	data, err := fetchJSON(url)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("could not unmarshal: %w", err)
	}

	// Caching successful HTTP Response ([]byte) for 5 seconds each entry
	(*apiConf).Cache.Add(url, data)

	return result, nil
}

// Inner function to Fetch the API response's JSON bytes in slice (of type []byte)
func fetchJSON(url string) ([]byte, error) {
	// Create HTTP GET Request. Shorthand //	res, err := http.Get(apiConf.Next)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("could not create request: %w", err)
	}
	// Set Headers
	req.Header.Set("Content-Type", "application/json")

	// Create a Client object, make the HTTP Request and receive the HTTP Response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("network Error: %v", err)
	}
	defer res.Body.Close()

	// Verify successful HTTP GET Request
	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("could not retrieve those location areas. non-OK HTTP status: %s", res.Status)
	}
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("Response failed with status code: %d and\nbody: %v\n", res.StatusCode, res.Body)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("could not read response body: %w", err)
	}
	return data, nil
}
