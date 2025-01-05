package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/JohannesKaufmann/html-to-markdown/v2"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

type GraphQLResponse struct {
	Data struct {
		Question struct {
			QuestionTitle string `json:"questionTitle"`
			Content       string `json:"content"` // HTML format
			Difficulty    string `json:"difficulty"`
		} `json:"question"`
	} `json:"data"`
}

func main() {
	problemSlug := "substrings-of-size-three-with-distinct-characters" // Replace with desired problem slug
	url := "https://leetcode.com/graphql"

	query := fmt.Sprintf(`
	{
		question(titleSlug: "%s") {
			questionTitle
			content
			difficulty
		}
	}`, problemSlug)

	// Create the GraphQL request
	requestBody, _ := json.Marshal(GraphQLRequest{Query: query})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, _ := io.ReadAll(resp.Body)
	var response GraphQLResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// Convert HTML content to Markdown or Plain Text
	htmlContent := response.Data.Question.Content
	// converter := converter.NewConverter("", true, nil) // HTML-to-Markdown converter
	markdownContent, err := htmltomarkdown.ConvertString(htmlContent)
	if err != nil {
		panic(err)
	}

	// Output the problem details
	question := response.Data.Question
	fmt.Println("Title:", question.QuestionTitle)
	fmt.Println("Difficulty:", question.Difficulty)
	fmt.Println("\nDescription:\n", markdownContent)
}
