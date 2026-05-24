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

type MemberStatus string

const (
	MemberStatusActive  MemberStatus = "active"
	MemberStatusRemoved MemberStatus = "removed"
	MemberStatusLeft    MemberStatus = "left"
)

type MemberSource string

const (
	MemberSourceApplication MemberSource = "application"
	MemberSourceInvitation  MemberSource = "invitation"
)

type InvitationStatus string

const (
	InvitationStatusPending   InvitationStatus = "pending"
	InvitationStatusAccepted  InvitationStatus = "accepted"
	InvitationStatusDeclined  InvitationStatus = "declined"
	InvitationStatusCancelled InvitationStatus = "cancelled"
	InvitationStatusExpired   InvitationStatus = "expired"
)
