// controllers/chat_controller.go
package controller

import (
	"chat/model"
	"chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"chat/service"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ChatController struct {
	Service     service.ChatService
	AuthService service.AuthService
}

var chats = make(map[uuid.UUID][]*websocket.Conn)
var mu sync.Mutex

func (handler *ChatController) handleConnections(chatID uuid.UUID, primaryChatID uint, username string, userid uint, conn *websocket.Conn) {

	defer func() {
		mu.Lock()
		for i, c := range chats[chatID] {
			if c == conn {
				chats[chatID] = append(chats[chatID][:i], chats[chatID][i+1:]...)
				break
			}
		}
		mu.Unlock()
		onlineUsersMutex.Lock()
		onlineUsers[userid] = false
		onlineUsersMutex.Unlock()
		conn.Close()
	}()
	mu.Lock()
	chats[chatID] = append(chats[chatID], conn)
	mu.Unlock()

	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			break
		}
		chatMessage := model.ChatMessage{
			Sender:   msg.Username,
			ChatID:   primaryChatID,
			SenderID: userid,
			Content:  msg.Text,
		}
		err = handler.Service.SaveMessage(&chatMessage)
		if err != nil {
			log.Printf("Error saving message: %v", err)
		}
		mu.Lock()
		for _, client := range chats[chatID] {
			if client != conn { // Skip the sender
				msg.Username = username
				msg.SenderID = userid
				msg.CreatedAt = time.Now()
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error: %v", err)
					client.Close()
				}
			}
		}
		mu.Unlock()
	}
}

var onlineUsers = make(map[uint]bool)
var onlineUsersMutex = sync.RWMutex{}

func (handler *ChatController) ServeWebSocket(ctx *gin.Context) {
	rawReceiverID := ctx.Query("receiverID")
	if rawReceiverID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Receiver id is empty")
		return
	}

	token := ctx.Query("token")
	if token == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "tokrn id is empty")
		return
	}
	receiverID, err := strconv.Atoi(rawReceiverID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid receiver id")
		return
	}
	usr, err := utils.ValidateJWT(token)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid token")
		return
	}
	senderID := int(usr)
	user, err := handler.AuthService.GetUserById(uint(senderID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving user")
		return
	}
	var chatID uuid.UUID
	var primaryChatID uint

	user1ID, user2ID := min(senderID, receiverID), max(senderID, receiverID)

	chat, _ := handler.Service.GetChatBySenderAndReceiverID(user1ID, user2ID)
	if chat == nil || chat.ID == 0 {
		chatID = uuid.New()
		if err := handler.Service.CreateChat(senderID, receiverID, chatID); err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Couldn't start a new chat")
			return
		}
		chat, restError := handler.Service.GetChatBySenderAndReceiverID(senderID, receiverID)
		if restError != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving chat")
			return
		}
		primaryChatID = chat.ID
	} else {
		chatID = chat.ChatID
		primaryChatID = chat.ID
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	onlineUsersMutex.Lock()
	onlineUsers[user.ID] = true
	onlineUsersMutex.Unlock()

	go handler.handleConnections(chatID, primaryChatID, user.Username, user.ID, ws)
}

func (handler *ChatController) GetChatHistory(ctx *gin.Context) {
	receiverID, err := strconv.Atoi(ctx.Param("receiverID"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "receiver id must be an integer")
		return
	}
	id, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "user_id not found in context")
		return
	}

	messages, restError := handler.Service.GetChatHistoryBetweenUsers(int(id.(uint)), receiverID)
	if restError != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Couldn't find chat history with this user")
		return
	}
	if messages == nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "No chat history found")
		return
	}
	utils.SuccessResponse(ctx, messages)
}

func (handler *ChatController) IsUserOnline(ctx *gin.Context) {
	idParam := ctx.Param("userID")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	onlineUsersMutex.RLock()
	online, exists := onlineUsers[uint(userID)]
	onlineUsersMutex.RUnlock()

	if exists && online {
		ctx.JSON(http.StatusOK, gin.H{"userID": userID, "online": true})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"userID": userID, "online": false})
	}
}
