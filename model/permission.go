package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/a5932016/go-ddd-example/util/copy"
	"github.com/pkg/errors"
)

const (
	DomPrefix = "div:"
	ObjPrefix = "obj:"
	ActPrefix = "act:"
)

const (
	PrefixedActionRead   = string(ActPrefix + ActionRead)
	PrefixedActionCreate = string(ActPrefix + ActionCreate)
	PrefixedActionUpdate = string(ActPrefix + ActionUpdate)
	PrefixedActionDelete = string(ActPrefix + ActionDelete)

	// root resources
	PrefixedResourceUser = string(ObjPrefix + ResourceUser)

	// admin resources
)

type Permission struct {
	Name    Resource         `json:"name" binding:"required"`
	Actions []ResourceAction `json:"actions" binding:"required,min=1"`
}

type ResourceAction struct {
	Name        Action `json:"name" binding:"required"`
	Status      bool   `json:"status" binding:"required"`
	IsAvailable bool   `json:"isAvailable" binding:"-"`
}

func NewPermissionsHandler(file string) (PermissionsHandler, error) {
	stuffedPermissionMap, err := readPermissionsFile(file)
	if err != nil {
		return PermissionsHandler{}, err
	}

	ph := PermissionsHandler{
		stuffedPermissionMap: stuffedPermissionMap,
	}

	ph.initAllAllowedPermissions()

	return ph, nil
}

func readPermissionsFile(path string) (map[Resource]map[Action]ResourceAction, error) {
	contentBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "buildPermissionsMap: os.ReadFile")
	}

	permissionsMap := make(map[Resource]map[Action]ResourceAction)
	err = json.Unmarshal(contentBytes, &permissionsMap)
	if err != nil {
		return nil, errors.Wrap(err, "buildPermissionsMap: json.Unmarshal")
	}

	return permissionsMap, nil
}

type PermissionsHandler struct {
	AllAllowedPermissions []Permission
	stuffedPermissionMap  map[Resource]map[Action]ResourceAction
}

func (ph *PermissionsHandler) initAllAllowedPermissions() {
	var permissions []Permission
	for resourceObject, mResourceAction := range ph.stuffedPermissionMap {
		permission := Permission{Name: resourceObject}
		var permissionActions []ResourceAction
		for _, permissionAction := range mResourceAction {
			permissionAction.Status = true
			permissionActions = append(permissionActions, permissionAction)
		}

		permission.Actions = permissionActions
	}

	ph.AllAllowedPermissions = permissions
}

func (ph PermissionsHandler) CasbinPoliciesToPermissions(policies [][]string) (permissions []Permission, err error) {
	// Copy stuffed permission map
	copiedPermissionMap := make(map[Resource]map[Action]ResourceAction)
	if err := copy.DeepCopy(&copiedPermissionMap, &ph.stuffedPermissionMap); err != nil {
		return nil, errors.Wrap(err, "copy.DeepCopy(&copiedPermissionMap, &stuffedPermissionMap)")
	}

	// Enable new policy rules
	for _, p := range policies {
		obj := Resource(strings.TrimPrefix(p[1], ObjPrefix))
		act := Action(strings.TrimPrefix(p[2], ActPrefix))
		permissionAction := copiedPermissionMap[obj][act]
		permissionAction.Status = true
		copiedPermissionMap[obj][act] = permissionAction
	}

	// Build permissions from map
	for resourceObject, mActions := range copiedPermissionMap {
		permission := Permission{Name: resourceObject}
		for _, permissionAction := range mActions {
			permission.Actions = append(permission.Actions, permissionAction)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (ph PermissionsHandler) PermissionsToCasbinPolicies(prefixedDivisionNameId string, permissions []Permission) [][]string {
	var rules [][]string
	for _, perm := range permissions {
		for _, permAction := range perm.Actions {
			if permAction.Status && ph.stuffedPermissionMap[perm.Name][permAction.Name].IsAvailable { // Only add rules for allowed actions
				rule := []string{
					prefixedDivisionNameId,
					perm.Name.Prefix(),
					permAction.Name.Prefix(),
				}

				rules = append(rules, rule)
			}
		}
	}

	return rules
}

type Resource string

func (ro Resource) Prefix() string {
	return fmt.Sprintf("%s%s", ObjPrefix, ro)
}

const (
	ResourceUser Resource = "user"
)

type Action string

func (a Action) Prefix() string {
	return fmt.Sprintf("%s%s", ActPrefix, a)
}

const (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

func GetDefaultDivisionCasbinPolicies(prefixedDivisionNameID string) [][]string {
	return [][]string{
		// {prefixedDivisionNameID, PrefixedResourceUser, PrefixedActionRead},
		// {prefixedDivisionNameID, PrefixedResourceUser, PrefixedActionCreate},
		// {prefixedDivisionNameID, PrefixedResourceUser, PrefixedActionUpdate},
		// {prefixedDivisionNameID, PrefixedResourceUser, PrefixedActionDelete},
	}
}
