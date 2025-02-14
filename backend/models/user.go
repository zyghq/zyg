package models

import (
	"encoding/json"
	"time"
)

type WorkOSUser struct {
	UserId            string    `json:"userId"`
	Email             string    `json:"email"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	EmailVerified     bool      `json:"emailVerified"`
	ProfilePictureURL *string   `json:"profilePictureUrl"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type WorkOSEvent struct {
	WorkOSUser WorkOSUser `json:"data"`
	Event      string     `json:"event"`
	CreatedAt  time.Time  `json:"created_at"`
}

func InitWorkOSEventFromPayload(payload []byte) (WorkOSEvent, error) {
	type payloadAux struct {
		Id    string `json:"id"`
		Event string `json:"event"`
		Data  struct {
			UserId            string    `json:"id"`
			Email             string    `json:"email"`
			FirstName         string    `json:"first_name"`
			LastName          string    `json:"last_name"`
			EmailVerified     bool      `json:"email_verified"`
			ProfilePictureURL *string   `json:"profile_picture_url"`
			CreatedAt         time.Time `json:"created_at"`
			UpdatedAt         time.Time `json:"updated_at"`
		} `json:"data"`
		CreatedAt time.Time `json:"created_at"`
	}

	var aux payloadAux
	err := json.Unmarshal(payload, &aux)
	if err != nil {
		return WorkOSEvent{}, err
	}
	return WorkOSEvent{
		WorkOSUser: WorkOSUser(aux.Data),
		Event:      aux.Event,
		CreatedAt:  aux.CreatedAt,
	}, nil
}

func (w *WorkOSUser) MarshalJSON() ([]byte, error) {
	type Alias WorkOSUser // Create an alias to avoid recursion
	return json.Marshal(&struct {
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		*Alias
	}{
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(w),
	})
}
