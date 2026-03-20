// gendoc generates docs/parameters.md by reflecting on registered plugin config structs.
//
// Usage: go run ./cmd/gendoc
package main

import (
	"fmt"
	"goddns/internal/plugin"
	"os"
	"reflect"
	"slices"
	"strings"
)

func main() {
	var b strings.Builder

	b.WriteString("# Plugin Parameters\n\n")
	b.WriteString("This file is auto-generated. Run `go run ./cmd/gendoc` to regenerate.\n\n")

	b.WriteString("## Retrievers\n\n")
	writeSection(&b, plugin.RetrieverConfigTypes())

	b.WriteString("## Providers\n\n")
	writeSection(&b, plugin.ProviderConfigTypes())

	const outPath = "docs/parameters.md"
	if err := os.MkdirAll("docs", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "creating docs directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "writing %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s\n", outPath)
}

func writeSection(b *strings.Builder, types map[string]reflect.Type) {
	names := make([]string, 0, len(types))
	for name := range types {
		names = append(names, name)
	}
	slices.Sort(names)

	for _, name := range names {
		t := types[name]
		fmt.Fprintf(b, "### `%s`\n\n", name)
		b.WriteString("| Parameter | Type | Required | Default | Description |\n")
		b.WriteString("|-----------|------|----------|---------|-------------|\n")

		for i := range t.NumField() {
			f := t.Field(i)
			paramName := f.Tag.Get("json")
			if paramName == "" || paramName == "-" {
				continue
			}
			if idx := strings.Index(paramName, ","); idx != -1 {
				paramName = paramName[:idx]
			}

			required := f.Tag.Get("required")
			if required == "true" {
				required = "Yes"
			} else {
				required = "No"
			}

			defaultVal := f.Tag.Get("default")
			if f.Type.String() == "string" && defaultVal != "" {
				defaultVal = fmt.Sprintf("%q", defaultVal)
			}

			doc := f.Tag.Get("doc")

			fmt.Fprintf(b, "| `%s` | `%s` | %s | %s | %s |\n",
				paramName,
				f.Type.Name(),
				required,
				defaultVal,
				doc,
			)
		}
		b.WriteString("\n")
	}
}
