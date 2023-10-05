/*
 * flyteidl/service/admin.proto
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package flyteadmin

// Additional launch plan attributes included in the LaunchPlanSpec not strictly required to launch the reference workflow.
type AdminLaunchPlanMetadata struct {
	Schedule      *AdminSchedule      `json:"schedule,omitempty"`
	Notifications []AdminNotification `json:"notifications,omitempty"`
}