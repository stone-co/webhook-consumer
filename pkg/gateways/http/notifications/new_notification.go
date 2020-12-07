package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http/responses"
)

const (
	EventIDHeader   = "X-Stone-Webhook-Event-Id"
	EventTypeHeader = "X-Stone-Webhook-Event-Type"
)

type NotificationRequest struct {
	EncryptedBody string `json:"encrypted_body" validate:"required"`
}

func (h Handler) New(w http.ResponseWriter, r *http.Request) {
	// Decode request body.
	var encryptedBody NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&encryptedBody); err != nil {
		h.log.WithError(err).Error("body is empty or has no valid fields")
		_ = responses.SendError(w, "body is empty or has no valid fields", http.StatusBadRequest)
		return
	}

	// Validate request body.
	if err := h.Validate(encryptedBody); err != nil {
		h.log.WithError(err).Error("invalid request body")
		_ = responses.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for mandatory headers.
	if r.Header.Get(EventIDHeader) == "" || r.Header.Get(EventTypeHeader) == "" {
		h.log.Errorf("%s and %s headers are mandatories", EventIDHeader, EventTypeHeader)
		_ = responses.SendError(w, fmt.Sprintf("%s and %s headers are mandatories", EventIDHeader, EventTypeHeader), http.StatusBadRequest)
		return
	}

	input := domain.NotificationInput{
		Header: domain.HeaderNotification{
			EventID:   r.Header.Get(EventIDHeader),
			EventType: r.Header.Get(EventTypeHeader),
		},
		EncryptedBody: encryptedBody.EncryptedBody,
	}

	// Call the usecase.
	if err := h.usecase.SendNotification(r.Context(), input); err != nil {
		h.log.WithError(err).Error("failed to send notification")
		_ = responses.SendError(w, "failed to send notification", http.StatusForbidden)
		return
	}

	_ = responses.Send(w, nil, http.StatusNoContent)
}
