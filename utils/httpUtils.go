package utils

import (
	"LLM-Chat/models/llm"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func RequestDeepSeek(request llm.ChatRequest) (*StreamResponse[llm.ChatCompletionChunk], error) {
	body := Marshal(request)
	log.Println("req body:" + string(body))

	headers := make(map[string]string)
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "Bearer " + apiKey

	return RequestAndParseStream[llm.ChatCompletionChunk]("https://api.deepseek.com/chat/completions", "POST", body, headers)
}

func RequestAndParseStream[T any](url string, method string, body []byte, headers map[string]string) (*StreamResponse[T], error) {
	client := &http.Client{}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errorText, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code: %d and respond with: %s", resp.StatusCode, errorText)
	}

	// 超时时间
	time.AfterFunc(time.Second*time.Duration(60), func() {
		// close the response body if timeout
		resp.Body.Close()
	})

	// debug
	//respBody, _ := ioutil.ReadAll(resp.Body)
	//log.Println(fmt.Println(string(respBody)))

	response := NewStreamResponse[T]()

	go func(response *StreamResponse[T]) {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			data := scanner.Bytes()
			if len(data) == 0 {
				continue
			}

			if bytes.HasPrefix(data, []byte("data:")) {
				// split
				data = data[5:]
			}

			// trim space
			data = bytes.TrimSpace(data)

			// unmarshal
			t := Unmarshal[T](data)

			response.Write(t)
		}
		response.Close()
	}(response)

	return response, nil
}
