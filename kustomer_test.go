package kustomersdk

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newHttpClient returns *http.Client with Transport replaced to avoid making real calls
func newHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func KustomerWithMockResponse(statusCode int, responseBody []byte) *Kustomer {
	client := newHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}
	})

	return &Kustomer{
		opts: &Options{
			APIKey: "testing123",
		},
		Client: client,
	}
}

var _ = Describe("Kustomer SDK", func() {
	Context("New", func() {
		Context("when the options are valid", func() {
			It("should return a new Kustomer instance", func() {
				opts := &Options{
					APIKey: "abc123",
				}

				k, err := New(opts)
				Expect(err).To(BeNil())
				Expect(k).ToNot(BeNil())
			})
		})

		Context("when the options are invalid", func() {
			It("should return an error", func() {
				opts := &Options{}

				k, err := New(opts)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("APIKey cannot be empty"))
				Expect(k).To(BeNil())
			})
		})
	})

	Context("GetCustomerByExternalID", func() {
		Context("when the externalID is empty", func() {
			It("should return an error", func() {
				k, _ := New(&Options{
					APIKey: "abc123",
				})

				resp, err := k.GetCustomerByExternalID(nil, "")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("externalID cannot be empty"))
				Expect(resp).To(BeNil())
			})
		})

		Context("It returns user data", func() {
			It("should return a response", func() {
				apiResp, err := os.ReadFile("./test-assets/responses/get_customer_by_external_id.json")
				Expect(err).To(BeNil())

				k := KustomerWithMockResponse(http.StatusOK, apiResp)
				resp, err := k.GetCustomerByExternalID(context.Background(), "123")
				Expect(err).To(BeNil())
				Expect(resp).ToNot(BeNil())
				Expect(resp.Data.Attributes.Name).To(Equal("John Smith Jr."))
			})
		})
	})

	Context("createConversation", func() {
		It("creates a conversation", func() {
			apiResp, err := os.ReadFile("./test-assets/responses/create_conversation.json")
			Expect(err).To(BeNil())

			req := &CreateConversationRequest{
				Direction: ConversationDirectionOut,
				Status:    ConversationStatusOpen,
			}

			k := KustomerWithMockResponse(http.StatusCreated, apiResp)
			resp, err := k.createConversation(context.Background(), "123", req)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Data.Id).To(Equal("58228e05d4e96711004eccb7"))
		})
	})

	Context("createDraftSMSMessage", func() {
		It("creates a draft message", func() {
			apiResp, err := os.ReadFile("./test-assets/responses/create_draft_message.json")
			Expect(err).To(BeNil())

			req := &CreateDraftSMSMessageRequest{
				App:     MessageAppTwilio,
				Channel: MessageChannelSMS,
				Body:    "Please schedule your blood draw",
				SendAt:  time.Now().UTC().Add(time.Minute).Format(time.RFC3339),
			}

			k := KustomerWithMockResponse(http.StatusCreated, apiResp)
			resp, err := k.createDraftSMSMessage(context.Background(), "123", req)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Data.Id).To(Equal("5823e4d3bff98a100085320e"))
			Expect(resp.Data.Type).To(Equal("draft"))
		})
	})

	Context("updateDraftMessage", func() {
		It("updates a draft SMS message", func() {
			apiResp, err := os.ReadFile("./test-assets/responses/update_draft_message.json")
			Expect(err).To(BeNil())

			req := &UpdateDraftMessageRequest{
				Channel: MessageChannelSMS,
				SendAt:  time.Now().UTC().Add(time.Minute).Format(time.RFC3339),
			}

			k := KustomerWithMockResponse(http.StatusOK, apiResp)
			resp, err := k.updateDraftMessage(context.Background(), "123", req)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Data.Id).To(Equal("5823e4d3bff98a100085320e"))
			Expect(resp.Data.Type).To(Equal("draft"))
		})
	})

	Context("createDraftAttachment", func() {
		It("creates an attachment", func() {
			apiResp, err := os.ReadFile("./test-assets/responses/create_attachment.json")
			Expect(err).To(BeNil())

			req := &CreateAttachmentRequest{
				Name:          "requisition.pdf",
				ContentType:   "application/pdf",
				ContentLength: 12345,
			}

			k := KustomerWithMockResponse(http.StatusCreated, apiResp)
			resp, err := k.createDraftAttachment(context.Background(), "1234", req)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())

			// This is the URL where uploads will be POSTed to
			Expect(resp.Data.Links.Related).To(Equal("http://www.cool.com"))
		})
	})

	Context("chunkSMSMessage", func() {
		It("returns a single element when the message is less than 1600 chars", func() {
			msg := "This is a short message"
			chunks := chunkSMSMessage(msg)
			Expect(len(chunks)).To(Equal(1))
			Expect(chunks[0]).To(Equal(msg))
		})

		It("returns multiple elements when the message is greater than 1600 chars", func() {
			msg := strings.Repeat("a", 1600) + strings.Repeat("b", 1600) + strings.Repeat("c", 1600)
			chunks := chunkSMSMessage(msg)
			Expect(len(chunks)).To(Equal(3))
			Expect(chunks[0]).To(Equal(strings.Repeat("a", 1600)))
			Expect(chunks[1]).To(Equal(strings.Repeat("b", 1600)))
			Expect(chunks[2]).To(Equal(strings.Repeat("c", 1600)))
		})
	})
})
