package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/karim-w/go-azure-communication-services/chat"
	"github.com/karim-w/go-azure-communication-services/identity"
	// "github.com/stretchr/testify/assert"
)

func main() {
	// redirect stdout
	// file, err := os.OpenFile("./log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	file, err := os.OpenFile("./log", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	if err == nil {
		os.Stdout = file
		os.Stderr = file
	}

	resourceHost := "" // microsoft calls this endpoint
	accessKey := ""

	// roomsClient := rooms.NewClient(resourceHost, accessKey)
	identityClient := identity.New(resourceHost, accessKey)
	identity, err := identityClient.CreateIdentity(
		context.Background(),
		&identity.CreateIdentityOptions{
			CreateTokenWithScopes: []string{"chat"},
			ExpiresInMinutes:      60,
		},
	)
	log("identity", identity, err)

	chatClient, err := chat.New(resourceHost, accessKey)
	if err != nil {
		fmt.Println("err = ", err)
	}

	chatThreadRes, err := chatClient.CreateChatThread(
		context.Background(),
		"test",
		chat.ChatUser{ID: identity.ID, DisplayName: "test"},
	)
	log("chatThreadRes", chatThreadRes, err)

	chatThreadRes, err = chatClient.CreateChatThread(
		context.Background(),
		"test2",
		chat.ChatUser{ID: identity.ID, DisplayName: "test2"},
	)
	log("chatThreadRes", chatThreadRes, err)

	listChatThreadsRes, err := chatClient.ListChatThreads(
		context.Background(),
		&chat.ListChatThreadsOptions{},
	)
	log("listChatThreadsRes", listChatThreadsRes, err)

	listChatParticipantsRes, err := chatClient.ListChatParticipants(
		context.Background(),
		&chat.ListChatParticipantsOptions{
			ChatThreadId: chatThreadRes.ChatThread.ID,
		},
	)
	log("listChatParticipantsRes", listChatParticipantsRes, err)

	sendChatMessageRes, err := chatClient.SendChatMessage(
		context.Background(),
		&chat.SendChatMessageOptions{
			ChatThreadId: chatThreadRes.ChatThread.ID,
			Request: chat.SendChatMessageRequest{
				Content: "golang msg test",
				// Metadata: map[string]string{"key1": "value1", "key2": "value2"},
				Type: chat.ChatMessageType_Text,
			},
		},
	)
	log("sendChatMessageRes", sendChatMessageRes, err)

	getChatMessageRes, err := chatClient.GetChatMessage(
		context.Background(),
		sendChatMessageRes.ID,
		chatThreadRes.ChatThread.ID,
	)
	log("getChatMessageRes", getChatMessageRes, err)

	err = chatClient.UpdateChatMessages(
		context.Background(),
		sendChatMessageRes.ID,
		chatThreadRes.ChatThread.ID,
		&chat.UpdateChatMessageOptions{
			Content:  "golang msg test updated",
			Metadata: map[string]string{"key3": "value3"},
		},
	)
	log("UpdateChatMessages", nil, err)

	// time.Sleep(2 * time.Second)

	getChatMessageRes, err = chatClient.GetChatMessage(
		context.Background(),
		sendChatMessageRes.ID,
		chatThreadRes.ChatThread.ID,
	)
	log("getChatMessageRes", getChatMessageRes, err)

	listChatMessagesRes, err := chatClient.ListChatMessages(
		context.Background(),
		&chat.ListChatMessagesOptions{
			ChatThreadId: chatThreadRes.ChatThread.ID,
		},
	)
	log("listChatMessagesRes", listChatMessagesRes, err)

	err = chatClient.DeleteChatMessage(
		context.Background(),
		sendChatMessageRes.ID,
		chatThreadRes.ChatThread.ID,
	)
	log("deleteChatMessageRes", nil, err)

	listChatMessagesRes, err = chatClient.ListChatMessages(
		context.Background(),
		&chat.ListChatMessagesOptions{
			ChatThreadId: chatThreadRes.ChatThread.ID,
		},
	)
	log("listChatMessagesRes", listChatMessagesRes, err)
}

func log(varName string, res interface{}, err error) {
	formatted, _ := json.MarshalIndent(res, "", "    ")
	fmt.Println(varName, " = ", string(formatted))
	if err != nil {
		fmt.Println("err = ", err)
	}
}

//ed069156eae94f83a7e32a0071eedbbf
