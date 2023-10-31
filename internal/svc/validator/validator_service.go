package validator

import (
	"context"
	"fmt"

	"github.com/ArminGh02/go-auth-system/internal/broker"
	"github.com/ArminGh02/go-auth-system/internal/imagga"
	"github.com/ArminGh02/go-auth-system/internal/mailgun"
	"github.com/ArminGh02/go-auth-system/internal/model"
	"github.com/ArminGh02/go-auth-system/internal/repository"
	"github.com/ArminGh02/go-auth-system/internal/s3"
)

type Listener struct {
	broker  broker.Broker
	users   repository.User
	s3      *s3.S3
	imagga  *imagga.Imagga
	mailGun *mailgun.MailGun
}

func NewListener(u repository.User, s3 *s3.S3, b broker.Broker, i *imagga.Imagga, mg *mailgun.MailGun) *Listener {
	return &Listener{broker: b, s3: s3, users: u, imagga: i, mailGun: mg}
}

func (l *Listener) Listen() error {
	return l.broker.Subscribe("registrations", "validator", func(message []byte) error {
		nationalID := string(message)

		user, err := l.users.GetByNationalID(context.TODO(), nationalID)
		if err != nil {
			return fmt.Errorf("error getting user by national id: %w", err)
		}

		faceDetected, similarity, err := l.validateImages(user)
		if err != nil {
			return fmt.Errorf("error validating images: %w", err)
		}

		if faceDetected && similarity > 80 {
			user.Status = model.StatusAccepted
		} else {
			user.Status = model.StatusRejected
		}

		err = l.mailGun.Send(
			"arminghorbanian02@gmail.com",
			"Registration Status: "+user.Status.String(),
			buildMessageBody(user, faceDetected, similarity),
			user.Email,
		)
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}

		err = l.users.Update(context.Background(), user)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}

		return nil
	})
}

func (l *Listener) validateImages(user *model.User) (faceDetected bool, similarity float64, err error) {
	faceID1, faceID2, err := l.detectedFaces(user)
	if err != nil {
		return false, 0, err
	}
	if faceID1 == "" || faceID2 == "" {
		return false, 0, nil
	}

	similarity, err = l.imagga.FaceSimilarity(faceID1, faceID2)
	if err != nil {
		return false, 0, err
	}
	return true, similarity, nil
}

func (l *Listener) detectedFaces(user *model.User) (faceID1, faceID2 string, err error) {
	img1 := l.s3.URL(user.FirstImage())
	img2 := l.s3.URL(user.SecondImage())

	faceID1, err = l.imagga.DetectFace(img1)
	if err != nil {
		return "", "", err
	}

	faceID2, err = l.imagga.DetectFace(img2)
	if err != nil {
		return "", "", err
	}

	return faceID1, faceID2, nil
}

func buildMessageBody(user *model.User, faceDetected bool, similarity float64) string {
	const msg = `Dear %s,
Your registration status has been updated to: %s
%s
Best regards,
Auth System
`
	reason := ""
	if !faceDetected {
		reason = "Reason: No face detected in one or both images."
	} else if similarity <= 80 {
		reason = "Reason: The faces in images are not similar enough. Score = " + fmt.Sprintf("%.2f", similarity)
	}
	return fmt.Sprintf(msg, user.Name, user.Status, reason)
}
