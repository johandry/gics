package schematics

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const (
	activityNameForCreate  = "WORKSPACE_CREATE"
	activityNameForDestroy = "DESTROY"
)

const (
	listWorkspaceActivitiesTimeout  = 50
	refreshWorkspaceActivityTimeout = 30
)

// NilActivity is an empty or nil activity that doesn't exists or already finished
var NilActivity = NewActivity(nil, "", nil)

// Activity encapsupate the information about a Schematics workspace activity
type Activity struct {
	ID          string    `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string    `json:"name,omitempty" yaml:"name,omitempty"`
	Status      string    `json:"status,omitempty" yaml:"status,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty" yaml:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty" yaml:"end_time,omitempty"`
	PerformedAt time.Time `json:"performed_at,omitempty" yaml:"performed_at,omitempty"`
	PerformedBy string    `json:"performed_by,omitempty" yaml:"performed_by,omitempty"`
	Message     string    `json:"message,omitempty" yaml:"message,omitempty"`
	WorkspaceID string    `json:"workspace_id,omitempty" yaml:"workspace_id,omitempty"`

	service *Service
	logger  *log.Logger
}

// NewActivity creates a new Activity from a WorkspaceActivity
func NewActivity(service *Service, workspaceID string, act *apiv1.WorkspaceActivity) Activity {
	if workspaceID == "" || act == nil {
		return Activity{}
	}

	activity := Activity{
		WorkspaceID: workspaceID,
		service:     service,
	}
	activity.update(act)

	return activity
}

func (a *Activity) update(act *apiv1.WorkspaceActivity) {
	template := apiv1.WorkspaceActivityTemplate{}
	if act.Templates != nil && len(*act.Templates) > 0 {
		template = (*act.Templates)[0]
	}

	a.ID = stringValue(act.ActionId)
	a.PerformedBy = stringValue(act.PerformedBy)
	a.Message = stringValue(template.Message)

	if act.Name != nil {
		a.Name = string(*act.Name)
	}
	if act.Status != nil {
		a.Status = string(*act.Status)
	}
	if template.StartTime != nil {
		a.StartTime = *template.StartTime
	}
	if template.EndTime != nil {
		a.EndTime = *template.EndTime
	}
	if act.PerformedAt != nil {
		a.PerformedAt = *act.PerformedAt
	}
}

func (a *Activity) isNil() bool {
	return len(a.ID) == 0
}

// getActivities gets all the activities in a given workspace
func getActivities(service *Service, workspaceID string) ([]Activity, error) {
	// GetActivities Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), listWorkspaceActivitiesTimeout*time.Second)
	defer cancelFunc()

	params := &apiv1.ListWorkspaceActivitiesParams{}
	resp, err := service.clientWithResponses.ListWorkspaceActivitiesWithResponse(ctx, workspaceID, params)
	if err != nil {
		return nil, err
	}
	if code := resp.StatusCode(); code != 200 {
		return nil, getAPIError("failed to list the workspace activities", resp.Body)
	}
	response := resp.JSON200 // WorkspaceActivities

	// No activities found
	if response.Actions == nil || len(*response.Actions) == 0 {
		// fmt.Printf("[DEBUG] no Activities for Workspace %s\n", workspaceID)
		return []Activity{}, nil
	}

	// fmt.Printf("[DEBUG] all Activities %+v\n", *response.Actions)

	wID := *response.WorkspaceId
	activities := []Activity{}
	for _, act := range *response.Actions {
		activity := NewActivity(service, wID, &act)
		// fmt.Printf("[DEBUG] appending Activity %+v\n", activity)
		activities = append(activities, activity)
	}

	return activities, nil
}

// SetOutput sets a logger for the activity. It won't log by default
func (a *Activity) SetOutput(w io.Writer) *Activity {
	if w == nil {
		return a
	}

	a.logger = log.New(w, fmt.Sprintf("[%s, %s]", a.WorkspaceID, a.ID), log.Ldate|log.Ltime|log.Lshortfile)
	return a
}

func (a *Activity) logPrintf(format string, v ...interface{}) {
	if a.logger == nil {
		// log.Printf(format, v...)
		return
	}

	a.logger.Printf(format, v...)
}

func (a *Activity) refresh() error {
	if a.isNil() {
		return fmt.Errorf("can't refresh a NilActivity")
	}
	if a.service == nil {
		return fmt.Errorf("this Activity don't have a service")
	}
	// refresh Activity Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), refreshWorkspaceActivityTimeout*time.Second)
	defer cancelFunc()

	params := &apiv1.GetWorkspaceActivityParams{}
	resp, err := a.service.clientWithResponses.GetWorkspaceActivityWithResponse(ctx, a.WorkspaceID, a.ID, params)
	if err != nil {
		return err
	}
	if code := resp.StatusCode(); code != 200 {
		return getAPIError("failed to refresh the workspace activity", resp.Body)
	}
	response := resp.JSON200 // WorkspaceActivity

	a.update(response)

	return nil
}

// Wait ...
func (a *Activity) Wait() error {
	if a.isNil() {
		return nil
	}

	for {
		if err := a.refresh(); err != nil {
			return err
		}
		if a.Status == "DONE" {
			a.logPrintf("completed. Status: %s", a.Status)
			return nil
		}
		a.logPrintf("waiting. Status: %s", a.Status)
	}
}

// Error ...
func (a *Activity) Error() error {
	errMsg := ""

	if len(errMsg) == 0 {
		return nil
	}
	return fmt.Errorf(errMsg)
}
