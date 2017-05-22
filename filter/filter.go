package filter

import (
	"fmt"
	"sort"

	"github.com/prizem-io/gateway/backend"
	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
)

type (
	Filter interface {
		Name() string
		Priority() int
		Evaluate(context.Context, interface{}) error
	}

	Execution struct {
		Filter        Filter
		Configuration interface{}
	}
)

var (
	_config   config.Configuration
	filterMap = map[string]Filter{}
)

func Initialize(config config.Configuration) {
	_config = config
}

func Register(filters ...Filter) {
	lookup := make(map[string]Filter, len(filters))
	for _, filter := range filters {
		if initializable, ok := filter.(config.Initializable); ok {
			initializable.Initialize(_config)
		}
		lookup[filter.Name()] = filter
	}
	filterMap = lookup
}

type Invocation struct {
	filter         Filter
	configurations []interface{}
}

type SortedByPriority []Invocation

func (s SortedByPriority) Len() int {
	return len(s)
}
func (s SortedByPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortedByPriority) Less(i, j int) bool {
	return s[i].filter.Priority() < s[j].filter.Priority()
}

func HasFilter(name string) bool {
	_, exists := filterMap[name]
	return exists
}

func GetConfig(configuration *config.PluginConfig) (interface{}, error) {
	filter, ok := filterMap[configuration.Name]
	if !ok {
		return nil, fmt.Errorf("Could not find plugin: %s", configuration.Name)
	}

	configurable, ok := filter.(config.Configurable)
	if !ok {
		return nil, nil
	}

	return configurable.DecodeConfig(configuration.Properties)
}

func Handler(ctx context.Context) error {
	filters := getFilterConfigurations(ctx)

	invocations, err := getFilterInvocations(ctx, filters)
	if err != nil {
		return err
	}

	return invokeFilters(ctx, invocations)
}

func getFilterConfigurations(ctx context.Context) []config.PluginConfig {
	consumer := ctx.Consumer()
	operation := ctx.Operation()
	numFilters := len(ctx.Service().Filters)

	if operation != nil {
		numFilters += len(operation.Filters)
	}
	if consumer != nil {
		numFilters += len(consumer.Filters)
	}

	filters := make([]config.PluginConfig, 0, numFilters)
	if consumer != nil {
		filters = append(filters, consumer.Filters...)
	}
	filters = append(filters, ctx.Service().Filters...)
	if operation != nil {
		filters = append(filters, operation.Filters...)
	}

	return filters
}

func getFilterInvocations(ctx context.Context, filters []config.PluginConfig) (SortedByPriority, error) {
	sorted := make(SortedByPriority, 0, len(filters))
	grouped := make(map[string][]interface{}, len(filters))

	for _, configuration := range filters {
		// First time we are referencing this policy
		// Look it up
		filter, ok := filterMap[configuration.Name]
		if !ok {
			return nil, ef.NewError(ctx, "unregisterdFilter", ef.Params{
				"filter": configuration.Name,
			})
		}

		list, ok := grouped[filter.Name()]
		if !ok {
			// Create slice of configurations
			// Append to grouped configurations and sorted slice
			list = make([]interface{}, 0, 10)
			grouped[filter.Name()] = list

			sorted = append(sorted, Invocation{
				filter: filter,
			})
		}

		list = append(list, configuration.Config)
		grouped[configuration.Name] = list
	}

	for i := 0; i < len(sorted); i++ {
		fi := &sorted[i]
		fi.configurations = grouped[fi.filter.Name()]
	}

	// Sorts by Filter.Priority()
	sort.Sort(sorted)

	return sorted, nil
}

func invokeFilters(ctx context.Context, invocations SortedByPriority) error {
	executions := make([]Execution, 0, len(invocations))

	for _, invocation := range invocations {
		filter := invocation.filter

		configurations := invocation.configurations
		var configuration interface{}

		if combinable, ok := filter.(config.ConfigurationCombiner); ok {
			// Give the policy the chance to combine configuration data
			config, err := combinable.Combine(configurations)
			if err != nil {
				return err
			}
			configuration = config
		} else if len(configurations) > 0 {
			configuration = configurations[0]
		}

		executions = append(executions, Execution{
			Filter:        filter,
			Configuration: configuration,
		})
	}

	handlerName := ctx.Service().Backend.Name
	handler, err := backend.GetHandler(ctx, &handlerName)
	if err != nil {
		return err
	}

	middleware := filterMiddleware{
		filters:        executions,
		backendHandler: handler,
	}

	ctx.SetMiddlewareHandler(&middleware)

	for !ctx.IsStopped() {
		err := ctx.Execute()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
