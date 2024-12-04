package kustomersdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/superpowerdotcom/go-lib-common/util"
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
}

type Options struct {
	APIKey string
}

type Kustomer struct {
	opts *Options
}

func New(opts *Options) (*Kustomer, error) {
	return &Kustomer{
		opts: opts,
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

	endpoint := BaseAPIURL + "/customers/externalId=" + externalID
	_, err := util.DoHTTP(ctx, "GET", endpoint, nil, resp, headers...)
	if err != nil {
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

	endpoint := BaseAPIURL + "/v1/customers/" + req.ID + "/conversation"
	_, err = util.DoHTTP(ctx, "POST", endpoint, body, resp, headers...)
	if err != nil {
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
	_, err = util.DoHTTP(ctx, "POST", endpoint, body, resp, headers...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create draft message for conversation '%s' and recipient '%s'", req.Conversation, req.To)
	}

	return resp, nil
}

// getHeaders returns all necessary headers for performing an HTTP call to Kustomer
func (k *Kustomer) getHeaders() []http.Header {
	return []http.Header{
		{"Authorization": []string{"Bearer " + k.opts.APIKey}},
		{"Content-Type": []string{"application/json"}},
	}
}
