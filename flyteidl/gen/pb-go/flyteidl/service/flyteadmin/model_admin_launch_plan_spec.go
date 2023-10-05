/*
 * flyteidl/service/admin.proto
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package flyteadmin

// User-provided launch plan definition and configuration values.
type AdminLaunchPlanSpec struct {
	WorkflowId     *CoreIdentifier          `json:"workflow_id,omitempty"`
	EntityMetadata *AdminLaunchPlanMetadata `json:"entity_metadata,omitempty"`
	// Input values to be passed for the execution. These can be overridden when an execution is created with this launch plan.
	DefaultInputs *CoreParameterMap `json:"default_inputs,omitempty"`
	// Fixed, non-overridable inputs for the Launch Plan. These can not be overridden when an execution is created with this launch plan.
	FixedInputs *CoreLiteralMap `json:"fixed_inputs,omitempty"`
	Role        string          `json:"role,omitempty"`
	// Custom labels to be applied to the execution resource.
	Labels *AdminLabels `json:"labels,omitempty"`
	// Custom annotations to be applied to the execution resource.
	Annotations *AdminAnnotations `json:"annotations,omitempty"`
	// Indicates the permission associated with workflow executions triggered with this launch plan.
	Auth            *AdminAuth           `json:"auth,omitempty"`
	AuthRole        *AdminAuthRole       `json:"auth_role,omitempty"`
	SecurityContext *CoreSecurityContext `json:"security_context,omitempty"`
	// Indicates the runtime priority of the execution.
	QualityOfService *CoreQualityOfService `json:"quality_of_service,omitempty"`
	// Encapsulates user settings pertaining to offloaded data (i.e. Blobs, Schema, query data, etc.).
	RawOutputDataConfig *AdminRawOutputDataConfig `json:"raw_output_data_config,omitempty"`
	// Controls the maximum number of tasknodes that can be run in parallel for the entire workflow. This is useful to achieve fairness. Note: MapTasks are regarded as one unit, and parallelism/concurrency of MapTasks is independent from this.
	MaxParallelism int32 `json:"max_parallelism,omitempty"`
	// Allows for the interruptible flag of a workflow to be overwritten for a single execution. Omitting this field uses the workflow's value as a default. As we need to distinguish between the field not being provided and its default value false, we have to use a wrapper around the bool field.
	Interruptible bool `json:"interruptible,omitempty"`
	// Allows for all cached values of a workflow and its tasks to be overwritten for a single execution. If enabled, all calculations are performed even if cached results would be available, overwriting the stored data once execution finishes successfully.
	OverwriteCache bool `json:"overwrite_cache,omitempty"`
	// Environment variables to be set for the execution.
	Envs *AdminEnvs `json:"envs,omitempty"`
}