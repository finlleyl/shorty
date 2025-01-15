package main

import (
	"bufio"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"
	data := url.Values{}

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)

	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	long = strings.TrimSuffix(long, "\n")
	data.Set("url", long)

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(data.Encode()).
		Post(endpoint)

	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status Code: ", response.StatusCode())
	fmt.Println("Response Body: ", response.String())

}
