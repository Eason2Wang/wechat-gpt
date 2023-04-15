package entity

type SavePromptRequest struct {
	Prompt string `json:"prompt"`
	UserId string `json:"userId"`
}
