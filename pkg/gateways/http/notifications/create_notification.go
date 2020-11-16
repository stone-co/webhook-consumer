package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http/responses"
)

type CreateNotificationRequest struct {
	Text string `json:"text" validate:"required"`
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Decode request body.
	var body CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.WithError(err).Error("body is empty or has no valid fields")
		response := responses.Error{Message: "body is empty or has no valid fields"}
		_ = responses.Send(w, response, http.StatusBadRequest)
		return
	}

	// Validate request body.
	if err := h.Validate(body); err != nil {
		h.log.WithError(err).Error("invalid request body")
		response := responses.Error{Message: err.Error()}
		_ = responses.Send(w, response, http.StatusBadRequest)
		return
	}

	input := domain.CreateNotificationInput{
		Body: body.Text,
	}

	// Call the usecase.
	err := h.usecase.CreateNotification(r.Context(), input)
	if err != nil {
		h.log.WithError(err).Error("failed to create notification")
		response := responses.Error{Message: "failed to create notification"}
		_ = responses.Send(w, response, http.StatusInternalServerError)
		return
	}

	_ = responses.Send(w, nil, http.StatusNoContent)
}
