package kustomersdk

import (
	"slices"

	"github.com/pkg/errors"
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

func validateUpdateDraftMessageRequest(req *UpdateDraftMessageRequest) error {
	if req.Channel == "" {
		return errors.New("Channel cannot be empty")
	}

	if req.SendAt == "" {
		return errors.New("SendAt must be specified.")
	}

	return nil
}

func validateAttachmentRequest(req *CreateAttachmentRequest) error {
	if req.Name == "" {
		return errors.New("Name cannot be empty")
	}

	if req.ContentType == "" {
		return errors.New("ContentType cannot be nil")
	}

	if req.ContentLength == 0 {
		return errors.New("ContentLength cannot be 0")
	}

	return nil
}

func validateUploadAttachmentRequest(req *UploadAttachmentRequest) error {
	if req.UploadURL == "" {
		return errors.New("UploadURL cannot be empty")
	}

	if req.ContentType == "" {
		return errors.New("ContentType cannot be empty")
	}

	if len(req.Data) == 0 {
		return errors.New("Data cannot be empty")
	}

	return nil
}
