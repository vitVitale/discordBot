package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"strings"
)

var (
	Token  = "YOUR_DISCORD_BOT_TOKEN_HERE"
	ApiKey = "YOUR_OPENAI_KEY_HERE"
)

type ChatGPTResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   ChatUsage    `json:"usage"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatUsage struct {
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
		chatResp, err := generateChatGPTResponse(m.Content)
		if err != nil {
			fmt.Println("Error generating chat response: ", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, chatResp.Choices[0].Message.Content)
	}
}

func generateChatGPTResponse(input string) (ChatGPTResponse, error) {
	var response ChatGPTResponse

	jsonValue, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": input},
		},
	})

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonValue))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
