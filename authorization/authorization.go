package authorization

import (
	"strings"
	//ef "github.com/prizem-io/gateway/errorfactory"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
)

type (
	SetExists struct{}
	StringSet map[string]SetExists
)

var exists = SetExists{}

func Handler(ctx context.Context) error {
	consumer := ctx.Consumer()
	identity := ctx.Identity()
	var consumerPermissions map[string]StringSet
	var identityPermissions map[string]StringSet
	var allPermissionIDs StringSet

	if consumer != nil || identity != nil {
		allSize := 0
		if consumer != nil {
			allSize += len(consumer.PermissionIDs)
		}
		if identity != nil {
			allSize += len(identity.PermissionIDs())
		}
		allPermissionIDs = make(StringSet, allSize)

		if consumer != nil {
			permissionIDs := consumer.PermissionIDs
			consumerPermissions := make(map[string]StringSet, len(permissionIDs))
			populatePermissionMap(permissionIDs, consumerPermissions, allPermissionIDs)
		}

		if identity != nil {
			permissionIDs := identity.PermissionIDs()
			identityPermissions := make(map[string]StringSet, len(permissionIDs))
			populatePermissionMap(permissionIDs, identityPermissions, allPermissionIDs)
		}

		claims := ctx.Claims()

		for permissionID, _ := range allPermissionIDs {
			permission, err := ctx.GetPermission(permissionID)
			if err != nil {
				return nil
				/*ef.New(ctx, "unknownPermission", ef.Params{
					"permissionId": permissionId,
				})*/
			}

			// Filter out permissions that have not been granted to
			// consumer, principal, or both
			if (permission.Scope == config.ScopeBoth ||
				permission.Scope == config.ScopeConsumer) &&
				!hasPermissionID(consumerPermissions, permission.ID) {
				continue
			}

			if (permission.Scope == config.ScopeBoth ||
				permission.Scope == config.ScopeUser) &&
				!hasPermissionID(identityPermissions, permission.ID) {
				continue
			}

			if permission.Type == config.TypeEntity {
				consumerActions := consumerPermissions[permission.ID]
				identityActions := identityPermissions[permission.ID]

				intersection := make(StringSet,
					maxInt(len(consumerActions), len(identityActions)))
				for permissionId, _ := range consumerActions {
					if _, ok := identityActions[permissionId]; ok {
						intersection[permissionId] = exists
					}
				}

				if len(intersection) == 0 {
					continue
				}

				if len(intersection) == 1 {
					for action, _ := range intersection {
						claims.Set(permission.ClaimPath, action)
					}
				} else {
					claims.Set(permission.ClaimPath, intersection)
				}
			} else {
				claims.Set(permission.ClaimPath, permission.ClaimValue)
			}
		}
	}

	return nil
}

func maxInt(left, right int) int {
	if left > right {
		return left
	} else {
		return right
	}
}

func hasPermissionID(permissions map[string]StringSet, id string) bool {
	_, ok := permissions[id]
	return ok
}

func populatePermissionMap(permissionIDs []string, permissionMap map[string]StringSet, allPermissionIDs StringSet) {
	for _, permissionId := range permissionIDs {
		index := strings.IndexByte(permissionId, ':')
		if index != -1 {
			id := permissionId[0:index]
			action := permissionId[:index+1]
			set := permissionMap[id]
			if set == nil {
				set = make(StringSet, 5)
				permissionMap[id] = set
				allPermissionIDs[id] = exists
			}
			set[action] = exists
		} else {
			_, ok := permissionMap[permissionId]
			if !ok {
				set := make(StringSet, 0)
				permissionMap[permissionId] = set
				allPermissionIDs[permissionId] = exists
			}
		}
	}
}
