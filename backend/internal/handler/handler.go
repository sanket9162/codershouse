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
	"github.com/sanket9162/codershouse/internal/utils"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Handler struct {
	App    *config.Config
	Logger *slog.Logger
}

func NewHandler(app *config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		App:    app,
		Logger: logger,
	}
}

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyOTPRequest struct {
	Phone     string `json:"Phone" validate:"required,e164"`
	OTP       string `json:"otp" validate:"required,numeric,len=4"`
	Hash      string `json:"hash" validate:"required"`
	ExpiresAt int64  `json:"expiresAt" validate:"required"`
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

	response := map[string]any{
		"message": "OTP verified successfully",
		"status":  http.StatusOK,
	}

	utils.WriteJSON(w, http.StatusOK, response)

}
