package identity

type (
	Claims map[string]interface{}

	Identity interface {
		Id() string
		Name() string
		Context() *string
		PermissionIDs() []string
		Claims() Claims
	}
)

func (cm Claims) Set(path []string, value interface{}) {
	m := map[string]interface{}(cm)
	for i := 0; i < len(path)-1; i++ {
		key := path[i]
		var nx map[string]interface{}
		next, ok := m[key]
		if ok {
			nx, ok = next.(map[string]interface{})
		}
		if !ok {
			nx = map[string]interface{}{}
			m[key] = nx
			m = nx
		}
	}

	m[path[len(path)-1]] = value
}
