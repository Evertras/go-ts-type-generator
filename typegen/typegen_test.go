package typegen

import (
	"strings"
	"testing"
)

type MockStructEmpty struct{}

type MockStructStrings struct {
	SomeField    string `json:"textstuff" tsdesc:"It's a field of some kind."`
	Another      string
	DontLookAtMe string `json:"-"`
}

type MockStructInts struct {
	A int
	B int8
	C int16
	D int32
	E int64
	F uint
	G uint8
	H uint16
	I uint32
	J uint64
	K float32
	L float64
}

type MockStructNestedInner struct {
	NestedField int `json:"x" tsdesc:"A really important value."`
}

type MockStructNestedOuter struct {
	Inner MockStructNestedInner `json:"inner"`
}

type MockStructPointer struct {
	Val *int `json:"x"`
}

type MockStructNestedCircular struct {
	Itself *MockStructNestedCircular `json:"circular,omitempty"`
}

type MockStructNestedBadFieldOuter struct {
	Inner MockStructNestedBadFieldInner
}

type MockStructNestedBadFieldInner struct {
	BadField complex128
}

type AliasedInt int
type MockStructAliased struct {
	X AliasedInt
}

type MockStructExplicitType struct {
	X AliasedInt `json:"x" tstype:"ExplicitType"`
}

type MockStructAnyInterface struct {
	Anything interface{}
}

type MockStructWithUnexportedFields struct {
	PublicStuff string
	secret      string
}

func TestGeneratesBasicInterfacesCorrectly(t *testing.T) {
	tests := []struct {
		Input       interface{}
		Output      string
		ExpectError bool
		Config      Config
	}{
		{
			Input:  MockStructEmpty{},
			Output: "export interface IMockStructEmpty {\n}",
		},
		{
			Input: MockStructStrings{},
			Output: `export interface IMockStructStrings {
	/**
	 * It's a field of some kind.
	 */
	textstuff: string;
	Another: string;
}`,
		},
		{
			Input: MockStructInts{},
			Output: `export interface IMockStructInts {
	A: number;
	B: number;
	C: number;
	D: number;
	E: number;
	F: number;
	G: number;
	H: number;
	I: number;
	J: number;
	K: number;
	L: number;
}`,
		},
		{
			Input: MockStructNestedOuter{},
			Output: `export interface IMockStructNestedOuter {
	inner: IMockStructNestedInner;
}

export interface IMockStructNestedInner {
	/**
	 * A really important value.
	 */
	x: number;
}`,
		},
		{
			Input: MockStructPointer{},
			Output: `export interface IMockStructPointer {
	x: number | null;
}`,
		},
		{
			Input: MockStructNestedCircular{},
			Output: `export interface IMockStructNestedCircular {
	circular: IMockStructNestedCircular | null | undefined;
}`,
		},
		{
			Input:       MockStructNestedBadFieldOuter{},
			ExpectError: true,
		},
		{
			Input: MockStructAliased{},
			Output: `export interface IMockStructAliased {
	X: number;
}`,
		},
		{
			Input: MockStructStrings{},
			Config: Config{
				Indentation: "  ",
			},
			Output: `export interface IMockStructStrings {
  /**
   * It's a field of some kind.
   */
  textstuff: string;
  Another: string;
}`,
		},
		{
			Input: MockStructExplicitType{},
			Output: `export interface IMockStructExplicitType {
	x: ExplicitType;
}`,
		},
		{
			Input: MockStructEmpty{},
			Config: Config{
				Prefix: "Message",
			},
			Output: `export interface IMessageMockStructEmpty {
}`,
		},
		{
			Input: MockStructNestedOuter{},
			Config: Config{
				Prefix: "Message",
			},
			Output: `export interface IMessageMockStructNestedOuter {
	inner: IMessageMockStructNestedInner;
}

export interface IMessageMockStructNestedInner {
	/**
	 * A really important value.
	 */
	x: number;
}`,
		},
		{
			Input: MockStructAnyInterface{},
			Output: `export interface IMockStructAnyInterface {
	Anything: any;
}`,
		},
		{
			Input: MockStructWithUnexportedFields{},
			Output: `export interface IMockStructWithUnexportedFields {
	PublicStuff: string;
}`,
		},
	}

	for _, test := range tests {
		g := NewWithConfig(test.Config)
		builder := strings.Builder{}

		err := g.GenerateSingle(&builder, test.Input)

		if test.ExpectError {
			if err == nil {
				t.Errorf("Expected error but instead got output:\n%s", builder.String())
			}

			// Either way, we're done
			continue
		}

		if err != nil {
			t.Error("failed to generate:", err)
			continue
		}

		str := builder.String()

		if str != test.Output {
			t.Errorf("\n----Expected:\n%s\n----but got:\n%s", test.Output, str)
		}
	}
}

func TestGenerateTypesWorks(t *testing.T) {
	builder := strings.Builder{}
	g := New()

	err := g.GenerateTypes(
		&builder,
		MockStructEmpty{},
		MockStructEmpty{}, // double up to make sure we don't write the same interface twice
		MockStructStrings{},
		MockStructNestedOuter{}, // don't explicitly include the inner
	)

	if err != nil {
		t.Error(err)
	}

	str := builder.String()

	expected := `export interface IMockStructEmpty {
}

export interface IMockStructStrings {
	/**
	 * It's a field of some kind.
	 */
	textstuff: string;
	Another: string;
}

export interface IMockStructNestedOuter {
	inner: IMockStructNestedInner;
}

export interface IMockStructNestedInner {
	/**
	 * A really important value.
	 */
	x: number;
}`

	if str != expected {
		t.Errorf("\n----Expected:\n%s\n----but got:\n%s", expected, str)
	}
}
