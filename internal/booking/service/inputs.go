package service

import "github.com/google/uuid"

type CreateInput struct {
	SlotID               uuid.UUID
	UserID               uuid.UUID
	CreateConferenceLink bool
}
