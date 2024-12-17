package kustomersdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

const BaseAPIURL = "https://api.kustomerapp.com/v1"

type IKustomer interface {
	// GetCustomerByExternalID retrieves a customer by their external ID, which in our case is the database ID
	GetCustomerByExternalID(ctx context.Context, databaseID string) (*GetCustomerByExternalIDResponse, error)

	// SendEmail sends an email to a Kustomer
	SendEmail(ctx context.Context, customerID string, req *SendEmailRequest) error

	// SendSMS sends an SMS to a Kustomer
	SendSMS(ctx context.Context, customerID string, req *SendSMSRequest) error
}

type Options struct {
	APIKey string
	Client *http.Client
}

type Kustomer struct {
	opts   *Options
	Client *http.Client
}

func New(opts *Options) (*Kustomer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, errors.Wrap(err, "unable to validate options")
	}

	// Client is an option so that we can inject a mock client for testing
	// or if we want to use a custom client for any reason
	if opts.Client == nil {
		opts.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Kustomer{
		opts:   opts,
		Client: opts.Client,
	}, nil
}

func validateOptions(opts *Options) error {
	if opts.APIKey == "" {
		return errors.New("APIKey cannot be empty")
	}

	return nil
}

func (k *Kustomer) GetCustomerByExternalID(ctx context.Context, externalID string) (*GetCustomerByExternalIDResponse, error) {
	if externalID == "" {
		return nil, errors.New("externalID cannot be empty")
	}

	headers := k.getHeaders()
	resp := &GetCustomerByExternalIDResponse{}

	endpoint := fmt.Sprintf("%s/customers/externalId=%s", BaseAPIURL, externalID)

	if _, err := k.call(ctx, endpoint, http.MethodGet, nil, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to get customer by external ID '%s'", externalID)
	}

	return resp, nil
}

func (k *Kustomer) SendEmail(ctx context.Context, customerID string, req *SendEmailRequest) error {
	if err := validateSendEmailRequest(req); err != nil {
		return errors.Wrap(err, "failed to validate SendEmailRequest")
	}

	if req.Attachments == nil {
		req.Attachments = make([]*Attachment, 0)
	}

	// Create Conversation
	conversationReq := &CreateConversationRequest{
		Direction: ConversationDirectionOut,
		Status:    ConversationStatusOpen,
		Variables: map[string]string{},
	}
	if req.BodyHTML != "" {
		conversationReq.Variables["htmlContentStr"] = req.BodyHTML
	}

	conversationResp, err := k.createConversation(ctx, customerID, conversationReq)
	if err != nil {
		return errors.Wrap(err, "unable to create conversation")
	}

	// Create Draft Email Message
	draftReq := &CreateDraftEmailMessageRequest{
		Channel:      MessageChannelEmail,
		App:          MessageAppGmail,
		Conversation: conversationResp.Data.Id,
		To:           req.To,
		From:         req.From,
		Body:         req.BodyText,
		TemplateID:   req.TemplateID,
		Subject:      req.Subject,
	}

	draftResp, err := k.createDraftEmailMessage(ctx, customerID, draftReq)
	if err != nil {
		return errors.Wrap(err, "unable to create draft email message")
	}

	for _, att := range req.Attachments {
		contentType := http.DetectContentType(att.Content)

		// Create Attachment Request, this will return a pre-signed URL
		attachmentResp, err := k.createDraftAttachment(ctx, &CreateAttachmentRequest{
			Name:          att.Name,
			ContentType:   att.ContentType,
			ContentLength: int64(len(att.Content)),
		})
		if err != nil {
			return errors.Wrap(err, "unable to create draft attachment")
		}

		// Perform the actual upload to the pre-signed URL
		uploadErr := k.uploadAttachment(ctx, &UploadAttachmentRequest{
			UploadURL:   attachmentResp.Data.Links.Related,
			ContentType: contentType,
			Data:        att.Content,
		})
		if uploadErr != nil {
			return errors.Wrap(uploadErr, "unable to upload attachment")
		}
	}

	sendAt := time.Now().UTC().Add(time.Minute).Format(time.RFC3339)

	// Update the draft message with the sendAt
	_, err = k.updateDraftMessage(ctx, draftResp.Data.Id, &UpdateDraftMessageRequest{
		SendAt:  sendAt,
		Channel: MessageChannelEmail,
	})
	if err != nil {
		return errors.Wrap(err, "unable to update draft message")
	}

	return nil
}

func (k *Kustomer) SendSMS(ctx context.Context, customerID string, req *SendSMSRequest) error {
	if err := validateSendSMSRequest(req); err != nil {
		return errors.Wrap(err, "failed to validate SendSMSRequest")
	}

	// Create a conversation under the customer, and then we'll take
	// the resulting conversation ID and make a draft
	convoResp, err := k.createConversation(ctx, customerID, &CreateConversationRequest{
		Direction: ConversationDirectionOut,
		Status:    ConversationStatusOpen,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create conversation for kustomer '%s'", customerID)
	}

	// Create Draft Message with a timestamp 1 minute into the future.
	// Kustomer requires it to be sent this way for some silly reason.
	_, err = k.createDraftSMSMessage(ctx, customerID, &CreateDraftSMSMessageRequest{
		App:          MessageAppTwilio,
		Channel:      MessageChannelSMS,
		Conversation: convoResp.Data.Id,
		To:           req.To,
		From:         req.From,
		Body:         req.Body,
		SendAt:       time.Now().UTC().Add(time.Minute).Format(time.RFC3339),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create draft message for conversation '%s'", convoResp.Data.Id)
	}

	return nil
}

func (k *Kustomer) createConversation(ctx context.Context, customerID string, req *CreateConversationRequest) (*CreateConversationResponse, error) {
	if err := validateCreateConversationRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateConversationRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateConversationRequest")
	}

	headers := k.getHeaders()
	resp := &CreateConversationResponse{}
	endpoint := fmt.Sprintf("%s/customers/%s/conversations", BaseAPIURL, customerID)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create conversation for kustomer '%s'", customerID)
	}

	return resp, nil
}

