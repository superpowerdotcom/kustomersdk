package kustomersdk

import "time"

type ConversationDirection string

const (
	ConversationDirectionIn  ConversationDirection = "in"
	ConversationDirectionOut ConversationDirection = "out"
)

type ConversationStatus string

const (
	ConversationStatusOpen ConversationStatus = "open"
	ConversationStatusDone ConversationStatus = "done"
)

type MessageChannel string

const (
	MessageChannelSMS   MessageChannel = "sms"
	MessageChannelEmail MessageChannel = "email"
)

type MessageApp string

const (
	MessageAppTwilio  MessageApp = "twilio"
	MessageAppZipwhip MessageApp = "zipwhip"
	MessageAppGmail   MessageApp = "gmail"
)

type CustomAttrKey string

const (
	CustomAttrIsSubscribed                CustomAttrKey = "subscribedBool"
	CustomAttrIsIntakeCompleted           CustomAttrKey = "initialIntakeCompleteBool"
	CustomAttrIntakeReminderScheduledBool CustomAttrKey = "intakeReminderScheduledBool"
	CustomAttrInitialIntakeCompleteBool   CustomAttrKey = "initialIntakeCompleteBool"
)

type SendEmailRequest struct {
	To          []*EmailContact
	From        *EmailContact
	Subject     string
	TemplateID  string
	BodyText    string
	BodyHTML    string
	Attachments []*Attachment
}

type SendSMSRequest struct {
	To   string
	From string
	Body string
}

type CreateConversationRequest struct {
	Direction ConversationDirection `json:"direction,omitempty"`
	Status    ConversationStatus    `json:"status,omitempty"`

	// Used for setting custom variables, ie. html content
	Variables map[string]string `json:"custom,omitempty"`
}

type CreateConversationResponse struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

type CreateDraftSMSMessageRequest struct {
	App          MessageApp     `json:"app,omitempty"`
	Channel      MessageChannel `json:"channel"` // REQUIRED
	Conversation string         `json:"conversation,omitempty"`
	To           string         `json:"to,omitempty"`
	From         string         `json:"from,omitempty"`
	Body         string         `json:"body,omitempty"`
	SendAt       string         `json:"sendAt,omitempty"`
}

type CreateDraftEmailMessageRequest struct {
	App          MessageApp      `json:"app,omitempty"`
	Channel      MessageChannel  `json:"channel"` // REQUIRED
	Conversation string          `json:"conversation,omitempty"`
	To           []*EmailContact `json:"to,omitempty"`
	From         *EmailContact   `json:"from,omitempty"`
	Subject      string          `json:"subject,omitempty"`
	TemplateID   string          `json:"template,omitempty"`
	Body         string          `json:"body,omitempty"`
	SendAt       string          `json:"sendAt,omitempty"`
}

type EmailContact struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateDraftMessageResponse struct {
	Data struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Channel    string    `json:"channel"`
			Status     string    `json:"status"`
			CreatedAt  time.Time `json:"createdAt"`
			UpdatedAt  time.Time `json:"updatedAt"`
			ModifiedAt time.Time `json:"modifiedAt"`
		} `json:"attributes"`
		Relationships Relationships `json:"relationships"`
	} `json:"data"`
}

type UpdateDraftMessageRequest struct {
	Channel MessageChannel `json:"channel"` // REQUIRED
	SendAt  string         `json:"sendAt"`
}

type UpdateDraftMessageResponse struct {
	Data struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Channel    string    `json:"channel"`
			Body       string    `json:"body"`
			Status     string    `json:"status"`
			CreatedAt  time.Time `json:"createdAt"`
			UpdatedAt  time.Time `json:"updatedAt"`
			ModifiedAt time.Time `json:"modifiedAt"`
		} `json:"attributes"`
		Relationships Relationships `json:"relationships"`
	} `json:"data"`
}

