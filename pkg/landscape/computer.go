package landscape

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// ComputerService provides methods for the /computers endpoints.
// Access it through Client.Computer rather than constructing it directly.
type ComputerService struct {
	Service
}

// DistributionInfo describes the OS release a computer is running.
type DistributionInfo struct {
	CodeName    *string `json:"code_name,omitempty" yaml:"code_name,omitempty"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	Distributor *string `json:"distributor,omitempty" yaml:"distributor,omitempty"`
	Release     *string `json:"release,omitempty" yaml:"release,omitempty"`
}

// ComputerProfile is a profile (package, repository, security, etc.)
// associated with a computer.
type ComputerProfile struct {
	ID    int     `json:"id" yaml:"id"`
	Name  *string `json:"name,omitempty" yaml:"name,omitempty"`
	Title *string `json:"title,omitempty" yaml:"name,omitempty"`
	Type  *string `json:"type,omitempty" yaml:"name,omitempty"`
}

// Computer is a single computer/instance managed by Landscape. Some fields
// are only populated when the corresponding "with_*" option is requested,
// and a few (vm_info, container_info, ubuntu_pro_info, cloud_init, ...)
// vary in shape enough that they're left as `any` rather than fully typed.
type Computer struct {
	ID                    int               `json:"id" yaml:"id"`
	Title                 *string           `json:"title,omitempty" yaml:"title,omitempty"`
	Comment               *string           `json:"comment,omitempty" yaml:"comment,omitempty"`
	Hostname              *string           `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	TotalMemory           int               `json:"total_memory,omitempty" yaml:"total_memory,omitempty"`
	TotalSwap             int               `json:"total_swap,omitempty" yaml:"total_swap,omitempty"`
	RebootRequiredFlag    bool              `json:"reboot_required_flag,omitempty" yaml:"reboot_required_flag,omitempty"`
	UpdateManagerPrompt   *string           `json:"update_manager_prompt,omitempty" yaml:"update_manager_prompt,omitempty"`
	CloneID               int               `json:"clone_id" yaml:"clone_id"`
	LastExchangeTime      *time.Time        `json:"last_exchange_time,omitempty" yaml:"last_exchange_time,omitempty"`
	LastPingTime          *time.Time        `json:"last_ping_time,omitempty" yaml:"last_ping_time,omitempty"`
	RegisteredAt          *time.Time        `json:"registered_at,omitempty" yaml:"registered_at,omitempty"`
	Tags                  []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	AccessGroup           *string           `json:"access_group,omitempty" yaml:"access_group,omitempty"`
	Distribution          *string           `json:"distribution,omitempty" yaml:"distribution,omitempty"`
	DistributionInfo      *DistributionInfo `json:"distribution_info,omitempty" yaml:"distribution_info,omitempty"`
	CloudInstanceMetadata map[string]any    `json:"cloud_instance_metadata,omitempty" yaml:"cloud_instance_metadata,omitempty"`
	VMInfo                any               `json:"vm_info,omitempty" yaml:"vm_info,omitempty"`
	ContainerInfo         any               `json:"container_info,omitempty" yaml:"container_info,omitempty"`
	UbuntuProInfo         any               `json:"ubuntu_pro_info,omitempty" yaml:"ubuntu_pro_info,omitempty"`
	IsWSLInstance         bool              `json:"is_wsl_instance,omitempty" yaml:"is_wsl_instance,omitempty"`
	Children              []Computer        `json:"children,omitempty" yaml:"children,omitempty"`
	Parent                *Computer         `json:"parent,omitempty" yaml:"parent,omitempty"`
	Archived              bool              `json:"archived,omitempty" yaml:"archived,omitempty"`
	NumChild              int               `json:"num_child,omitempty" yaml:"num_child,omitempty"`
	Annotations           map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Profiles              []ComputerProfile `json:"profiles,omitempty" yaml:"profiles,omitempty"`
	CloudInit             map[string]any    `json:"cloud_init,omitempty" yaml:"cloud_init,omitempty"`
}

// ComputerListOptions holds the query parameters for ComputerService.List.
// Zero-valued fields are simply omitted from the request, so callers only
// need to set what they care about. See the Landscape docs for the full
// set of "query" selector prefixes (tag:, hostname:, distribution:, ...).
type ComputerListOptions struct {
	Query               string
	Limit               int
	Offset              int
	WithAlerts          bool
	WithUpgrades        bool
	WithReleaseUpgrades bool
	WithRebootPackages  bool
	WithNetwork         bool
	WithAllNetwork      bool
	WithHardware        bool
	WithAnnotations     bool
	WithGroupedHardware bool
	WithWSLProfiles     bool
	ArchivedOnly        bool
	// RootOnly defaults to true server-side; leave nil to use that
	// default, or set explicitly (e.g. to false) to override it.
	RootOnly    *bool
	WSLParents  bool
	WSLChildren bool
}

