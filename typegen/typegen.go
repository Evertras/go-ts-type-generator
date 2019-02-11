package typegen

import (
	"reflect"
	"strings"
	"text/template"
)

// Generator can inspect a struct type and generate Typescript definitions.
type Generator struct {
}

const interfaceTemplate = `interface I<<.Name>> {<<range .Fields>>
	<<.Name>>: <<.TypescriptType>>;<<end>>
}`

type fieldTemplateData struct {
	Name           string
	TypescriptType string
}

type typeTemplateData struct {
	Name string

	Fields []fieldTemplateData
}

var templateInterface *template.Template

var typeMapping map[string]string

func init() {
	templateInterface = template.Must(template.New("interface").Delims("<<", ">>").Parse(interfaceTemplate))
	typeMapping = map[string]string{
		"string": "string",
		"int":    "number",
		"int8":   "number",
		"int16":  "number",
		"int32":  "number",
		"int64":  "number",
		"uint":   "number",
		"uint8":  "number",
		"uint16": "number",
		"uint32": "number",
		"uint64": "number",
	}
}

// GenerateSingle takes in a single type and returns a full Typescript definition
// based on that type.
func (g *Generator) GenerateSingle(t interface{}) (string, error) {

	r := reflect.TypeOf(t)

	data := typeTemplateData{
		Name: r.Name(),
	}

	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		fieldName := field.Tag.Get("json")

		// Skip if it's explicitly set to -
		if fieldName == "-" {
			continue
		}

		if fieldName == "" {
			fieldName = field.Name
		}

		var typeName string
		var ok bool

		if typeName, ok = typeMapping[field.Type.Name()]; !ok {
			// For now, if we don't know it, just set it to any
			typeName = "any"
		}

		data.Fields = append(data.Fields, fieldTemplateData{
			Name:           fieldName,
			TypescriptType: typeName,
		})
	}

	var builder strings.Builder

	if err := templateInterface.Execute(&builder, data); err != nil {
		return "", err
	}

	return builder.String(), nil
}
