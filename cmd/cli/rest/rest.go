package rest

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Get requests http using url and bearer token
func Get(url, token string) {
	// Create a Resty Client
	client := resty.New()
	client.SetAuthToken(token)

	resp, err := client.R().
		EnableTrace().
		Get(url)

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}
