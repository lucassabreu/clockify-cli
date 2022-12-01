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
	HourlyRate  Rate              `json:"hourlyRate"`
	Memberships []Membership
}

// Membership DTO
type Membership struct {
	HourlyRate *Rate            `json:"hourlyRate"`
	CostRate   *Rate            `json:"costRate"`
	Status     MembershipStatus `json:"membershipStatus"`
	Type       string           `json:"membershipType"`
	TargetID   string           `json:"targetId"`
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
	AdminOnlyPages                     []string      `json:"adminOnlyPages"`
	AutomaticLock                      AutomaticLock `json:"automaticLock"`
	CanSeeTimeSheet                    bool          `json:"canSeeTimeSheet"`
	DefaultBillableProjects            bool          `json:"defaultBillableProjects"`
	ForceDescription                   bool          `json:"forceDescription"`
	ForceProjects                      bool          `json:"forceProjects"`
	ForceTags                          bool          `json:"forceTags"`
	ForceTasks                         bool          `json:"forceTasks"`
	LockTimeEntries                    time.Time     `json:"lockTimeEntries"`
	OnlyAdminsCreateProject            bool          `json:"onlyAdminsCreateProject"`
	OnlyAdminsCreateTag                bool          `json:"onlyAdminsCreateTag"`
	OnlyAdminsCreateTask               bool          `json:"onlyAdminsCreateTask"`
	OnlyAdminsSeeAllTimeEntries        bool          `json:"onlyAdminsSeeAllTimeEntries"`
	OnlyAdminsSeeBillableRates         bool          `json:"onlyAdminsSeeBillableRates"`
	OnlyAdminsSeeDashboard             bool          `json:"onlyAdminsSeeDashboard"`
	OnlyAdminsSeePublicProjectsEntries bool          `json:"onlyAdminsSeePublicProjectsEntries"`
	ProjectFavorites                   bool          `json:"projectFavorites"`
	ProjectGroupingLabel               string        `json:"projectGroupingLabel"`
	ProjectPickerSpecialFilter         bool          `json:"projectPickerSpecialFilter"`
	Round                              Round         `json:"round"`
	TimeRoundingInReports              bool          `json:"timeRoundingInReports"`
	TrackTimeDownToSecond              bool          `json:"trackTimeDownToSecond"`
	IsProjectPublicByDefault           bool          `json:"isProjectPublicByDefault"`
	CanSeeTracker                      bool          `json:"canSeeTracker"`
	FeatureSubscriptionType            string        `json:"featureSubscriptionType"`
}

// AutomaticLock DTO
type AutomaticLock struct {
	ChangeDay       string `json:"changeDay"`
	DayOfMonth      int    `json:"dayOfMonth"`
	FirstDay        string `json:"firstDay"`
	OlderThanPeriod string `json:"olderThanPeriod"`
	OlderThanValue  int    `json:"olderThanValue"`
	Type            string `json:"type"`
}

// Round DTO
type Round struct {
	Minutes string `json:"minutes"`
	Round   string `json:"round"`
}

// Rate DTO
type Rate struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

// TimeEntry DTO
type TimeEntry struct {
	ID            string       `json:"id"`
	Billable      bool         `json:"billable"`
	Description   string       `json:"description"`
	HourlyRate    Rate         `json:"hourlyRate"`
	IsLocked      bool         `json:"isLocked"`
	Project       *Project     `json:"project"`
	ProjectID     string       `json:"projectId"`
	Tags          []Tag        `json:"tags"`
	Task          *Task        `json:"task"`
	TimeInterval  TimeInterval `json:"timeInterval"`
	TotalBillable int64        `json:"totalBillable"`
	User          *User        `json:"user"`
	WorkspaceID   string       `json:"workspaceId"`
}

// TimeInterval DTO
type TimeInterval struct {
	Duration string     `json:"duration"`
	End      *time.Time `json:"end"`
	Start    time.Time  `json:"start"`
}

// Tag DTO
type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
}

func (e Tag) GetID() string   { return e.ID }
func (e Tag) GetName() string { return e.Name }
func (e Tag) String() string  { return e.Name + " (" + e.ID + ")" }

// TaskStatus task status
type TaskStatus string

// TaskStatusActive task is Active
const TaskStatusActive = TaskStatus("ACTIVE")

