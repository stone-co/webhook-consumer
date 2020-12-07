package usecase

import (
	"context"
	"fmt"

	"github.com/stone-co/webhook-consumer/pkg/domain"
	"gopkg.in/square/go-jose.v2"
)

func (uc NotificationUsecase) SendNotification(ctx context.Context, input domain.NotificationInput) error {
	encryptedPayload, err := uc.verify(input.EncryptedBody)
	if err != nil {
		return fmt.Errorf("unable to verify signature: %v", err)
	}

	payload, err := uc.decode(encryptedPayload)
	if err != nil {
		return fmt.Errorf("unable to decode payload: %v", err)
	}

	for _, notifier := range uc.notifiers {
		err := notifier.Send(ctx, input.Header.EventType, input.Header.EventID, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc NotificationUsecase) verify(signedBody string) (string, error) {
	obj, err := jose.ParseSigned(signedBody)
	if err != nil {
		return "", fmt.Errorf("unable to parse message: %v", err)
	}

	if len(obj.Signatures) != 1 {
		return "", fmt.Errorf("multi signature not supported")
	}

	// Verify will all keys.
	var plainText []byte
	for _, verificationKey := range uc.keys.VerificationKeyList {
		plainText, err = obj.Verify(verificationKey)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("invalid signature: %v", err)
	}

	return string(plainText), nil
}

func (uc NotificationUsecase) decode(encryptedBody string) (string, error) {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(encryptedBody)
	if err != nil {
		return "", fmt.Errorf("parsing encrypted: %v", err)
	}

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate the the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt(uc.keys.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("decrypting: %v", err)
	}

	return string(decrypted), nil
}
