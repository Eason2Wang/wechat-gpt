package entity

type ChatGptReq struct {
	Prompt         string `xml:"prompt" json:"prompt"`
	ParentId       string `xml:"parentId" json:"parentId"`
	ConversationId string `xml:"conversationId" json:"conversationId"`
	Streamed       string `xml:"streamed" json:"streamed"`
}
