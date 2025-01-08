package main

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
