package handler

import "time"

type PATReqPayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type WorkspaceReqPayload struct {
	Name string `json:"name"`
}

type CrLabelReqPayload struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type CrLabelRespPayload struct {
	LabelId   string `json:"labelId"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
