package dto

import (
	"fmt"
)

// Error api errors
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}

// Workspace DTO
type Workspace struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	ImageURL    string            `json:"imageUrl"`
	Settings    WorkspaceSettings `json:"workspaceSettings"`
	HourlyRate  HourlyRate        `json:"hourlyRate"`
	Memberships []Membership
}

// Membership DTO
type Membership struct {
	HourlyRate HourlyRate       `json:"hourlyRate"`
	Status     MembershipStatus `json:"membershipStatus"`
	Type       string           `json:"membershipType"`
	Target     string           `json:"target"`
	UserID     string           `json:"userId"`
}

// MembershipStatus possible Membership Status
type MembershipStatus string

// MembershipStatusPending membership is Pending
const MembershipStatusPending = MembershipStatus("PENDING")

// MembershipStatusActive membership is Active
const MembershipStatusActive = MembershipStatus("ACTIVE")

// MembershipStatusDeclined membership is Declined
const MembershipStatusDeclined = MembershipStatus("DECLINED")

// MembershipStatusInactive membership is Inactive
const MembershipStatusInactive = MembershipStatus("INACTIVE")

// WorkspaceSettings DTO
type WorkspaceSettings struct {
	CanSeeTimeSheet                    bool   `json:"canSeeTimeSheet"`
	DefaultBillableProjects            bool   `json:"defaultBillableProjects"`
	ForceDescription                   bool   `json:"forceDescription"`
	ForceProjects                      bool   `json:"forceProjects"`
	ForceTags                          bool   `json:"forceTags"`
	ForceTasks                         bool   `json:"forceTasks"`
	LockTimeEntries                    string `json:"lockTimeEntries"`
	OnlyAdminsCreateProject            bool   `json:"onlyAdminsCreateProject"`
	OnlyAdminsSeeAllTimeEntries        bool   `json:"onlyAdminsSeeAllTimeEntries"`
	OnlyAdminsSeeBillableRates         bool   `json:"onlyAdminsSeeBillableRates"`
	OnlyAdminsSeeDashboard             bool   `json:"onlyAdminsSeeDashboard"`
	OnlyAdminsSeePublicProjectsEntries bool   `json:"onlyAdminsSeePublicProjectsEntries"`
	ProjectFavorites                   bool   `json:"projectFavorites"`
	ProjectPickerSpecialFilter         bool   `json:"projectPickerSpecialFilter"`
	Round                              Round  `json:"round"`
	TimeRoundingInReports              bool   `json:"timeRoundingInReports"`
}

// Round DTO
type Round struct {
	Minutes string `json:"minutes"`
	Round   string `json:"round"`
}

// HourlyRate DTO
type HourlyRate struct {
	Amount   int32  `json:"amount"`
	Currency string `json:"currency"`
}
