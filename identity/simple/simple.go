package simple

import (
	"github.com/prizem-io/gateway/identity"
)

type (
	SimpleIdentity struct {
		id            string
		name          string
		context       *string
		permissionIDs []string
		claims        identity.Claims
	}
)

func New(subject string) (identity.Identity, error) {
	return &SimpleIdentity{
		id:            subject,
		name:          subject,
		context:       &subject,
		permissionIDs: nil,
		claims:        nil,
	}, nil
}

func (i *SimpleIdentity) ID() string {
	return i.id
}

func (i *SimpleIdentity) Name() string {
	return i.name
}

func (i *SimpleIdentity) Context() *string {
	return i.context
}

func (i *SimpleIdentity) PermissionIDs() []string {
	return i.permissionIDs
}

func (i *SimpleIdentity) Claims() identity.Claims {
	return i.claims
}
