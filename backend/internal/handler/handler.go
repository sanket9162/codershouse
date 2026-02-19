package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/middleware"
	"github.com/sanket9162/codershouse/internal/models"
	"github.com/sanket9162/codershouse/internal/repository"
	"github.com/sanket9162/codershouse/internal/utils"
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
	Avatar string `json:"avatar" validate:"required"`
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
	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: h.App.TwilioSID,
	// 	Password: h.App.TwilioToken,
	// })

	// params := &twilioApi.CreateMessageParams{}
	// params.SetTo(req.Phone)
	// params.SetFrom(h.App.TwilioPhone)
	// params.SetBody(fmt.Sprintf("Your coder's house OTP is %d ", otp))

	// _, err = client.Api.CreateMessage(params)
	// if err != nil {
	// 	h.Logger.Error("Error sending OTP", "error", err)
	// 	utils.ErrorJSON(w, err, http.StatusInternalServerError)
	// 	return
	// }

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
	}

	// validate
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// process image

	// Update user
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Activated = true

	if err := h.DB.UpdateUser(user); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Response
	response := map[string]any{
		"message": "User activated successfully",
		"user":    user,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
