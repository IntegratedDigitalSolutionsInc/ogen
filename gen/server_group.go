package gen

import (
	"strings"

	"github.com/ogen-go/ogen/gen/ir"
)

// ServerGroup represents a group of operations that will be generated as a separate server.
type ServerGroup struct {
	Name       string // e.g., "Images"
	ServerName string // e.g., "ImagesServer"
	Operations []*ir.Operation
	Router     Router // Filtered router for this group
}

// BuildServerGroups creates ServerGroup structures from operation groups when
// the ServerPerOperationGroup feature is enabled.
func BuildServerGroups(groups []*ir.OperationGroup, router Router) []ServerGroup {
	var serverGroups []ServerGroup

	for _, group := range groups {
		// Filter router to only include operations from this group
		filteredRouter := Router{
			Tree:               RouteTree{},
			MaxParametersCount: router.MaxParametersCount,
		}

		// Add routes for this group's operations
		for _, op := range group.Operations {
			route := Route{
				Method:    strings.ToUpper(op.Spec.HTTPMethod),
				Path:      op.Spec.Path.String(),
				Operation: op,
			}
			_ = filteredRouter.Add(route) // Ignore error as routes were already validated
		}

		serverGroups = append(serverGroups, ServerGroup{
			Name:       group.Name,
			ServerName: group.Name + "Server",
			Operations: group.Operations,
			Router:     filteredRouter,
		})
	}

	return serverGroups
}
