package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/karim-w/go-azure-communication-services/client"
	"github.com/karim-w/go-azure-communication-services/identity"
	"github.com/karim-w/stdlib/httpclient"
)

type Chat interface {
	CreateChatThread(
		ctx context.Context,
		topic string,
		participants ...ChatUser,
	) (*CreateChatThreadResponse, error)
	DeleteChatThread(
		ctx context.Context,
		threadID string,
	) error
	AddChatParticipants(
		ctx context.Context,
		threadID string,
		participants ...ChatUser,
	) error
	RemoveChatParticipant(
		ctx context.Context,
		threadID string,
		acsId string,
	) error
	WithToken(
		token string,
		ExpiresAt time.Time,
	) Chat
	GetToken() (string, error)
	SetTokenFetcher(
		fetcher func() (string, error),
	)
	SendChatMessage(
		ctx context.Context,
		opts *SendChatMessageOptions,
	) (*SendChatMessageResponse, error)
	GetChatMessage(
		ctx context.Context,
		messageID string,
		threadID string,
	) (*ChatMessage, error)
	ListChatMessages(
		ctx context.Context,
		opts *ListChatMessagesOptions,
	) (*ChatMessagesCollection, error)
	DeleteChatMessage(
		ctx context.Context,
		messageID string,
		threadID string,
	) error
	UpdateChatMessages(
		ctx context.Context,
		messageID string,
		threadID string,
		opts *UpdateChatMessageOptions,
	) error
	ListChatThreads(
		ctx context.Context,
		opts *ListChatThreadsOptions,
	) (*ChatThreadsItemCollection, error)
	ListChatParticipants(
		ctx context.Context,
		opts *ListChatParticipantsOptions,
	) (*ChatParticipantsCollection, error)
}

type _chat struct {
	host         string
	client       *client.Client
	token        string
	validUntil   time.Time
	idc          *identity.Identity
	id           string
	tokenFetcher *func() (string, error)
}

func New(host string, key string) (Chat, error) {
	identityClient := identity.New(host, key)
	user, err := identityClient.CreateIdentity(
		context.Background(),
		&identity.CreateIdentityOptions{
			CreateTokenWithScopes: []string{"chat", "voip"},
			ExpiresInMinutes:      1440,
		},
	)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("failed to create identity")
	}
	client := client.New(key)
	return &_chat{
		host:       host,
		client:     client,
		token:      user.Token,
		validUntil: user.ExpiresOn,
		idc:        &identityClient,
		id:         user.ID,
	}, nil
}

func (c *_chat) refreshToken() error {
	client := *c.idc
	user, err := client.IssueAccessToken(
		context.Background(),
		c.id,
		&identity.IssueTokenOptions{
			ExpiresInMinutes: 1440,
			Scopes:           []string{"chat", "void"},
		},
	)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("failed to create identity")
	}
	c.token = user.Token
	c.validUntil = user.ExpiresOn
	return nil
}

func NewWithToken(host string, token string, expiresAt time.Time) (Chat, error) {
	return &_chat{
		host:       host,
		client:     nil,
		token:      token,
		validUntil: expiresAt,
	}, nil
}

func (c *_chat) getToken() (string, error) {
	if time.Now().After(c.validUntil) {
		if c.idc != nil && c.id != "" {
			err := c.refreshToken()
			if err != nil {
				return c.token, err
			} else {
				return "", ERR_EXPIRED_TOKEN
			}
		}
	}
	if c.token == "" {
		return "", ERR_NO_TOKEN_PROVIDED
	}
	return c.token, nil
}

func (c *_chat) GetToken() (string, error) {
	if c.tokenFetcher != nil {
		return (*c.tokenFetcher)()
	}
	return c.getToken()
}

func (c *_chat) SetTokenFetcher(
	fetcher func() (string, error),
) {
	c.tokenFetcher = &fetcher
}

func (c *_chat) CreateChatThread(
	ctx context.Context,
	topic string,
	participants ...ChatUser,
) (*CreateChatThreadResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	req := CreateChatThread{
		Topic: topic,
	}
	for _, p := range participants {
		req.Participants = append(req.Participants, Participant{
			CommunicationIdentifier: identity.CommunicationIdentifier{
				RawID: p.ID,
				CommunicationUser: identity.CommunicationUser{
					ID: p.ID,
				},
			},
			DisplayName: p.DisplayName,
		})
	}
	response := CreateChatThreadResponse{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		AddBody(req).Post()
	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}

func (c *_chat) DeleteChatThread(
	ctx context.Context,
	threadID string,
) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}

	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		Del()
	if res.IsSuccess() {
		return nil
	}
	if res.GetStatusCode() == 401 {
		return ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return err
}

func (c *_chat) WithToken(
	token string,
	ExpiresAt time.Time,
) Chat {
	c.token = token
	c.validUntil = ExpiresAt
	return c
}