func (o ComputerListOptions) queryParams() map[string]string {
	params := map[string]string{}
	if o.Query != "" {
		params["query"] = o.Query
	}
	if o.Limit > 0 {
		params["limit"] = strconv.Itoa(o.Limit)
	}
	if o.Offset > 0 {
		params["offset"] = strconv.Itoa(o.Offset)
	}
	setTrue := func(key string, v bool) {
		if v {
			params[key] = "true"
		}
	}
	setTrue("with_alerts", o.WithAlerts)
	setTrue("with_upgrades", o.WithUpgrades)
	setTrue("with_release_upgrades", o.WithReleaseUpgrades)
	setTrue("with_reboot_packages", o.WithRebootPackages)
	setTrue("with_network", o.WithNetwork)
	setTrue("with_all_network", o.WithAllNetwork)
	setTrue("with_hardware", o.WithHardware)
	setTrue("with_annotations", o.WithAnnotations)
	setTrue("with_grouped_hardware", o.WithGroupedHardware)
	setTrue("with_wsl_profiles", o.WithWSLProfiles)
	setTrue("archived_only", o.ArchivedOnly)
	setTrue("wsl_parents", o.WSLParents)
	setTrue("wsl_children", o.WSLChildren)
	if o.RootOnly != nil {
		params["root_only"] = strconv.FormatBool(*o.RootOnly)
	}
	return params
}

// ComputerListResponse is the paginated response from GET /computers.
type ComputerListResponse struct {
	Count    int        `json:"count"`
	Results  []Computer `json:"results"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
}

// List returns computers associated with the account, optionally filtered
// and expanded via opts. Corresponds to GET /computers.
func (s *ComputerService) List(ctx context.Context, opts ComputerListOptions) (*ComputerListResponse, error) {
	var result ComputerListResponse
	var failure Error

	resp, err := s.client.R().
		SetContext(ctx).
		SetQueryParams(opts.queryParams()).
		SetResult(&result).
		SetResultError(&failure).
		Get("/computers")
	if err != nil {
		return nil, fmt.Errorf("list computers failed: %w", err)
	}
	if resp.IsStatusFailure() {
		if failure.Message != "" {
			return nil, &failure
		}
		return nil, fmt.Errorf("list computers failed: status %d: %s", resp.StatusCode(), resp.String())
	}
	return &result, nil
}

// ComputerReadOptions holds the query parameters for ComputerService.Read.
type ComputerReadOptions struct {
	WithAnnotations     bool
	WithGroupedHardware bool
	WithHardware        bool
	WithNetwork         bool
	WithProfiles        bool
	WithAlerts          bool
}

func (o ComputerReadOptions) queryParams() map[string]string {
	params := map[string]string{}
	setTrue := func(key string, v bool) {
		if v {
			params[key] = "true"
		}
	}
	setTrue("with_annotations", o.WithAnnotations)
	setTrue("with_grouped_hardware", o.WithGroupedHardware)
	setTrue("with_hardware", o.WithHardware)
	setTrue("with_network", o.WithNetwork)
	setTrue("with_profiles", o.WithProfiles)
	setTrue("with_alerts", o.WithAlerts)
	return params
}

// Read returns a single computer by ID. Corresponds to
// GET /computers/<int:computer_id>.
func (s *ComputerService) Read(ctx context.Context, id int, opts ComputerReadOptions) (*Computer, error) {
	var result Computer
	var failure Error

	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		SetQueryParams(opts.queryParams()).
		SetResult(&result).
		SetResultError(&failure).
		Get("/computers/{id}")
	if err != nil {
		return nil, fmt.Errorf("read computer %d failed: %w", id, err)
	}
	if resp.IsStatusFailure() {
		if failure.Message != "" {
			return nil, &failure
		}
		return nil, fmt.Errorf("read computer %d failed: status %d: %s", id, resp.StatusCode(), resp.String())
	}
	return &result, nil
}

// computersDeleteRequest is the body for POST /computers:delete.
type computersDeleteRequest struct {
	ComputerIDs []int `json:"computer_ids,omitempty" yaml:"computer_ids,omitempty"`
}

// Delete removes the given computers. Corresponds to POST /computers:delete.
//
// Note: this endpoint is available from Landscape Server 26.04 LTS onwards
// for select accounts; older servers will reject it.
func (s *ComputerService) Delete(ctx context.Context, computerIDs ...int) error {
	if len(computerIDs) == 0 {
		return fmt.Errorf("delete computers: at least one computer ID is required")
	}

	var failure Error

	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(computersDeleteRequest{ComputerIDs: computerIDs}).
		SetResultError(&failure).
		Post("/computers:delete")
	if err != nil {
		return fmt.Errorf("delete computers failed: %w", err)
	}
	if resp.IsStatusFailure() {
		if failure.Message != "" {
			return &failure
		}
		return fmt.Errorf("delete computers failed: status %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
}
