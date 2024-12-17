package kustomersdk

import (
	"slices"

	"github.com/pkg/errors"
)

func validateCreateConversationRequest(req *CreateConversationRequest) error {
	if req.Direction != "" {
		ways := []ConversationDirection{ConversationDirectionIn, ConversationDirectionOut}
		if !slices.Contains(ways, req.Direction) {
			return errors.New("Direction is invalid")
		}
	}

	if req.Status != "" {
		statuses := []ConversationStatus{ConversationStatusOpen, ConversationStatusDone}
		if !slices.Contains(statuses, req.Status) {
			return errors.New("Status is invalid")
		}
	}

	return nil
}

func validateCreateDraftSMSMessageRequest(req *CreateDraftSMSMessageRequest) error {
	if req.Channel == "" {
		return errors.New("Channel cannot be empty")
	}

	channels := []MessageChannel{MessageChannelSMS}
	if !slices.Contains(channels, req.Channel) {
		return errors.New("Channel must be 'sms'")
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

func validateCreateDraftEmailMessageRequest(req *CreateDraftEmailMessageRequest) error {
	if req.Channel == "" {
		return errors.New("Channel cannot be empty")
	}

	channels := []MessageChannel{MessageChannelEmail}
	if !slices.Contains(channels, req.Channel) {
		return errors.New("Channel must be 'email'")
	}

	if req.Body == "" {
		return errors.New("Body cannot be empty")
	}

	if req.App != "" {
		apps := []MessageApp{MessageAppGmail}
		if !slices.Contains(apps, req.App) {
			return errors.New("App is invalid")
		}
	}

	if req.To == nil {
		return errors.New("To cannot be nil")
	}

	for _, recipient := range req.To {
		if recipient.Email == "" {
			return errors.New("To's email cannot be empty")
		}
	}

	if req.From == nil {
		return errors.New("From cannot be nil")
	}

	if req.From.Email == "" {
		return errors.New("From's email cannot be empty")
	}

	if req.From.Name == "" {
		return errors.New("From's name cannot be empty")
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

func validateSendEmailRequest(req *SendEmailRequest) error {
	if len(req.To) == 0 {
		return errors.New("To cannot be empty")
	}

	for _, recipient := range req.To {
		if recipient.Email == "" {
			return errors.New("To's email cannot be empty")
		}
		if recipient.Name == "" {
			return errors.New("To's name cannot be empty")
		}
	}

	if req.From == nil {
		return errors.New("From cannot be nil")
	}

	if req.From.Email == "" {
		return errors.New("From's email cannot be empty")
	}

	if req.From.Name == "" {
		return errors.New("From's name cannot be empty")
	}

	if req.Subject == "" {
		return errors.New("Subject cannot be empty")
	}

	if req.Attachments == nil {
		req.Attachments = []*Attachment{}
	}

	return nil
}

func validateSendSMSRequest(req *SendSMSRequest) error {
	if req.To == "" {
		return errors.New("To cannot be empty")
	}

	if req.From == "" {
		return errors.New("From cannot be empty")
	}

	if req.Body == "" {
		return errors.New("Body cannot be empty")
	}

	return nil
}
