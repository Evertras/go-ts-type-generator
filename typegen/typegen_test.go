package typegen

import "testing"

type MockStructEmpty struct{}

type MockStructStrings struct {
	SomeField string `json:"textstuff"`
}

func TestGeneratesBasicInterfacesCorrectly(t *testing.T) {
	tests := []struct {
		Input  interface{}
		Output string
	}{
		{
			Input:  MockStructEmpty{},
			Output: "interface MockStructEmpty {\n}",
		},
		{
			Input: MockStructStrings{},
			Output: `interface MockStructStrings {
	textstuff: string;
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
