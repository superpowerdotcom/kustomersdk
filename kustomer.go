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

	// CreateConversation creates a new conversation for a Kustomer
	CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error)

	// CreateDraftMessage creates a new draft message for a Kustomer. CreateConversation must be called
	// before calling this method, and the conversation ID must be passed in the request
	CreateDraftMessage(ctx context.Context, customerID string, req *CreateDraftMessageRequest) (*CreateDraftMessageResponse, error)

	// UpdateDraftMessage updates an existing draft message
	// This method is used to update the send time of a draft message after attachments have been uploaded
	UpdateDraftMessage(ctx context.Context, draftID string, req *UpdateDraftMessageRequest) (*UpdateDraftMessageResponse, error)

	// CreateDraftAttachment creates a new attachment for a draft message
	// The draft must first be created before calling this method
	CreateDraftAttachment(ctx context.Context, req *CreateAttachmentRequest) (*CreateAttachmentResponse, error)

	// UploadAttachment creates a new attachment for a draft message
	// CreateDraftAttachment must be called first to get the signed URL that  this request will be POSTed to
	UploadAttachment(ctx context.Context, req *UploadAttachmentRequest) error
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

func (k *Kustomer) CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error) {
	if err := validateCreateConversationRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateConversationRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateConversationRequest")
	}

	headers := k.getHeaders()
	resp := &CreateConversationResponse{}
	endpoint := fmt.Sprintf("%s/customers/%s/conversations", BaseAPIURL, req.ID)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create conversation for kustomer '%s'", req.ID)
	}

	return resp, nil
}

func (k *Kustomer) CreateDraftMessage(ctx context.Context, customerID string, req *CreateDraftMessageRequest) (*CreateDraftMessageResponse, error) {
	if err := validateCreateDraftMessageRequest(req); err != nil {
		return nil, errors.Wrap(err, "failed to validate CreateDraftMessageRequest")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal CreateDraftMessageRequest")
	}

	headers := k.getHeaders()
	resp := &CreateDraftMessageResponse{}
	endpoint := fmt.Sprintf("%s/customers/%s/drafts", BaseAPIURL, customerID)

	if _, err := k.call(ctx, endpoint, http.MethodPost, body, resp, headers...); err != nil {
		return nil, errors.Wrapf(err, "unable to create draft message for conversation '%s' and recipient '%s'", req.Conversation, req.To)
	}

	return resp, nil
}

func (k *Kustomer) UpdateDraftMessage(ctx context.Context, draftID string, req *UpdateDraftMessageRequest) (*UpdateDraftMessageResponse, error) {
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

func (k *Kustomer) CreateDraftAttachment(ctx context.Context, req *CreateAttachmentRequest) (*CreateAttachmentResponse, error) {
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

func (k *Kustomer) UploadAttachment(ctx context.Context, req *UploadAttachmentRequest) error {
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
	_ context.Context,
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
	request, err := http.NewRequest(method, endpoint, bodyBuffer)
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
