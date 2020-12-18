package schematics

import (
	"context"
	"fmt"
	"io"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const (
	activityNameForCreate = "WORKSPACE_CREATE"
)

const (
	listWorkspaceActivitiesTimeout = 50
)

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
		return nil, fmt.Errorf(`{"status_code": %d, "status": %q}`, code, resp.Status())
	}
	response := resp.JSON200 // WorkspaceActivities

	// No activities found
	if response.Actions == nil || len(*response.Actions) == 0 {
		return []Activity{}, nil
	}

	wID := *response.WorkspaceId
	activities := []Activity{}
	for _, act := range *response.Actions {
		template := apiv1.WorkspaceActivityTemplate{}
		if act.Templates != nil && len(*act.Templates) > 0 {
			template = (*act.Templates)[0]
		}
		activity := Activity{
			ID:          stringValue(act.ActionId),
			PerformedBy: stringValue(act.PerformedBy),
			Message:     stringValue(template.Message),
			WorkspaceID: wID,
		}
		if act.Name != nil {
			activity.Name = string(*act.Name)
		}
		if act.Status != nil {
			activity.Status = string(*act.Status)
		}
		if template.StartTime != nil {
			activity.StartTime = *template.StartTime
		}
		if template.EndTime != nil {
			activity.EndTime = *template.EndTime
		}
		if act.PerformedAt != nil {
			activity.PerformedAt = *act.PerformedAt
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// WriteLog ...
func (a *Activity) WriteLog(w io.Writer) *Activity {
	return a
}

// Wait ...
func (a *Activity) Wait() error {
	return nil
}

// Error ...
func (a *Activity) Error() error {
	errMsg := ""

	if len(errMsg) == 0 {
		return nil
	}
	return fmt.Errorf(errMsg)
}
