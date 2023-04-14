package entity

type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}
