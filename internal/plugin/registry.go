package plugin

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
)

type providerFactory func(params map[string]any) (Provider, error)
type retrieverFactory func(params map[string]any) (Retriever, error)

type providerEntry struct {
	factory    providerFactory
	configType reflect.Type
}

type retrieverEntry struct {
	factory    retrieverFactory
	configType reflect.Type
}

var providerRegistry = map[string]providerEntry{}
var retrieverRegistry = map[string]retrieverEntry{}

func RegisterProvider(name string, f providerFactory, configExample any) {
	providerRegistry[name] = providerEntry{
		factory:    f,
		configType: reflect.TypeOf(configExample),
	}
}

func RegisterRetriever(name string, f retrieverFactory, configExample any) {
	retrieverRegistry[name] = retrieverEntry{
		factory:    f,
		configType: reflect.TypeOf(configExample),
	}
}

// ProviderConfigTypes returns the registered config struct type for each provider name.
// Used by cmd/gendoc to generate parameter documentation.
func ProviderConfigTypes() map[string]reflect.Type {
	types := make(map[string]reflect.Type, len(providerRegistry))
	for name, entry := range providerRegistry {
		types[name] = entry.configType
	}
	return types
}

// RetrieverConfigTypes returns the registered config struct type for each retriever name.
// Used by cmd/gendoc to generate parameter documentation.
func RetrieverConfigTypes() map[string]reflect.Type {
	types := make(map[string]reflect.Type, len(retrieverRegistry))
	for name, entry := range retrieverRegistry {
		types[name] = entry.configType
	}
	return types
}

func BuildProviders(pool map[string]map[string]any, refs []map[string]any) ([]Provider, error) {
	providers := make([]Provider, 0, len(refs))
	for _, ref := range refs {
		prv, err := buildProvider(pool, ref)
		if err != nil {
			return nil, err
		}
		providers = append(providers, prv)
	}
	return providers, nil
}

func buildProvider(pool map[string]map[string]any, ref map[string]any) (Provider, error) {
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

	entry, ok := providerRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown provider type %q", typeName)
	}

	params := maps.Clone(def)
	maps.Insert(params, maps.All(ref))
	delete(params, "type")
	delete(params, "name")

	return entry.factory(params)
}

func BuildRetriever(pool map[string]map[string]any, ref map[string]any) (Retriever, error) {
	name, _ := ref["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("retriever: name is required")
	}

	def, ok := pool[name]
	if !ok {
		return nil, fmt.Errorf("retriever %q not found", name)
	}

	typeName, _ := def["type"].(string)
	if typeName == "" {
		return nil, fmt.Errorf("retriever %q: type is required", name)
	}

	entry, ok := retrieverRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown retriever type %q", typeName)
	}

	params := maps.Clone(def)
	maps.Insert(params, maps.All(ref))
	delete(params, "type")
	delete(params, "name")

	return entry.factory(params)
}

func Decode[T any](params map[string]any) (T, error) {
	var result T
	b, err := json.Marshal(params)
	if err != nil {
		return result, err
	}
	return result, json.Unmarshal(b, &result)
}