// TaskStatusDone task is Done
const TaskStatusDone = TaskStatus("DONE")

// Task DTO
type Task struct {
	AssigneeIDs  []string   `json:"assigneeIds"`
	UserGroupIDs []string   `json:"userGroupIds"`
	Estimate     Duration   `json:"estimate"`
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	ProjectID    string     `json:"projectId"`
	Billable     bool       `json:"billable"`
	HourlyRate   *Rate      `json:"hourlyRate"`
	CostRate     *Rate      `json:"costRate"`
	Status       TaskStatus `json:"status"`
	Duration     *Duration  `json:"duration"`
	Favorite     bool       `json:"favorite"`
}

func (e Task) GetID() string   { return e.ID }
func (e Task) GetName() string { return e.Name }

// Client DTO
type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

func (e Client) GetID() string   { return e.ID }
func (e Client) GetName() string { return e.Name }

// CustomField DTO
type CustomField struct {
	CustomFieldID string `json:"customFieldId"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Value         string `json:"value"`
}

// Project DTO
type Project struct {
	WorkspaceID string `json:"workspaceId"`

	ID    string `json:"id"`
	Name  string `json:"name"`
	Note  string `json:"note"`
	Color string `json:"color"`

	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`

	HourlyRate Rate  `json:"hourlyRate"`
	CostRate   *Rate `json:"costRate"`
	Billable   bool  `json:"billable"`

	TimeEstimate   TimeEstimate `json:"timeEstimate"`
	BudgetEstimate BaseEstimate `json:"budgetEstimate"`
	Duration       *Duration    `json:"duration"`

	Archived bool `json:"archived"`
	Template bool `json:"template"`
	Public   bool `json:"public"`
	Favorite bool `json:"favorite"`

	Memberships []Membership `json:"memberships"`

	// Hydrated indicates if the attributes CustomFields and Tasks are filled
	Hydrated     bool          `json:"-"`
	CustomFields []CustomField `json:"customFields,omitempty"`
	Tasks        []Task        `json:"tasks,omitempty"`
}

func (p Project) GetID() string   { return p.ID }
func (p Project) GetName() string { return p.Name }

// EstimateType possible Estimate types
type EstimateType string

// EstimateTypeAuto estimate is Auto
const EstimateTypeAuto = EstimateType("AUTO")

// EstimateTypeManual estimate is Manual
const EstimateTypeManual = EstimateType("MANUAL")

// EstimateResetOption possible Estimate Reset Options
type EstimateResetOption string

// EstimateResetOptionMonthly estimate is Auto
const EstimateResetOptionMonthly = EstimateResetOption("MONTHLY")

// BaseEstimate DTO
type BaseEstimate struct {
	Type         EstimateType         `json:"type"`
	Active       bool                 `json:"active"`
	ResetOptions *EstimateResetOption `json:"resetOptions"`
}

// TimeEstimate DTO
type TimeEstimate struct {
	BaseEstimate
	Estimate           Duration `json:"estimate"`
	IncludeNonBillable bool     `json:"includeNonBillable"`
}

// BudgetEstimate DTO
type BudgetEstimate struct {
	BaseEstimate
	Estimate uint `json:"estimate"`
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
	Roles            *[]Role      `json:"roles"`
}

func (e User) GetID() string   { return e.ID }
func (e User) GetName() string { return e.Name }

// Role DTO
type Role struct {
	Role     string       `json:"role"`
	Entities []RoleEntity `json:"entities"`
}

type RoleEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// TimeEntriesList DTO
type TimeEntriesList struct {
	AllEntriesCount int64           `json:"allEntriesCount"`
	GotAllEntries   bool            `json:"gotAllEntries"`
	TimeEntriesList []TimeEntryImpl `json:"timeEntriesList"`
}

// TimeEntryImpl DTO
type TimeEntryImpl struct {
	Billable     bool         `json:"billable"`
	Description  string       `json:"description"`
	ID           string       `json:"id"`
	IsLocked     bool         `json:"isLocked"`
	ProjectID    string       `json:"projectId"`
	TagIDs       []string     `json:"tagIds"`
	TaskID       string       `json:"taskId"`
	TimeInterval TimeInterval `json:"timeInterval"`
	UserID       string       `json:"userId"`
	WorkspaceID  string       `json:"workspaceId"`
}
