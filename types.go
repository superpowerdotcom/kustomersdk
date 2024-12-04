package kustomersdk

import "time"

type ConversationDirection string

const (
	ConversationDirectionIn  ConversationDirection = "in"
	ConversationDirectionOut ConversationDirection = "out"
)

type ConversationStatus string

const (
	ConversationStatusOpen   ConversationStatus = "open"
	ConversationStatusClosed ConversationStatus = "closed"
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
)

type CreateConversationRequest struct {
	ID        string                `json:"id"`
	Direction ConversationDirection `json:"direction,omitempty"`
	Status    ConversationStatus    `json:"status,omitempty"`
}

type CreateConversationResponse struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

type CreateDraftMessageRequest struct {
	App          MessageApp     `json:"app,omitempty"`
	Channel      MessageChannel `json:"channel"` // REQUIRED
	Conversation string         `json:"conversation,omitempty"`
	To           string         `json:"to,omitempty"`
	From         string         `json:"from,omitempty"`
	Body         string         `json:"body,omitempty"`
	SendAt       string         `json:"send_at,omitempty"`
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
