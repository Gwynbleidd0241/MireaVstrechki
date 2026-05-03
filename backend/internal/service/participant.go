package service

import (
	"errors"

	"go.uber.org/zap"

	"meeting-service/internal/model"
	"meeting-service/internal/notification"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrParticipantUserRequired = errors.New("participant user required")
	ErrInvalidParticipantRole  = errors.New("invalid participant role")
	ErrInvalidRSVPStatus       = errors.New("invalid rsvp status")
	ErrParticipantNotFound     = errors.New("participant not found")
)

type ParticipantService struct {
	participantRepo *postgres.ParticipantRepository
	eventRepo       *postgres.EventRepository
	userRepo        *postgres.UserRepository
	emailSender     *notification.EmailSender
	logger          *zap.Logger
}

func NewParticipantService(
	participantRepo *postgres.ParticipantRepository,
	eventRepo *postgres.EventRepository,
	userRepo *postgres.UserRepository,
	emailSender *notification.EmailSender,
	logger *zap.Logger,
) *ParticipantService {
	return &ParticipantService{
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
		userRepo:        userRepo,
		emailSender:     emailSender,
		logger:          logger,
	}
}

type AddParticipantRequest struct {
	EventID        int64
	UserID         int64
	Role           string
	ActingUserID   int64
	ActingUserRole string
}

type UpdateParticipantRequest struct {
	EventID        int64
	ParticipantID  int64
	Role           string
	ActingUserID   int64
	ActingUserRole string
}

type UpdateRSVPRequest struct {
	EventID        int64
	ParticipantID  int64
	RSVPStatus     string
	ActingUserID   int64
	ActingUserRole string
}

func (s *ParticipantService) Add(req AddParticipantRequest) (*model.Participant, error) {
	if err := s.checkManagePermission(req.EventID, req.ActingUserID, req.ActingUserRole); err != nil {
		return nil, err
	}

	if req.UserID <= 0 {
		return nil, ErrParticipantUserRequired
	}

	role := req.Role
	if role == "" {
		role = "participant"
	}

	if !isValidParticipantRole(role) {
		return nil, ErrInvalidParticipantRole
	}

	participant, err := s.participantRepo.Add(model.Participant{
		EventID:    req.EventID,
		UserID:     req.UserID,
		Role:       role,
		RSVPStatus: "pending",
	})
	if err != nil {
		return nil, err
	}

	// Отправляем приглашение асинхронно, чтобы не замедлять ответ API
	go s.sendInvite(req.EventID, req.UserID)

	return participant, nil
}

func (s *ParticipantService) sendInvite(eventID, userID int64) {
	if !s.emailSender.Enabled() {
		return
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil || event == nil {
		s.logger.Error("invite: failed to get event", zap.Int64("event_id", eventID), zap.Error(err))
		return
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		s.logger.Error("invite: failed to get user", zap.Int64("user_id", userID), zap.Error(err))
		return
	}

	if err := s.emailSender.SendInvite(notification.InviteData{
		RecipientEmail: user.Email,
		EventTitle:     event.Title,
		StartTime:      event.StartTime,
		Location:       event.Location,
		MeetingURL:     event.MeetingURL,
	}); err != nil {
		s.logger.Error("invite: failed to send email",
			zap.String("email", user.Email),
			zap.Error(err),
		)
	} else {
		s.logger.Info("invite sent", zap.String("email", user.Email), zap.Int64("event_id", eventID))
	}
}

func (s *ParticipantService) Update(req UpdateParticipantRequest) (*model.Participant, error) {
	participant, err := s.participantRepo.GetByID(req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if participant == nil || participant.EventID != req.EventID {
		return nil, ErrParticipantNotFound
	}

	if err := s.checkManagePermission(req.EventID, req.ActingUserID, req.ActingUserRole); err != nil {
		return nil, err
	}

	if !isValidParticipantRole(req.Role) {
		return nil, ErrInvalidParticipantRole
	}

	return s.participantRepo.UpdateRole(req.ParticipantID, req.Role)
}

func (s *ParticipantService) UpdateRSVP(req UpdateRSVPRequest) (*model.Participant, error) {
	participant, err := s.participantRepo.GetByID(req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if participant == nil || participant.EventID != req.EventID {
		return nil, ErrParticipantNotFound
	}

	if req.ActingUserRole != "admin" && participant.UserID != req.ActingUserID {
		return nil, ErrPermissionDenied
	}

	if !isValidRSVPStatus(req.RSVPStatus) {
		return nil, ErrInvalidRSVPStatus
	}

	return s.participantRepo.UpdateRSVP(req.ParticipantID, req.RSVPStatus)
}

func (s *ParticipantService) Remove(eventID, participantID, actingUserID int64, actingUserRole string) error {
	participant, err := s.participantRepo.GetByID(participantID)
	if err != nil {
		return err
	}

	if participant == nil || participant.EventID != eventID {
		return ErrParticipantNotFound
	}

	if err := s.checkManagePermission(eventID, actingUserID, actingUserRole); err != nil {
		return err
	}

	return s.participantRepo.Delete(participantID)
}

func (s *ParticipantService) ListByEventID(eventID int64) ([]model.Participant, error) {
	return s.participantRepo.ListByEventID(eventID)
}

func (s *ParticipantService) checkManagePermission(eventID, actingUserID int64, actingUserRole string) error {
	if actingUserRole == "admin" {
		return nil
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	if event.CreatorID != actingUserID {
		return ErrPermissionDenied
	}

	return nil
}

func isValidParticipantRole(role string) bool {
	return role == "participant" || role == "responsible"
}

func isValidRSVPStatus(status string) bool {
	return status == "pending" || status == "accepted" || status == "declined"
}
