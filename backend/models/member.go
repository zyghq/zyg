package models

import (
	"github.com/rs/xid"
	"github.com/zyghq/zyg"
	"strings"
	"time"
)

type Member struct {
	WorkspaceId string
	MemberId    string
	Name        string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Member) GenId() string {
	return "mm" + xid.New().String()
}

func (m Member) IsMemberSystem() bool {
	return m.Role == MemberRole{}.System()
}

func (m Member) AsMemberActor() MemberActor {
	return MemberActor{
		MemberId: m.MemberId,
		Name:     m.Name,
	}
}

func (m Member) CreateNewSystemMember(workspaceId string) Member {
	now := time.Now().UTC().UTC()
	return Member{
		MemberId:    m.GenId(), // generates a new ID
		WorkspaceId: workspaceId,
		Name:        "System",
		Role:        MemberRole{}.System(),
		CreatedAt:   now, // in same time space
		UpdatedAt:   now, // in same time space
	}
}

func (m Member) AvatarUrl() string {
	url := zyg.GetAvatarBaseURL()
	// url may or may not have a trailing slash
	// add a trailing slash if it doesn't have one
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url + m.MemberId
}

type MemberRole struct{}

func (mr MemberRole) Owner() string {
	return "owner"
}

func (mr MemberRole) System() string {
	return "system"
}

func (mr MemberRole) Admin() string {
	return "admin"
}

func (mr MemberRole) Support() string {
	return "support"
}

func (mr MemberRole) Viewer() string {
	return "viewer"
}

func (mr MemberRole) DefaultRole() string {
	return mr.Support()
}

func (mr MemberRole) IsValid(s string) bool {
	switch s {
	case mr.Owner(), mr.System(), mr.Admin(), mr.Support(), mr.Viewer():
		return true
	default:
		return false
	}
}