func (k *Kustomer) createDraftSMSMessage(ctx context.Context, customerID string, req *CreateDraftSMSMessageRequest) (*CreateDraftMessageResponse, error) {
	if err := validateCreateDraftSMSMessageRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateDraftSMSMessageRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateDraftSMSMessageRequest")
	}

	headers := k.getHeaders()
	resp := &CreateDraftMessageResponse{}
	endpoint := fmt.Sprintf("%s/customers/%s/drafts", BaseAPIURL, customerID)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create draft message for conversation '%s' and recipient '%s'", req.Conversation, req.To)
	}

	return resp, nil
}

func (k *Kustomer) createDraftEmailMessage(ctx context.Context, customerID string, req *CreateDraftEmailMessageRequest) (*CreateDraftMessageResponse, error) {
	if err := validateCreateDraftEmailMessageRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateDraftEmailMessageRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateDraftEmailMessageRequest")
	}

	headers := k.getHeaders()
	resp := &CreateDraftMessageResponse{}
	endpoint := fmt.Sprintf("%s/customers/%s/drafts", BaseAPIURL, customerID)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create draft message for conversation '%s' and recipient '%s'", req.Conversation, req.To[0].Email)
	}

	return resp, nil
}

func (k *Kustomer) updateDraftMessage(ctx context.Context, draftID string, req *UpdateDraftMessageRequest) (*UpdateDraftMessageResponse, error) {
	if err := validateUpdateDraftMessageRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate UpdateDraftMessageRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal UpdateDraftMessageRequest")
	}

	headers := k.getHeaders()
	resp := &UpdateDraftMessageResponse{}
	endpoint := fmt.Sprintf("%s/drafts/%s", BaseAPIURL, draftID)

	if _, err := k.call(ctx, endpoint, http.MethodPut, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to update draft message with ID '%s'", draftID)
	}

	return resp, nil
}

func (k *Kustomer) createDraftAttachment(ctx context.Context, req *CreateAttachmentRequest) (*CreateAttachmentResponse, error) {
	if err := validateAttachmentRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateAttachmentRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateAttachmentRequest")
	}

	headers := k.getHeaders()
	resp := &CreateAttachmentResponse{}
	endpoint := fmt.Sprintf("%s/attachments", BaseAPIURL)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create attachment for file '%s'", req.Name)
	}

	return resp, nil
}

func (k *Kustomer) uploadAttachment(ctx context.Context, req *UploadAttachmentRequest) error {
	if err := validateUploadAttachmentRequest(req); err != nil {
		return errors.Wrap(err, "failed to validate UploadAttachmentRequest")
	}

	headers := []http.Header{
		{"Content-Type": []string{req.ContentType}},
		{"Accept": []string{"application/json"}},
	}

	if _, err := k.call(ctx, req.UploadURL, http.MethodPut, req.Data, nil, headers...); err != nil {
		return errors.Wrapf(err, "unable to upload file to destination '%s'", req.UploadURL)
	}

	return nil
}

// getHeaders returns all necessary headers for performing an HTTP call to Kustomer
func (k *Kustomer) getHeaders() []http.Header {
	return []http.Header{
		{"Authorization": []string{"Bearer " + k.opts.APIKey}},
		{"Content-Type": []string{"application/json"}},
	}
}

func (k *Kustomer) call(
	ctx context.Context,
	endpoint,
	method string,
	requestBody []byte,
	target any,
	headers ...http.Header,
) (*http.Response, error) {

	if target != nil {
		if reflect.ValueOf(target).Kind() != reflect.Ptr {
			return nil, errors.New("target must be a pointer")
		}
	}

	// Automatically handles nil requestBody
	bodyBuffer := bytes.NewBuffer(requestBody)

	// Generate request
	request, err := http.NewRequestWithContext(ctx, method, endpoint, bodyBuffer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// Set headers
	for _, h := range headers {
		for key, values := range h {
			for _, value := range values {
				request.Header.Set(key, value)
			}
		}
	}

	// Perform the request
	resp, err := k.Client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}

	defer resp.Body.Close()

	body, err := getResponseBody(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-200 status code: %d; resp body: %s", resp.StatusCode, string(body))
	}

	// If there is no target, we are done
	if target == nil {
		return resp, nil
	}

	if err := json.Unmarshal(body, target); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return resp, nil
}

func getResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("response cannot be nil")
	}

	if resp.Body == nil {
		return nil, errors.New("response body cannot be nil")
	}

	defer resp.Body.Close()

	// Read the response body into a buffer so we can re-add the body back into
	// the response for others to use
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response body")
	}

	// Re-create body so it can be read again
	resp.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

	return buf.Bytes(), nil
}
