package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/middleware"
	"github.com/sanket9162/codershouse/internal/models"
	"github.com/sanket9162/codershouse/internal/repository"
	"github.com/sanket9162/codershouse/internal/utils"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Handler struct {
	App    *config.Config
	Logger *slog.Logger
	DB     repository.DatabaseRepo
}

func NewHandler(app *config.Config, logger *slog.Logger, db repository.DatabaseRepo) *Handler {
	return &Handler{
		App:    app,
		Logger: logger,
		DB:     db,
	}
}

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyOTPRequest struct {
	Phone     string `json:"phone" validate:"required,e164"`
	OTP       string `json:"otp" validate:"required,numeric,len=4"`
	Hash      string `json:"hash" validate:"required"`
	ExpiresAt int64  `json:"expiresAt" validate:"required"`
}

type ActivateRequest struct {
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar" validate:""`
}

type CreateRoomRequest struct {
	Topic    string `json:"topic" validate:"required"`
	RoomType string `json:"roomType" validate:"required"`
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from codershouse"))
}

func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	// parse incoming phone number
	var req SendOTPRequest

	err := utils.ReadJSON(w, r, &req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// generate random 4-digit OTP
	n, _ := rand.Int(rand.Reader, big.NewInt(9000))
	otp := n.Int64() + 1000

	// Expiry time 5 minutes
	expiresAt := time.Now().Add(5 * time.Minute).Unix()

	// create data string
	data := fmt.Sprintf("%s.%d.%d", req.Phone, otp, expiresAt)

	// encrypt data
	hash := utils.Encrypt(data, h.App.SecretKey)

	// send OTP using twilio
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: h.App.TwilioSID,
		Password: h.App.TwilioToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(req.Phone)
	params.SetFrom(h.App.TwilioPhone)
	params.SetBody(fmt.Sprintf("Your coder's house OTP is %d ", otp))

	_, err = client.Api.CreateMessage(params)
	if err != nil {
		h.Logger.Error("Error sending OTP", "error", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	h.Logger.Info("OTP Generate,", "otp", otp, "phone", req.Phone)

	response := map[string]interface{}{
		"hash":      hash,
		"expiresAt": expiresAt,
		"phone":     req.Phone,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	// parse json
	var req VerifyOTPRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate struct
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// check expiry
	if time.Now().Unix() > req.ExpiresAt {
		utils.ErrorJSON(w, errors.New("otp expired"), http.StatusBadRequest)
		return
	}

	// convert otp string to int
	// otp, err := strconv.Atoi(req.OTP)
	// if err != nil {
	// 	utils.ErrorJSON(w, err, http.StatusBadRequest)
	// 	return
	// }

	otpInt := 0
	fmt.Sscanf(req.OTP, "%d", &otpInt)

	data := fmt.Sprintf("%s.%d.%d", req.Phone, otpInt, req.ExpiresAt)
	hash := utils.Encrypt(data, h.App.SecretKey)

	if req.Hash != hash {
		utils.ErrorJSON(w, errors.New("invalid otp"), http.StatusBadRequest)
		return
	}

	// check if user exists
	user, err := h.DB.GetUserByPhone(req.Phone)
	if err != nil {
		user = &models.User{
			Phone:     req.Phone,
			Activated: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := h.DB.CreateUser(user); err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		// fetch user again to get ID
		user, err = h.DB.GetUserByPhone(req.Phone)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	// generate tokens
	jwtUser := &utils.JwtUser{
		ID: user.ID.Hex(),
	}

	tokens, err := h.App.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// set access cookie
	accessCookie := h.App.Auth.GetAccessTokenCookie(tokens.Token)
	http.SetCookie(w, accessCookie)

	// set refresh cookie
	refreshCookie := h.App.Auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	response := map[string]any{
		"message":      "OTP verified successfully",
		"access_token": tokens.Token,
		"userID":       user.ID.Hex(),
		"phone":        user.Phone,
		"activated":    user.Activated,
	}

	utils.WriteJSON(w, http.StatusOK, response)

}

func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// parse json
	var req ActivateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// process image
	var avatarPath string
	if req.Avatar != "" {
		var err error
		avatarPath, err = utils.SaveProfileImage(req.Avatar)
		if err != nil {
			utils.ErrorJSON(w, errors.New("could not save image"), http.StatusInternalServerError)
			return
		}
	} else {
		avatarPath = "/images/monkey-avatar.png"
	}

	// Update user
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	// build full url for avatar
	baseURL := "http://localhost:" + h.App.Port
	// Or use h.App.Config.Domain if it's explicitly set to an absolute URL. Assuming localhost for dev:

	user.Name = req.Name
	user.Avatar = baseURL + avatarPath
	user.Activated = true

	if err := h.DB.UpdateUser(user); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message": "User activated successfully",
		"user":    user,
		"auth":    true,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.ErrorJSON(w, errors.New("unauthorized - no refresh token"), http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	// validate the token
	claims, err := h.App.Auth.ValidateToken(refreshToken)
	if err != nil {
		utils.ErrorJSON(w, errors.New("unauthorized - invalid refresh token"), http.StatusUnauthorized)
		return
	}

	// get the user id
	userID := claims.ID

	// Check if user exists
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
		return
	}

	jwtUser := &utils.JwtUser{
		ID: user.ID.Hex(),
	}

	tokenPairs, err := h.App.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		utils.ErrorJSON(w, errors.New("error generating tokens"), http.StatusInternalServerError)
		return
	}

	// set cookies
	accessCookie := h.App.Auth.GetAccessTokenCookie(tokenPairs.Token)
	http.SetCookie(w, accessCookie)

	refreshCookie := h.App.Auth.GetRefreshCookie(tokenPairs.RefreshToken)
	http.SetCookie(w, refreshCookie)

	response := map[string]any{
		"access_token": tokenPairs.Token,
		"user":         user,
		"auth":         true,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, h.App.Auth.GetAccessTokenCookie(""))
	http.SetCookie(w, h.App.Auth.GetRefreshCookie(""))

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Logged out successfully",
	})
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// parse json
	var req CreateRoomRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate struct
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// create room
	room := &models.Room{
		OwnerID:    userID,
		Topic:      req.Topic,
		RoomType:   req.RoomType,
		SpeakerIDs: []string{userID},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.DB.CreateRoom(room); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message": "Room created successfully",
		"room":    room,
		"auth":    true,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.DB.GetAllRooms("open")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, rooms)
}

func (h *Handler) GetRoomById(w http.ResponseWriter, r *http.Request) {
	// get room id from url
	roomId := chi.URLParam(r, "roomId")

	// get room from db
	room, err := h.DB.GetRoomById(roomId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, room)
}
