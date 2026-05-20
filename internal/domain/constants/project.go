package constants

type ProjectStatus string

const (
	ProjectStatusOpen   ProjectStatus = "open"
	ProjectStatusClosed ProjectStatus = "closed"
)

type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusApproved  ApplicationStatus = "approved"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn ApplicationStatus = "withdrawn"
)
