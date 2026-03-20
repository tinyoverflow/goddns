package factory

import (
	"fmt"
	"goddns/internal/retriever"
	"maps"
)

type RetrieverFactory func(params map[string]any) (retriever.Retriever, error)

var retrieverFactories = map[string]RetrieverFactory{}

func RegisterRetriever(typeName string, f RetrieverFactory) {
	retrieverFactories[typeName] = f
}

func BuildRetriever(pool map[string]map[string]any, ref map[string]any) (retriever.Retriever, error) {
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

	factory, ok := retrieverFactories[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown retriever type %q", typeName)
	}

	params := maps.Clone(def)
	maps.Insert(params, maps.All(ref))

	delete(params, "type")
	delete(params, "name")

	return factory(params)
}
