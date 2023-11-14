package autostruct

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"

	"github.com/invopop/jsonschema"
)

var (
	Key           = ""
	DefaultClient = Client{
		Model:            "gpt-3.5-turbo-1106",
		Temperature:      0.5,
		MaxTokens:        4096,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
)

type FunctionCall struct {
	Arguments string `json:"arguments"`
}

type Message struct {
	FunctionCall FunctionCall `json:"function_call"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

type Client struct {
	Model            string
	Temperature      float64
	MaxTokens        int
	TopP             float64
	FrequencyPenalty float64
	PresencePenalty  float64
}

type CompletionRequest struct {
	Model            string                   `json:"model"`
	Messages         []map[string]string      `json:"messages"`
	Functions        []map[string]interface{} `json:"functions"`
	FunctionCall     map[string]string        `json:"function_call"`
	Temperature      float64                  `json:"temperature"`
	MaxTokens        int                      `json:"max_tokens"`
	TopP             float64                  `json:"top_p"`
	FrequencyPenalty float64                  `json:"frequency_penalty"`
	PresencePenalty  float64                  `json:"presence_penalty"`
}

type Wrapper struct {
	Wrapped interface{} `json:"wrapped"`
}

func Fill(prompt string, obj interface{}) error {
	return DefaultClient.Fill(prompt, obj)
}

func (client *Client) Fill(prompt string, obj interface{}) error {
	if Key == "" {
		return errors.New("no key set")
	}

	schema, shouldWrap, err := getSchema(obj)
	if err != nil {
		return err
	}

	completionRequest := CompletionRequest{
		Model: client.Model,
		Messages: []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		Functions: []map[string]interface{}{
			{
				"name":       "set_json",
				"parameters": schema,
			},
		},
		FunctionCall: map[string]string{
			"name": "set_json",
		},
		Temperature:      client.Temperature,
		MaxTokens:        client.MaxTokens,
		TopP:             client.TopP,
		FrequencyPenalty: client.FrequencyPenalty,
		PresencePenalty:  client.PresencePenalty,
	}

	res, err := makeRequest(completionRequest)
	if err != nil {
		return err
	}

	if len(res.Choices) == 0 {
		return errors.New("no choices returned from completions api call")
	}

	argument := []byte(res.Choices[0].Message.FunctionCall.Arguments)

	if shouldWrap {
		wrappedItem := &Wrapper{
			Wrapped: obj,
		}

		return json.Unmarshal(argument, wrappedItem)
	}

	return json.Unmarshal(argument, obj)
}

func makeRequest(postBody interface{}) (Response, error) {
	postBodyJson, err := json.Marshal(postBody)
	if err != nil {
		return Response{}, err
	}

	r, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(postBodyJson))
	if err != nil {
		return Response{}, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+Key)

	cli := http.Client{}

	resp, err := cli.Do(r)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	res := Response{}
	return res, json.Unmarshal(body, &res)
}

func getSchema(obj interface{}) (interface{}, bool, error) {
	var schema interface{}
	for _, v := range jsonschema.Reflect(obj).Definitions {
		schema = v
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return nil, false, errors.New("object must be a pointer")
	}

	pointsToValue := reflect.Indirect(reflect.ValueOf(obj))

	if pointsToValue.Kind() == reflect.Struct {
		return schema, false, nil
	}

	if pointsToValue.Kind() == reflect.Slice {
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"wrapped": map[string]interface{}{
					"type":                 "array",
					"items":                schema,
					"additionalProperties": false,
				},
			},
			"additionalProperties": false,
		}, true, nil
	}

	return nil, false, errors.New("type is not a struct or array")
}
