package handler

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/utils"
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
	Phone string `json:"phone" validate:"required,numeric,len=10"`
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
	otp := rand.Intn(9999)

	// Expiry time 5 minutes
	expiresAt := time.Now().Add(5 * time.Minute)

	// create data string
	data := fmt.Sprintf("%s.%s.%d", req.Phone, otp, expiresAt)

	// encrypt data
	hash := utils.Encrypt(data, h.App.SecretKey)
	if err != nil {
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
