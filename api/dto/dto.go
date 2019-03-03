package dto

import (
	"fmt"
	"time"
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

// UserStatus possible user status
type UserStatus string

// UserStatusActive when the user is Active
const UserStatusActive = UserStatus("ACTIVE")

// UserStatusPendingEmailVerification when the user is Pending Email Verification
const UserStatusPendingEmailVerification = UserStatus("PENDING_EMAIL_VERIFICATION")

// UserStatusDeleted when the user is Deleted
const UserStatusDeleted = UserStatus("DELETED")

// User DTO
type User struct {
	ID               string       `json:"id"`
	ActiveWorkspace  string       `json:"activeWorkspace"`
	DefaultWorkspace string       `json:"defaultWorkspace"`
	Email            string       `json:"email"`
	Memberships      []Membership `json:"memberships"`
	Name             string       `json:"name"`
	ProfilePicture   string       `json:"profilePicture"`
	Settings         UserSettings `json:"settings"`
	Status           UserStatus   `json:"status"`
}

// WeekStart when the week starts
type WeekStart string

// WeekStartMonday when start at Monday
const WeekStartMonday = WeekStart("MONDAY")

// WeekStartTuesday when start at Tuesday
const WeekStartTuesday = WeekStart("TUESDAY")

// WeekStartWednesday when start at Wednesday
const WeekStartWednesday = WeekStart("WEDNESDAY")

// WeekStartThursday when start at Thursday
const WeekStartThursday = WeekStart("THURSDAY")

// WeekStartFriday when start at Friday
const WeekStartFriday = WeekStart("FRIDAY")

// WeekStartSaturday when start at Saturday
const WeekStartSaturday = WeekStart("SATURDAY")

// WeekStartSunday when start at Sunday
const WeekStartSunday = WeekStart("SUNDAY")

// UserSettings DTO
type UserSettings struct {
	DateFormat            string                `json:"dateFormat"`
	IsCompactViewOn       bool                  `json:"isCompactViewOn"`
	LongRunning           bool                  `json:"longRunning"`
	SendNewsletter        bool                  `json:"sendNewsletter"`
	SummaryReportSettings SummaryReportSettings `json:"summaryReportSettings"`
	TimeFormat            string                `json:"timeFormat"`
	TimeTrackingManual    bool                  `json:"timeTrackingManual"`
	TimeZone              string                `json:"timeZone"`
	WeekStart             string                `json:"weekStart"`
	WeeklyUpdates         bool                  `json:"weeklyUpdates"`
}

// SummaryReportSettings DTO
type SummaryReportSettings struct {
	Group    string `json:"group"`
	Subgroup string `json:"subgroup"`
}

// InvitedUser DTO
type InvitedUser struct {
	ID          string       `json:"id"`
	Email       string       `json:"email"`
	Invitation  Invitation   `json:"invitation"`
	Memberships []Membership `json:"memberships"`
}

// Invitation DTO
type Invitation struct {
	Creation       time.Time  `json:"creation"`
	InvitationCode string     `json:"invitationCode"`
	Membership     Membership `json:"membership"`
	WorkspaceID    string     `json:"workspaceId"`
	WorkspaceName  string     `json:"workspaceName"`
}
