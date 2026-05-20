package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listApplicationsByProjectUsecase struct {
	projectRepository     port.ProjectRepository
	applicationRepository port.ProjectApplicationRepository
	professorRepository   port.ProfessorRepository
}

func NewListApplicationsByProjectUsecase(
	projectRepository port.ProjectRepository,
	applicationRepository port.ProjectApplicationRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.ListApplicationsByProjectInput, inputoutput.ListProjectApplicationsByProjectOutput] {
	return &listApplicationsByProjectUsecase{
		projectRepository:     projectRepository,
		applicationRepository: applicationRepository,
		professorRepository:   professorRepository,
	}
}

func (u *listApplicationsByProjectUsecase) Execute(ctx context.Context, input *inputoutput.ListApplicationsByProjectInput) (*inputoutput.ListProjectApplicationsByProjectOutput, error) {
	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	professor, err := u.professorRepository.GetByUserID(ctx, input.ProfessorUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load professor", err)
	}
	if professor == nil {
		return nil, domain.NotFound("professor profile not found")
	}

	if project.Professor == nil || project.Professor.ID != professor.ID {
		return nil, domain.Forbidden("you are not the owner of this project")
	}

	applications, err := u.applicationRepository.ListByProject(ctx, project.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list applications", err)
	}

	acceptedApplications := make([]*inputoutput.ProjectApplicationOutput, 0)
	pendingApplications := make([]*inputoutput.ProjectApplicationOutput, 0)
	for _, a := range applications {
		switch a.Status {
		case constants.ApplicationStatusApproved:
			acceptedApplications = append(acceptedApplications, &inputoutput.ProjectApplicationOutput{
				ID: a.ID,
				Status: string(a.Status),
				Message: a.Message,
				Student: &inputoutput.ProjectStudentBrief{
					Name:  a.Student.User.Name,
					Email: a.Student.User.Email,
					UserID: a.Student.User.ID,
					University: a.Student.University,
					Department: a.Student.Department,
					ResearchInterests: a.Student.ResearchInterests,
				},
				Project: &inputoutput.ProjectShortForApp{
					Title:       a.Project.Title,
					Description: a.Project.Description,
					Status:      string(a.Project.Status),
				},
			})
		case constants.ApplicationStatusPending:
			pendingApplications = append(pendingApplications, &inputoutput.ProjectApplicationOutput{
				ID: a.ID,
				Status: string(a.Status),
				Message: a.Message,
				Student: &inputoutput.ProjectStudentBrief{
					Name:  a.Student.User.Name,
					Email: a.Student.User.Email,
					UserID: a.Student.User.ID,
					University: a.Student.University,
					Department: a.Student.Department,
					ResearchInterests: a.Student.ResearchInterests,
				},
				Project: &inputoutput.ProjectShortForApp{
					Title:       a.Project.Title,
					Description: a.Project.Description,
					Status:      string(a.Project.Status),
				},
			})
		}
	}

	out := make([]*inputoutput.ProjectApplicationOutput, 0, len(applications))
	for _, a := range applications {
		out = append(
			out, 
			&inputoutput.ProjectApplicationOutput{
				ID:        a.ID,
				Status:    string(a.Status),
				Message:   a.Message,
				Student: &inputoutput.ProjectStudentBrief{
					Name:  a.Student.User.Name,
					Email: a.Student.User.Email,
					UserID: a.Student.User.ID,
					University: a.Student.University,
					Department: a.Student.Department,
					ResearchInterests: a.Student.ResearchInterests,
					StudentID: a.Student.ID,
				},
				Project: &inputoutput.ProjectShortForApp{
					Title:       a.Project.Title,
					Description: a.Project.Description,
					Status:      string(a.Project.Status),
				},
			},
		)
	}

	return &inputoutput.ListProjectApplicationsByProjectOutput{
		ApprovedApplications: acceptedApplications,
		PendingApplications:  pendingApplications,
	}, nil
}
