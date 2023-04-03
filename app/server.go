package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	Token  = "YOUR_DISCORD_BOT_TOKEN_HERE"
	ApiKey = "YOUR_OPENAI_KEY_HERE"
)

type Gpt3Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text     string      `json:"text"`
	Index    int         `json:"index"`
	Logprobs interface{} `json:"logprobs"`
	Reason   string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.ToLower(m.Content) == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if strings.ToLower(m.Content) != "ping" && strings.ToLower(m.Content) != "pong" {
		gpt3Response, err := generateGpt3Response(m.Content)
		if err != nil {
			fmt.Println("Error generating GPT-3 response: ", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, gpt3Response.Choices[0].Text)
	}
}

func generateGpt3Response(input string) (Gpt3Response, error) {
	var response Gpt3Response

	jsonValue, _ := json.Marshal(map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      input,
		"max_tokens":  1024,
		"temperature": 0.5,
	})

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonValue))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
