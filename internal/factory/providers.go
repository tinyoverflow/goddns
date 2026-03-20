package factory

import (
	"fmt"
	"goddns/internal/provider"
	"maps"
)

type ProviderFactory func(params map[string]any) (provider.Provider, error)

var providerFactories = map[string]ProviderFactory{}

func RegisterProvider(typeName string, f ProviderFactory) {
	providerFactories[typeName] = f
}

func BuildProviders(pool map[string]map[string]any, refs []map[string]any) ([]provider.Provider, error) {
	providers := make([]provider.Provider, 0, len(refs))

	for _, ref := range refs {
		prv, err := buildProvider(pool, ref)

		if err != nil {
			return nil, err
		}

		providers = append(providers, prv)
	}

	return providers, nil
}

func buildProvider(pool map[string]map[string]any, ref map[string]any) (provider.Provider, error) {
	name, _ := ref["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("provider name is required")
	}

	def, ok := pool[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not found", name)
	}

	typeName, _ := def["type"].(string)
	if typeName == "" {
		return nil, fmt.Errorf("type is required for provider %q", name)
	}

	factory, ok := providerFactories[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown provider type %q", typeName)
	}

	params := maps.Clone(def)
	maps.Insert(params, maps.All(ref))

	delete(params, "type")
	delete(params, "name")

	return factory(params)
}