type GetCustomerByExternalIDResponse struct {
	Data struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
			Name         string `json:"name"`
			DisplayName  string `json:"displayName"`
			DisplayColor string `json:"displayColor"`
			DisplayIcon  string `json:"displayIcon"`
			ExternalId   string `json:"externalId"`
			ExternalIds  []struct {
				ExternalId string      `json:"externalId"`
				Verified   bool        `json:"verified"`
				Id         interface{} `json:"id"`
			} `json:"externalIds"`
			SharedExternalIds []interface{} `json:"sharedExternalIds"`
			SignedUpAt        interface{}   `json:"signedUpAt"`
			AvatarUrl         interface{}   `json:"avatarUrl"`
			Username          interface{}   `json:"username"`
			Emails            []struct {
				Email    string      `json:"email"`
				Verified bool        `json:"verified"`
				Type     string      `json:"type"`
				Id       interface{} `json:"id"`
			} `json:"emails"`
			SharedEmails []struct {
				Email    string      `json:"email"`
				Verified bool        `json:"verified"`
				Type     string      `json:"type"`
				Id       interface{} `json:"id"`
			} `json:"sharedEmails"`
			Phones         []interface{} `json:"phones"`
			SharedPhones   []interface{} `json:"sharedPhones"`
			Socials        []interface{} `json:"socials"`
			SharedSocials  []interface{} `json:"sharedSocials"`
			Urls           []interface{} `json:"urls"`
			Locations      []interface{} `json:"locations"`
			ActiveUsers    []interface{} `json:"activeUsers"`
			Watchers       []interface{} `json:"watchers"`
			RecentLocation struct {
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"recentLocation"`
			Locale                 interface{} `json:"locale"`
			TimeZone               interface{} `json:"timeZone"`
			BirthdayAt             interface{} `json:"birthdayAt"`
			Gender                 interface{} `json:"gender"`
			Presence               string      `json:"presence"`
			CreatedAt              time.Time   `json:"createdAt"`
			UpdatedAt              time.Time   `json:"updatedAt"`
			ModifiedAt             time.Time   `json:"modifiedAt"`
			LastSeenAt             time.Time   `json:"lastSeenAt"`
			LastActivityAt         time.Time   `json:"lastActivityAt"`
			LastCustomerActivityAt time.Time   `json:"lastCustomerActivityAt"`
			LastMessageIn          struct {
				SentAt    time.Time   `json:"sentAt"`
				Sentiment interface{} `json:"sentiment"`
			} `json:"lastMessageIn"`
			LastConversation struct {
				Id        string `json:"id"`
				Sentiment struct {
					Confidence float64 `json:"confidence"`
					Polarity   int     `json:"polarity"`
				} `json:"sentiment"`
				Tags []interface{} `json:"tags"`
			} `json:"lastConversation"`
			ConversationCounts struct {
				All     int `json:"all"`
				Snoozed int `json:"snoozed"`
				Open    int `json:"open"`
				Done    int `json:"done"`
			} `json:"conversationCounts"`
			Preview struct {
				PreviewAt time.Time `json:"previewAt"`
				Type      string    `json:"type"`
				Text      string    `json:"text"`
				Subject   string    `json:"subject"`
			} `json:"preview"`
			Tags      []interface{} `json:"tags"`
			Sentiment struct {
				Polarity   int     `json:"polarity"`
				Confidence float64 `json:"confidence"`
			} `json:"sentiment"`
			ProgressiveStatus string `json:"progressiveStatus"`
			Verified          bool   `json:"verified"`
			Rev               int    `json:"rev"`
		} `json:"attributes"`
		Relationships Relationships `json:"relationships"`
		Links         Links         `json:"links"`
	} `json:"data"`
}

type Attachment struct {
	Name          string `json:"name"`
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
	Content       []byte `json:"-"`
}

type CreateAttachmentRequest struct {
	Name          string `json:"name"`
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
}

type CreateAttachmentResponse struct {
	Data struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name          string `json:"name"`
			ContentType   string `json:"contentType"`
			ContentLength int    `json:"contentLength"`
			Uploaded      bool   `json:"uploaded"`
			Redacted      bool   `json:"redacted"`
		} `json:"attributes"`
		Relationships Relationships `json:"relationships"`
		Links         struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"links"`
	} `json:"data"`
	Meta struct {
		Upload struct {
			Url    string `json:"url"`
			Fields struct {
				Key                  string `json:"key"`
				Acl                  string `json:"acl"`
				ContentType          string `json:"Content-Type"`
				XAmzMetaAttachmentId string `json:"X-Amz-Meta-Attachment-Id"`
				Bucket               string `json:"bucket"`
				XAmzAlgorithm        string `json:"X-Amz-Algorithm"`
				XAmzCredential       string `json:"X-Amz-Credential"`
				XAmzDate             string `json:"X-Amz-Date"`
				Policy               string `json:"Policy"`
				XAmzSignature        string `json:"X-Amz-Signature"`
			} `json:"fields"`
		} `json:"upload"`
	} `json:"meta"`
}

type UploadAttachmentRequest struct {
	// UploadURL is the presigned URL to upload the attachment
	// this is obtained from CreateAttachmentResponse.Data.Links.Related
	UploadURL string `json:"-"`

	// ContentType is the mimetype of the file
	ContentType string `json:"-"`

	// Data is the actual file contents to upload
	Data []byte `json:"-"`
}

// ----------------- Common shared structures between requests

type Links struct {
	Self string `json:"self"`
}

type RelationshipData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

// Relationships is a struct that contains the OpenAPI relationship associations
// This structure is shared between multiple requests
type Relationships struct {
	Org struct {
		Links Links            `json:"links"`
		Data  RelationshipData `json:"data"`
	} `json:"org,omitempty"`
	Customer struct {
		Links Links            `json:"links"`
		Data  RelationshipData `json:"data"`
	} `json:"customer,omitempty"`
	Messages struct {
		Links Links `json:"links"`
	} `json:"messages,omitempty"`
	CreatedBy struct {
		Links Links            `json:"links"`
		Data  RelationshipData `json:"data"`
	} `json:"createdBy,omitempty"`
	ModifiedBy struct {
		Links Links            `json:"links"`
		Data  RelationshipData `json:"data"`
	} `json:"modifiedBy,omitempty"`
}

type GetCustomerAttributesResponse struct {
	Data struct {
		Attributes struct {
			Custom map[CustomAttrKey]interface{} `json:"custom"`
		} `json:"attributes"`
	} `json:"data"`
}

type SetCustomerAttributesRequest struct {
	Custom map[CustomAttrKey]interface{} `json:"custom"`
}

type SetCustomerAttributesResponse struct {
	// TODO: what does this return?
}
