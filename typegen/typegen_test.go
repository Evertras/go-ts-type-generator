package typegen

import "testing"

type MockStructEmpty struct{}

type MockStructStrings struct {
	SomeField    string `json:"textstuff"`
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
	NestedField int `json:"x"`
}

type MockStructNestedOuter struct {
	Inner MockStructNestedInner `json:"inner"`
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
}`,
		},
	}

	for _, test := range tests {
		g := Generator{}

		out, err := g.GenerateSingle(test.Input)

		if err != nil {
			t.Error("failed to generate:", err)
		}

		if out != test.Output {
			t.Errorf("Expected %q but got %q", test.Output, out)
		}
	}
}
