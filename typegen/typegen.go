package typegen

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"text/template"
)

// Generator can inspect a struct type and generate Typescript definitions.
type Generator struct {
}

// Using custom delimiters here to avoid {} collisions
const interfaceTemplate = `interface I<<.Name>> {<<range .Fields>><<if .Desc>>
	/**
	 * <<.Desc>>
	 */<<end>>
	<<.Name>>: <<.TypescriptType>>;<<end>>
}`

type fieldTemplateData struct {
	Name           string
	TypescriptType string
	Desc           string
}

type typeTemplateData struct {
	Name string

	Fields []fieldTemplateData
}

var templateInterface *template.Template
var typeMapping map[string]string
var typesDefined = make(map[string]bool)

func init() {
	templateInterface = template.Must(template.New("interface").Delims("<<", ">>").Parse(interfaceTemplate))

	// Go type ==> Typescript type
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
func (g *Generator) GenerateSingle(out io.Writer, t interface{}) error {
	r := reflect.TypeOf(t)

	return g.generateSingle(out, r)
}

func (g *Generator) generateSingle(out io.Writer, t reflect.Type) error {
	// Already defined earlier, don't redefine
	if typesDefined[t.Name()] {
		return nil
	}

	typesDefined[t.Name()] = true

	data := typeTemplateData{
		Name: t.Name(),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("json")
		canBeUndefined := false
		canBeNull := false

		// Skip if it's explicitly set to -
		if fieldName == "-" {
			continue
		}

		if fieldName == "" {
			fieldName = field.Name
		} else {
			split := strings.Split(fieldName, ",")

			for _, s := range split {
				if s == "omitempty" {
					canBeUndefined = true
				}
			}

			fieldName = split[0]
		}

		var fieldTypescriptType string
		var ok bool
		var fieldType = field.Type

		kind := fieldType.Kind()

		if kind == reflect.Ptr {
			canBeNull = true
			fieldType = fieldType.Elem()
			kind = fieldType.Kind()
		}

		if kind == reflect.Struct {
			fieldTypescriptType = "I" + fieldType.Name()

			// After we're done, make sure to include this type recursively
			defer func() {
				if !typesDefined[fieldType.Name()] {
					out.Write([]byte("\n\n"))
					g.generateSingle(out, fieldType)
				}
			}()
		} else if fieldTypescriptType, ok = typeMapping[fieldType.Name()]; !ok {
			return errors.New("cannot map typescript type from " + fieldType.Name())
		}

		if canBeNull {
			fieldTypescriptType = fieldTypescriptType + " | null"
		}

		if canBeUndefined {
			fieldTypescriptType = fieldTypescriptType + " | undefined"
		}

		data.Fields = append(data.Fields, fieldTemplateData{
			Name:           fieldName,
			TypescriptType: fieldTypescriptType,
			Desc:           field.Tag.Get("tsdesc"),
		})
	}

	return templateInterface.Execute(out, data)
}
