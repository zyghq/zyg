package handler

import (
	"github.com/zyghq/zyg/internal/ports"
)

type WorkspaceHandler struct {
	ws ports.WorkspaceServicer
}

func NewWorkspaceHandler(ws ports.WorkspaceServicer) *WorkspaceHandler {
	return &WorkspaceHandler{ws: ws}
}
