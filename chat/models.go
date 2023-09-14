package chat

import (
	"fmt"

	"github.com/karim-w/go-azure-communication-services/identity"
)

type CreateChatThread struct {
	Topic        string        `json:"topic"`
	Participants []Participant `json:"participants"`
}

const _apiVersion = "2021-09-07"

var (
	ERR_UNAUTHORIZED      = fmt.Errorf("unauthorized")
	ERR_EXPIRED_TOKEN     = fmt.Errorf("token expired")
	ERR_NO_TOKEN_PROVIDED = fmt.Errorf("no token provided")
)

type Participant struct {
	CommunicationIdentifier identity.CommunicationIdentifier `json:"communicationIdentifier"`
	DisplayName             string                           `json:"displayName"`
}

type ChatUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type CreateChatThreadResponse struct {
	ChatThread          ChatThread           `json:"chatThread"`
	InvalidParticipants []InvalidParticipant `json:"invalidParticipants"`
}

type ChatThread struct {
	ID                               string                           `json:"id"`
	Topic                            string                           `json:"topic"`
	CreatedOn                        string                           `json:"createdOn"`
	CreatedByCommunicationIdentifier CreatedByCommunicationIdentifier `json:"createdByCommunicationIdentifier"`
}

type CreatedByCommunicationIdentifier struct {
	RawID             string            `json:"rawId"`
	CommunicationUser CommunicationUser `json:"communicationUser"`
}

type CommunicationUser struct {
	ID string `json:"id"`
}

type InvalidParticipant struct {
	Target  string `json:"target"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ChatMessageType string

const (
	ChatMessageType_Html               = "html"
	ChatMessageType_ParticipantAdded   = "participantAdded"
	ChatMessageType_ParticipantRemoved = "participantRemoved"
	ChatMessageType_Text               = "text"
	ChatMessageType_TopicUpdated       = "topicUpdated"
)

type SendChatMessageOptions struct {
	ChatThreadId string                 `json:"chatThreadId"`
	Request      SendChatMessageRequest `json:"request"`
}

type SendChatMessageRequest struct {
	Content           string            `json:"content"`
	Metadata          map[string]string `json:"metadata"`
	SenderDisplayName string            `json:"senderDisplayName"`
	Type              ChatMessageType   `json:"type"`
}

type SendChatMessageResponse struct {
	ID string `json:"id"`
}

type ChatMessageContent struct {
	InitiatorCommunicationIdentifier identity.CommunicationIdentifier `json:"initiatorCommunicationIdentifier"`
	Message                          string                           `json:"message"`
	Participants                     []Participant                    `json:"participants"`
	Topic                            string                           `json:"topic"`
}

type ChatMessage struct {
	Content                       ChatMessageContent               `json:"content"`
	CreatedOn                     string                           `json:"createdOn"`
	DeletedOn                     string                           `json:"deletedOn"`
	EditedOn                      string                           `json:"editedOn"`
	ID                            string                           `json:"id"`
	Metadata                      map[string]string                `json:"metadata"`
	SenderCommunicationIdentifier identity.CommunicationIdentifier `json:"senderCommunicationIdentifier"`
	SenderDisplayName             string                           `json:"senderDisplayName"`
	SequenceId                    string                           `json:"sequenceId"`
	Type                          ChatMessageType                  `json:"type"`
	Version                       string                           `json:"version"`
}

type ListChatMessagesOptions struct {
	ChatThreadId string `json:"chatThreadId"`
	MaxPageSize  int    `json:"maxPageSize"`
	StartTime    string `json:"startTime"`
}

type ChatMessagesCollection struct {
	NextLink string        `json:"nextLink"`
	Value    []ChatMessage `json:"value"`
}

type UpdateChatMessageOptions struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type ListChatThreadsOptions struct {
	MaxPageSize int    `json:"maxPageSize"`
	StartTime   string `json:"startTime"`
}

type ChatThreadsItem struct {
	DeletedOn             string `json:"deletedOn"`
	ID                    string `json:"id"`
	LastMessageReceivedOn string `json:"lastMessageReceivedOn"`
	Topic                 string `json:"topic"`
}

type ChatThreadsItemCollection struct {
	NextLink string            `json:"nextLink"`
	Value    []ChatThreadsItem `json:"value"`
}

type ListChatParticipantsOptions struct {
	ChatThreadId string `json:"chatThreadId"`
	MaxPageSize  int    `json:"maxPageSize"`
	Skip         int    `json:"skip"`
}

type ChatParticipant struct {
	CommunicationIdentifier identity.CommunicationIdentifier `json:"communicationIdentifier"`
	DisplayName             string                           `json:"displayName"`
	ShareHistoryTime        string                           `json:"shareHistoryTime"`
}

type ChatParticipantsCollection struct {
	NextLink string            `json:"nextLink"`
	Value    []ChatParticipant `json:"value"`
}
