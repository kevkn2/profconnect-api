package entities

import "profconnect-api/internal/domain/constants"

type Project struct {
	ID          string                  `json:"id"`
	Professor   *Professor              `json:"professor"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Slots       int                     `json:"slots"`
	Status      constants.ProjectStatus `json:"status"`
}

type ProjectApplication struct {
	ID        string                      `json:"id"`
	Project  *Project                    `json:"project"`
	Student   *Student                    `json:"student"`
	Status    constants.ApplicationStatus `json:"status"`
	Message   string                      `json:"message"`
}
