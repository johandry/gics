package schematics

// WorkspaceStatus is the status of a Schematics workspace
type WorkspaceStatus string

const (
	// WorkspaceStatusNew is the status when the Workspace being created
	WorkspaceStatusNew = WorkspaceStatus("NEW")

	// WorkspaceStatusInactive is the status when it's alredy created waiting to do the planning
	WorkspaceStatusInactive = WorkspaceStatus("INACTIVE")

	// WorkspaceStatusPlaning is the status when the workspace is doing the planning
	WorkspaceStatusPlaning = WorkspaceStatus("PLANING")

	// WorkspaceStatusPlanned is the status when the workspace planing is completed
	WorkspaceStatusPlanned = WorkspaceStatus("PLANNED")

	// WorkspaceStatusApplying is the status when the workspace is applying the changes
	WorkspaceStatusApplying = WorkspaceStatus("APPLYING")

	// WorkspaceStatusDestroyed is the status when the workspace resources were destroyed and the workspace deleted
	WorkspaceStatusDestroyed = WorkspaceStatus("DESTROYED")

	// WorkspaceStatusDeleted is the status when the workspace was deleted
	WorkspaceStatusDeleted = WorkspaceStatus("DELETED")
)

// Status is the status of a Schematics workspace
// type Status int

// const (
// 	StatusNew Status = iota
// 	StatusInactive
// )

// func (s Status) String() string {
// 	return [...]string{
// 		statusNew,
// 		statusInactive,
// 	}[s]
// }
