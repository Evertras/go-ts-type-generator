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

func TestGeneratesBasicInterfacesCorrectly(t *testing.T) {
	tests := []struct {
		Input  interface{}
		Output string
	}{
		{
			Input:  MockStructEmpty{},
			Output: "interface IMockStructEmpty {\n}",
		},
		{
			Input: MockStructStrings{},
			Output: `interface IMockStructStrings {
	/**
	 * It's a field of some kind.
	 */
	textstuff: string;
	Another: string;
}`,
		},
		{
			Input: MockStructInts{},
			Output: `interface IMockStructInts {
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
}`,
		},
		{
			Input: MockStructNestedOuter{},
			Output: `interface IMockStructNestedOuter {
	inner: IMockStructNestedInner;
}

interface IMockStructNestedInner {
	/**
	 * A really important value.
	 */
	x: number;
}`,
		},
		{
			Input: MockStructPointer{},
			Output: `interface IMockStructPointer {
	x: number | null;
}`,
		},
		{
			Input: MockStructNestedCircular{},
			Output: `interface IMockStructNestedCircular {
	circular: IMockStructNestedCircular | null | undefined;
}`,
		},
	}

	for _, test := range tests {
		g := Generator{}
		builder := strings.Builder{}

		err := g.GenerateSingle(&builder, test.Input)

		if err != nil {
			t.Error("failed to generate:", err)
		}

		str := builder.String()

		if str != test.Output {
			t.Errorf("\n----Expected:\n%s\n----but got:\n%s", test.Output, str)
		}
	}
}
