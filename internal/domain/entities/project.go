package entities

import (
	"time"

	"profconnect-api/internal/domain/constants"
)

type Project struct {
	ID          string                  `json:"id"`
	Professor   *Professor              `json:"professor"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Slots       int                     `json:"slots"`
	Status      constants.ProjectStatus `json:"status"`
}

type ProjectApplication struct {
	ID      string                      `json:"id"`
	Project *Project                    `json:"project"`
	Student *Student                    `json:"student"`
	Status  constants.ApplicationStatus `json:"status"`
	Message string                      `json:"message"`
}

type ProjectMember struct {
	ID          string                 `json:"id"`
	Project     *Project               `json:"project"`
	Student     *Student               `json:"student"`
	Source      constants.MemberSource `json:"source"`
	SourceRefID string                 `json:"source_ref_id"`
	Status      constants.MemberStatus `json:"status"`
	JoinedAt    time.Time              `json:"joined_at"`
	LeftAt      *time.Time             `json:"left_at,omitempty"`
}

type ProjectInvitation struct {
	ID          string                     `json:"id"`
	Project     *Project                   `json:"project"`
	Student     *Student                   `json:"student"`
	Status      constants.InvitationStatus `json:"status"`
	Message     string                     `json:"message"`
	RespondedAt *time.Time                 `json:"responded_at,omitempty"`
}
