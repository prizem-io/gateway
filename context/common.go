package context

import (
	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/identity"
	"github.com/satori/go.uuid"
)

type Common struct {
	//prizem.Context
	DataAccessor
	// Gateway context
	requestID         string
	subjectType       string
	dataAccessor      DataAccessor
	credential        *config.Credential
	identity          identity.Identity
	consumer          *config.Consumer
	plan              *config.Plan
	service           *config.Service
	operation         *config.Operation
	version           string
	claims            identity.Claims
	middlewareHandler MiddlewareHandler
	err               error
}

func (c *Common) Initialize(subjectType string) {
	c.requestID = uuid.NewV4().String()
	c.subjectType = subjectType
	c.claims = identity.Claims{}
}

func (c *Common) Reset() {
	c.dataAccessor = nil
	c.credential = nil
	c.identity = nil
	c.consumer = nil
	c.plan = nil
	c.service = nil
	c.operation = nil
	c.claims = nil
	c.middlewareHandler = nil
	c.err = nil
}

func (c *Common) RequestID() string {
	return c.requestID
}

func (c *Common) SubjectType() string {
	return c.subjectType
}

func (c *Common) GetDataAccessor() DataAccessor {
	return c.DataAccessor
}

func (c *Common) SetDataAccessor(dataAccessor DataAccessor) {
	c.DataAccessor = dataAccessor
}

func (c *Common) Credential() *config.Credential {
	return c.credential
}

func (c *Common) SetCredential(credential *config.Credential) {
	c.credential = credential
}

func (c *Common) Identity() identity.Identity {
	return c.identity
}

func (c *Common) SetIdentity(identity identity.Identity) {
	c.identity = identity
}

func (c *Common) Consumer() *config.Consumer {
	return c.consumer
}

func (c *Common) SetConsumer(consumer *config.Consumer) {
	c.consumer = consumer
}

func (c *Common) Plan() *config.Plan {
	return c.plan
}

func (c *Common) SetPlan(plan *config.Plan) {
	c.plan = plan
}

func (c *Common) Service() *config.Service {
	return c.service
}

func (c *Common) SetService(service *config.Service) {
	c.service = service
}

func (c *Common) Operation() *config.Operation {
	return c.operation
}

func (c *Common) SetOperation(operation *config.Operation) {
	c.operation = operation
}

func (c *Common) Version() string {
	return c.version
}

func (c *Common) SetVersion(version string) {
	c.version = version
}

func (c *Common) Claims() identity.Claims {
	return c.claims
}

func (c *Common) GetError() error {
	return c.err
}

func (c *Common) SetError(err error) {
	c.err = err
	c.middlewareHandler.Stop()
}

func (c *Common) IsStopped() bool {
	return c.middlewareHandler.IsStopped()
}

func (c *Common) Stop() {
	c.middlewareHandler.Stop()
}

func (c *Common) SetMiddlewareHandler(handler MiddlewareHandler) {
	c.middlewareHandler = handler
}

func (c *Common) DoExecute(ctx Context) error {
	err := c.middlewareHandler.Execute(ctx)
	if err != nil {
		c.err = err
		return err
	}

	return nil
}

func (c *Common) DoNext(ctx Context) error {
	err := c.middlewareHandler.Next(ctx)
	if err != nil {
		c.err = err
		return err
	}

	return nil
}
