package kustomersdk

import (
	"github.com/pkg/errors"
	"slices"
)

func validateCreateConversationRequest(req *CreateConversationRequest) error {
	if req.ID == "" {
		return errors.New("ID cannot be empty")
	}

	if req.Direction != "" {
		ways := []ConversationDirection{ConversationDirectionIn, ConversationDirectionOut}
		if !slices.Contains(ways, req.Direction) {
			return errors.New("Direction is invalid")
		}
	}

	if req.Status != "" {
		statuses := []ConversationStatus{ConversationStatusOpen, ConversationStatusClosed}
		if !slices.Contains(statuses, req.Status) {
			return errors.New("Status is invalid")
		}
	}

	return nil
}

func validateCreateDraftMessageRequest(req *CreateDraftMessageRequest) error {
	if req.Channel == "" {
		return errors.New("Channel cannot be empty")
	}

	channels := []MessageChannel{MessageChannelSMS, MessageChannelEmail}
	if !slices.Contains(channels, req.Channel) {
		return errors.New("Channel must be one of 'sms', 'email'")
	}

	if req.Body == "" {
		return errors.New("Body cannot be empty")
	}

	if req.App != "" {
		apps := []MessageApp{MessageAppTwilio, MessageAppZipwhip}
		if !slices.Contains(apps, req.App) {
			return errors.New("App is invalid")
		}
	}

	return nil
}
