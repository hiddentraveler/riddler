package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

type CodeSnippet struct {
	Code string `json:"code"` // Language (e.g., "golang", "c")
	Lang string `json:"lang"` // Starter code
}

type GraphQLResponse struct {
	Data struct {
		Question struct {
			QuestionFrontendID  string        `json:"questionFrontendId"`
			QuestionTitle       string        `json:"questionTitle"`
			Content             string        `json:"content"` // HTML format
			Difficulty          string        `json:"difficulty"`
			ExampleTestcaseList []string      `json:"exampleTestcaseList"` // Raw sample input
			CodeSnippets        []CodeSnippet `json:"codeSnippets"`
		} `json:"question"`
	} `json:"data"`
}

func RemoveInvisibleCharacters(input string) string {
	input = strings.ReplaceAll(input, "\u200B", "") // Zero-width space
	return input
}

func main() {
	author := "Neo Orez"
	if len(os.Args) < 2 {
		fmt.Println("Please provide input.")
		return
	}

	// os.Args[0] is the program name; os.Args[1:] are the arguments
	var questionLink string
	for _, arg := range os.Args[1:] {
		questionLink = arg
	}

	var selectedLangOpt string
	var selectedLang string
	fmt.Println("Which language do you want to use?[Go:1][C++:2]")
	fmt.Scanln(&selectedLangOpt)

	var fileExtension string
	var langDir string
	switch selectedLangOpt {
	case "1":
		fileExtension = "go"
		langDir = "golang"
		selectedLang = "Go"
	case "2":
		fileExtension = "cpp"
		langDir = "cpp"
		selectedLang = "C++"

	}

	problemSlug := strings.Split(questionLink, "/")[4]

	url := "https://leetcode.com/graphql"

	query := fmt.Sprintf(`
	{
		question(titleSlug: "%s") {
			questionFrontendId
			questionTitle
			content
			difficulty
			exampleTestcaseList
			codeSnippets {
				lang
				code
		}
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
	markdownContent = RemoveInvisibleCharacters(markdownContent)
	if err != nil {
		panic(err)
	}

	// Output the problem details
	question := response.Data.Question

	found := false
	var code CodeSnippet
	for _, code = range question.CodeSnippets {
		if code.Lang == selectedLang {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("code snippet not found for this problem.")
	}

	solutionFolder := fmt.Sprintf(langDir + "/" + question.QuestionFrontendID)
	if err = os.MkdirAll(solutionFolder, os.ModePerm); err != nil {
		panic(err)
	}

	solutionFilePath := fmt.Sprintf(solutionFolder + "/solution." + fileExtension)
	fileContent := fmt.Sprintf("// Source: " + questionLink + "\n" + "// Author: " + author + "\n" + "// Difficulty: " + question.Difficulty + "\n\n" + "/*" + "\n" + markdownContent + "\n" + "*/" + "\n\n" + code.Code)

	if err = os.WriteFile(solutionFilePath, []byte(fileContent), os.ModePerm); err != nil {
		panic(err)
	}

	cmd := exec.Command("nvim", solutionFilePath)
	// Set the standard input, output, and error to match the parent process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Launching terminal application...")

	// Start the terminal application
	err = cmd.Run() // Run blocks until the subprocess exits
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var ans string
	fmt.Println("Link the solution to the Repo Readme?[y/N]")
	fmt.Scanln(&ans)
	if ans == "y" || ans == "Y" {
		file, err := os.OpenFile("Readme.md", os.O_APPEND|os.O_WRONLY, 0)

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		file.WriteString(fmt.Sprintf(`|%s|[%s](%s) | [%s](./%s)|%s|`, question.QuestionFrontendID, question.QuestionTitle, questionLink, code.Lang, solutionFilePath, question.Difficulty))
		fmt.Println("Done.")
	} else {
		fmt.Println("Bye.")
		return

	}

}
