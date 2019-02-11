package typegen

import (
	"errors"
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

		var fieldTypescriptType string
		var ok bool

		if field.Type.Kind() == reflect.Struct {
			// Going to assume that the inner type is also exported
			fieldTypescriptType = "I" + field.Type.Name()
		} else if fieldTypescriptType, ok = typeMapping[field.Type.Name()]; !ok {
			return "", errors.New("cannot map typescript type from " + field.Type.Name())
		}

		data.Fields = append(data.Fields, fieldTemplateData{
			Name:           fieldName,
			TypescriptType: fieldTypescriptType,
		})
	}

	var builder strings.Builder

	if err := templateInterface.Execute(&builder, data); err != nil {
		return "", err
	}

	return builder.String(), nil
}