func (c *_chat) AddChatParticipants(
	ctx context.Context,
	threadID string,
	participants ...ChatUser,
) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	req := []Participant{}
	for _, p := range participants {
		req = append(req, Participant{
			CommunicationIdentifier: identity.CommunicationIdentifier{
				RawID: p.ID,
				CommunicationUser: identity.CommunicationUser{
					ID: p.ID,
				},
			},
			DisplayName: p.DisplayName,
		})
	}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"/participants/:add?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		AddBody(map[string]interface{}{
			"participants": req,
		}).Post()
	if res.IsSuccess() {
		return nil
	}
	if res.GetStatusCode() == 401 {
		return ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return err
}

func (c *_chat) RemoveChatParticipant(
	ctx context.Context,
	threadID string,
	acsId string,
) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"/participants/:remove?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		AddBody(identity.CommunicationIdentifier{
			RawID: acsId,
			CommunicationUser: identity.CommunicationUser{
				ID: acsId,
			},
		}).Post()
	if res.IsSuccess() {
		return nil
	}
	if res.GetStatusCode() == 401 {
		return ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return err
}

func (c *_chat) SendChatMessage(
	ctx context.Context,
	opts *SendChatMessageOptions,
) (*SendChatMessageResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	req := SendChatMessageRequest{
		Content:           opts.Request.Content,
		Metadata:          opts.Request.Metadata,
		SenderDisplayName: opts.Request.SenderDisplayName,
		Type:              opts.Request.Type,
	}

	response := SendChatMessageResponse{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+opts.ChatThreadId+"/messages?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		AddBody(req).Post()

	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}

func (c *_chat) GetChatMessage(
	ctx context.Context,
	messageID string,
	threadID string,
) (*ChatMessage, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	response := ChatMessage{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"/messages/"+messageID+"?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").Get()

	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}

func (c *_chat) ListChatMessages(
	ctx context.Context,
	opts *ListChatMessagesOptions,
) (*ChatMessagesCollection, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	optionalParams := ""
	if opts.MaxPageSize > 0 {
		optionalParams += "maxPageSize=" + fmt.Sprint(opts.MaxPageSize) + "&"
	}
	if opts.StartTime != "" {
		optionalParams += "&startTime=" + opts.StartTime + "&"
	}
	response := ChatMessagesCollection{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+opts.ChatThreadId+"/messages?"+optionalParams+"api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").Get()

	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}

func (c *_chat) DeleteChatMessage(
	ctx context.Context,
	messageID string,
	threadID string,
) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}

	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"/messages/"+messageID+"?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").
		Del()
	if res.IsSuccess() {
		return nil
	}
	if res.GetStatusCode() == 401 {
		return ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return err
}

func (c *_chat) UpdateChatMessages(
	ctx context.Context,
	messageID string,
	threadID string,
	opts *UpdateChatMessageOptions,
) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}

	req := UpdateChatMessageOptions{
		Content:  opts.Content,
		Metadata: opts.Metadata,
	}

	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+threadID+"/messages/"+messageID+"?api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/merge-patch+json").
		AddBody(req).Patch()
	if res.IsSuccess() {
		return nil
	}
	fmt.Println(res)
	if res.GetStatusCode() == 401 {
		return ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return err
}

func (c *_chat) ListChatThreads(
	ctx context.Context,
	opts *ListChatThreadsOptions,
) (*ChatThreadsItemCollection, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	optionalParams := ""
	if opts.MaxPageSize > 0 {
		optionalParams += "maxPageSize=" + fmt.Sprint(opts.MaxPageSize) + "&"
	}
	if opts.StartTime != "" {
		optionalParams += "&startTime=" + opts.StartTime + "&"
	}
	response := ChatThreadsItemCollection{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads?"+optionalParams+"api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").Get()

	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}

func (c *_chat) ListChatParticipants(
	ctx context.Context,
	opts *ListChatParticipantsOptions,
) (*ChatParticipantsCollection, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	optionalParams := ""
	if opts.MaxPageSize > 0 {
		optionalParams += "maxPageSize=" + fmt.Sprint(opts.MaxPageSize) + "&"
	}
	if opts.Skip > 0 {
		optionalParams += "&skip=" + fmt.Sprint(opts.Skip) + "&"
	}
	response := ChatParticipantsCollection{}
	res := httpclient.Req(
		"https://"+c.host+"/chat/threads/"+opts.ChatThreadId+"/participants?"+optionalParams+"api-version="+_apiVersion,
	).AddBearerAuth(
		token,
	).AddHeader("Content-Type", "application/json").Get()

	if res.IsSuccess() {
		err := res.SetResult(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	if res.GetStatusCode() == 401 {
		return nil, ERR_UNAUTHORIZED
	}
	err = fmt.Errorf(string(res.GetBody()))
	return nil, err
}
