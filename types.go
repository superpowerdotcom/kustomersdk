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
	MessageChannelNote  MessageChannel = "note"
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
	To                    []*EmailContact
	From                  *EmailContact
	Subject               string
	TemplateID            string
	BodyText              string
	BodyHTML              string
	Attachments           []*Attachment
	AutoCloseConversation bool
	CategoryTree          string
}

type SendSMSRequest struct {
	To                    string
	From                  string
	Body                  string
	Subject               string // Optional. If empty, will default to the first 250 characters of the body plus '...'
	AutoCloseConversation bool
	CategoryTree          string
}

type CreateConversationRequest struct {
	Direction ConversationDirection `json:"direction,omitempty"`
	Status    ConversationStatus    `json:"status,omitempty"`
	Name      string                `json:"name,omitempty"`
	// Used for setting custom variables, ie. html content, categoryTree, etc.
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

// I am not using the Customer type because this is used in Windmill and I want to minimize imapct
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
			ProgressiveStatus string                 `json:"progressiveStatus"`
			Verified          bool                   `json:"verified"`
			Rev               int                    `json:"rev"`
			Custom            map[string]interface{} `json:"custom"`
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

type SendNoteRequest struct {
	Body                  string
	Subject               string
	AutoCloseConversation bool
	CategoryTree          string
}

type CreateDraftNoteMessageRequest struct {
	Channel      MessageChannel `json:"channel"` // REQUIRED
	Conversation string         `json:"conversation,omitempty"`
	Body         string         `json:"body,omitempty"`
	SendAt       string         `json:"sendAt,omitempty"`
}

// ConversationListResponse represents the response from the conversations list endpoint
type ConversationListResponse struct {
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
	Data  []Conversation  `json:"data"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	PageSize int `json:"pageSize"`
	Page     int `json:"page"`
}

// PaginationLinks contains pagination navigation links
type PaginationLinks struct {
	Self  string  `json:"self"`
	First string  `json:"first"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// Conversation represents a Kustomer conversation
type Conversation struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id"`
	Attributes    ConversationAttributes `json:"attributes"`
	Relationships Relationships          `json:"relationships"`
	Links         Links                  `json:"links"`
}

// ConversationAttributes contains all conversation attributes
type ConversationAttributes struct {
	Name                     string                 `json:"name"`
	Preview                  string                 `json:"preview"`
	Channels                 []string               `json:"channels"`
	Status                   string                 `json:"status"`
	Open                     *StatusTimestamp       `json:"open,omitempty"`
	MessageCount             int                    `json:"messageCount"`
	LastReceivedAt           *time.Time             `json:"lastReceivedAt"`
	NoteCount                int                    `json:"noteCount"`
	Satisfaction             int                    `json:"satisfaction"`
	SatisfactionLevel        SatisfactionLevel      `json:"satisfactionLevel"`
	CreatedAt                time.Time              `json:"createdAt"`
	UpdatedAt                time.Time              `json:"updatedAt"`
	ModifiedAt               time.Time              `json:"modifiedAt"`
	LastActivityAt           time.Time              `json:"lastActivityAt"`
	Spam                     bool                   `json:"spam"`
	Ended                    bool                   `json:"ended"`
	EndedAt                  *time.Time             `json:"endedAt,omitempty"`
	EndedReason              string                 `json:"endedReason,omitempty"`
	EndedByType              string                 `json:"endedByType,omitempty"`
	Tags                     []interface{}          `json:"tags"`
	SuggestedTags            []interface{}          `json:"suggestedTags"`
	Predictions              []interface{}          `json:"predictions"`
	Sentiment                *Sentiment             `json:"sentiment,omitempty"`
	SuggestedShortcuts       []interface{}          `json:"suggestedShortcuts"`
	FirstMessageIn           *MessageSummary        `json:"firstMessageIn,omitempty"`
	FirstMessageOut          *MessageSummary        `json:"firstMessageOut,omitempty"`
	LastMessageIn            *MessageSummary        `json:"lastMessageIn,omitempty"`
	LastMessageOut           *MessageSummary        `json:"lastMessageOut,omitempty"`
	LastMessageAt            *time.Time             `json:"lastMessageAt"`
	LastMessageUnrespondedTo *MessageSummary        `json:"lastMessageUnrespondedTo,omitempty"`
	AssignedUsers            []string               `json:"assignedUsers"`
	AssignedTeams            []string               `json:"assignedTeams"`
	FirstResponse            *ResponseSummary       `json:"firstResponse,omitempty"`
	LastResponse             *ResponseSummary       `json:"lastResponse,omitempty"`
	Direction                string                 `json:"direction"`
	LastMessageDirection     string                 `json:"lastMessageDirection"`
	OutboundMessageCount     int                    `json:"outboundMessageCount"`
	InboundMessageCount      int                    `json:"inboundMessageCount"`
	Rev                      int                    `json:"rev"`
	Priority                 int                    `json:"priority"`
	BillingStatus            string                 `json:"billingStatus"`
	BillableReason           string                 `json:"billableReason"`
	Custom                   map[string]interface{} `json:"custom,omitempty"`
	Phase                    string                 `json:"phase"`
}

// StatusTimestamp represents a status with timestamp
type StatusTimestamp struct {
	StatusAt time.Time `json:"statusAt"`
}

// SatisfactionLevel represents satisfaction survey data
type SatisfactionLevel struct {
	SentByTeams []interface{} `json:"sentByTeams"`
	Answers     []interface{} `json:"answers"`
}

// Sentiment represents sentiment analysis data
type Sentiment struct {
	Polarity   int     `json:"polarity"`
	Confidence float64 `json:"confidence"`
}

// MessageSummary represents a summary of a message
type MessageSummary struct {
	ID             string                 `json:"id"`
	SentAt         *time.Time             `json:"sentAt,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	DirectionType  string                 `json:"directionType,omitempty"`
	Channel        string                 `json:"channel,omitempty"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
	CreatedBy      string                 `json:"createdBy,omitempty"`
	CreatedByTeams []string               `json:"createdByTeams,omitempty"`
}

// ResponseSummary represents response time summary
type ResponseSummary struct {
	ID             string     `json:"id"`
	Time           int        `json:"time"`
	BusinessTime   int        `json:"businessTime"`
	CreatedAt      time.Time  `json:"createdAt"`
	CreatedBy      string     `json:"createdBy,omitempty"`
	CreatedByTeams []string   `json:"createdByTeams,omitempty"`
	AssignedTeams  []string   `json:"assignedTeams,omitempty"`
	AssignedUsers  []string   `json:"assignedUsers,omitempty"`
	SentAt         *time.Time `json:"sentAt,omitempty"`
	ResponseTime   int        `json:"responseTime,omitempty"`
}

// MessageListResponse represents the response from the messages list endpoint
type MessageListResponse struct {
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
	Data  []Message       `json:"data"`
}

// Message represents a Kustomer message
type Message struct {
	Type          string            `json:"type"`
	ID            string            `json:"id"`
	Attributes    MessageAttributes `json:"attributes"`
	Relationships Relationships     `json:"relationships"`
	Links         MessageLinks      `json:"links"`
}

// MessageAttributes contains all message attributes
type MessageAttributes struct {
	ExternalID       string                 `json:"externalId"`
	Channel          string                 `json:"channel"`
	App              string                 `json:"app"`
	Size             int                    `json:"size"`
	Direction        string                 `json:"direction"`
	DirectionType    string                 `json:"directionType"`
	Location         *MessageLocation       `json:"location,omitempty"`
	Preview          string                 `json:"preview"`
	Sentiment        *Sentiment             `json:"sentiment,omitempty"`
	Meta             map[string]interface{} `json:"meta"`
	Status           string                 `json:"status"`
	AssignedTeams    []string               `json:"assignedTeams"`
	AssignedUsers    []string               `json:"assignedUsers"`
	Auto             bool                   `json:"auto"`
	SentAt           *time.Time             `json:"sentAt,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
	ModifiedAt       *time.Time             `json:"modifiedAt,omitempty"`
	Redacted         bool                   `json:"redacted"`
	CreatedByTeams   []string               `json:"createdByTeams"`
	Rev              int                    `json:"rev"`
	Reactions        []interface{}          `json:"reactions"`
	IntentDetections []interface{}          `json:"intentDetections"`
	Error            *MessageError          `json:"error,omitempty"`
}

// MessageLocation represents location information for a message
type MessageLocation struct {
	CountryName string `json:"countryName"`
	RegionName  string `json:"regionName"`
	CityName    string `json:"cityName"`
	ZipCode     string `json:"zipCode"`
}

// MessageError represents error information for failed messages
type MessageError struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// MessageLinks contains links for a message
type MessageLinks struct {
	Self    string `json:"self"`
	Related string `json:"related"`
}

// UserResponse represents the response from the user endpoint
type UserResponse struct {
	Data User `json:"data"`
}

// User represents a Kustomer user (internal practitioner)
type User struct {
	Type          string         `json:"type"`
	ID            string         `json:"id"`
	Attributes    UserAttributes `json:"attributes"`
	Relationships Relationships  `json:"relationships"`
	Links         Links          `json:"links"`
}

// UserAttributes contains all user attributes
type UserAttributes struct {
	DisplayName          string                 `json:"displayName"`
	Name                 string                 `json:"name"`
	UserType             string                 `json:"userType"`
	Email                string                 `json:"email"`
	EmailSignature       *string                `json:"emailSignature,omitempty"`
	EmailVerifiedAt      *time.Time             `json:"emailVerifiedAt,omitempty"`
	MostRecentLoginAt    *time.Time             `json:"mostRecentLoginAt,omitempty"`
	FirstEmailVerifiedAt *time.Time             `json:"firstEmailVerifiedAt,omitempty"`
	Password             UserPassword           `json:"password"`
	Mobile               *string                `json:"mobile,omitempty"`
	RoleGroups           []string               `json:"roleGroups"`
	Roles                []string               `json:"roles"`
	Custom               map[string]interface{} `json:"custom"`
	UpdatedAt            time.Time              `json:"updatedAt"`
	ModifiedAt           time.Time              `json:"modifiedAt"`
	CreatedAt            time.Time              `json:"createdAt"`
}

// UserPassword represents password settings for a user
type UserPassword struct {
	AllowNew  bool      `json:"allowNew"`
	ForceNew  bool      `json:"forceNew"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Customer represents a Kustomer customer (external user/patient)
type Customer struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    CustomerAttributes `json:"attributes"`
	Relationships Relationships      `json:"relationships"`
	Links         Links              `json:"links"`
}

// CustomerAttributes contains all customer attributes
type CustomerAttributes struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	DisplayColor string `json:"displayColor"`
	DisplayIcon  string `json:"displayIcon"`
	ExternalID   string `json:"externalId"`
	ExternalIDs  []struct {
		Verified         bool   `json:"verified"`
		ExternalVerified bool   `json:"externalVerified"`
		ExternalID       string `json:"externalId"`
	} `json:"externalIds"`
	SharedExternalIDs []interface{} `json:"sharedExternalIds"`
	Emails            []struct {
		Type             string `json:"type"`
		Verified         bool   `json:"verified"`
		ExternalVerified bool   `json:"externalVerified"`
		Email            string `json:"email"`
	} `json:"emails"`
	SharedEmails []interface{} `json:"sharedEmails"`
	Phones       []struct {
		Type     string `json:"type"`
		Phone    string `json:"phone"`
		Verified bool   `json:"verified"`
	} `json:"phones"`
	SharedPhones    []interface{} `json:"sharedPhones"`
	Whatsapps       []interface{} `json:"whatsapps"`
	SharedWhatsapps []interface{} `json:"sharedWhatsapps"`
	FacebookIDs     []interface{} `json:"facebookIds"`
	InstagramIDs    []interface{} `json:"instagramIds"`
	Socials         []interface{} `json:"socials"`
	SharedSocials   []interface{} `json:"sharedSocials"`
	URLs            []interface{} `json:"urls"`
	Locations       []interface{} `json:"locations"`
	ActiveUsers     []interface{} `json:"activeUsers"`
	Watchers        []interface{} `json:"watchers"`
	RecentLocation  struct {
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"recentLocation"`
	Locale                 interface{} `json:"locale"`
	BirthdayAt             *time.Time  `json:"birthdayAt"`
	Gender                 string      `json:"gender"`
	CreatedAt              time.Time   `json:"createdAt"`
	UpdatedAt              time.Time   `json:"updatedAt"`
	ModifiedAt             time.Time   `json:"modifiedAt"`
	LastActivityAt         time.Time   `json:"lastActivityAt"`
	LastCustomerActivityAt time.Time   `json:"lastCustomerActivityAt"`
	Deleted                bool        `json:"deleted"`
	ImportedAt             time.Time   `json:"importedAt"`
	LastMessageIn          struct {
		Sentiment interface{} `json:"sentiment"`
		Channel   string      `json:"channel"`
		SentAt    time.Time   `json:"sentAt"`
	} `json:"lastMessageIn"`
	LastMessageOut struct {
		SentAt time.Time `json:"sentAt"`
	} `json:"lastMessageOut"`
	LastMessageUnrespondedTo struct{}  `json:"lastMessageUnrespondedTo"`
	LastMessageAt            time.Time `json:"lastMessageAt"`
	LastConversation         struct {
		Sentiment struct {
			Polarity   int     `json:"polarity"`
			Confidence float64 `json:"confidence"`
		} `json:"sentiment"`
		ID       string        `json:"id"`
		Channels []string      `json:"channels"`
		Tags     []interface{} `json:"tags"`
	} `json:"lastConversation"`
	ConversationCounts struct {
		Open    int `json:"open"`
		Done    int `json:"done"`
		Snoozed int `json:"snoozed"`
		All     int `json:"all"`
	} `json:"conversationCounts"`
	Preview struct {
		Channel   string    `json:"channel"`
		Subject   string    `json:"subject"`
		Text      string    `json:"text"`
		Type      string    `json:"type"`
		PreviewAt time.Time `json:"previewAt"`
	} `json:"preview"`
	Tags      []interface{} `json:"tags"`
	Sentiment struct {
		Polarity   int     `json:"polarity"`
		Confidence float64 `json:"confidence"`
	} `json:"sentiment"`
	Custom            map[string]interface{} `json:"custom"`
	ProgressiveStatus string                 `json:"progressiveStatus"`
	Verified          bool                   `json:"verified"`
	Rev               int                    `json:"rev"`
	RecentItems       []struct {
		Meta      map[string]interface{} `json:"meta"`
		UpdatedAt time.Time              `json:"updatedAt"`
		Type      string                 `json:"type"`
		ID        string                 `json:"id"`
	} `json:"recentItems"`
	DefaultLang       string `json:"defaultLang"`
	SatisfactionLevel struct {
		FirstSatisfaction struct {
			SentByTeams []interface{} `json:"sentByTeams"`
		} `json:"firstSatisfaction"`
		LastSatisfaction struct {
			SentByTeams []interface{} `json:"sentByTeams"`
		} `json:"lastSatisfaction"`
	} `json:"satisfactionLevel"`
	RoleGroupVersions []interface{}          `json:"roleGroupVersions"`
	AccessOverride    []interface{}          `json:"accessOverride"`
	FirstName         string                 `json:"firstName"`
	LastName          string                 `json:"lastName"`
	UnmaskingWindow   map[string]interface{} `json:"unmaskingWindow"`
}

// GetCustomerByIDResponse represents the response from getting a customer by ID
type GetCustomerByIDResponse struct {
	Data Customer `json:"data"`
}
