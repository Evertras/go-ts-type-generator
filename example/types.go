package example

//go:generate go run generate.go

// SomeData is some data
type SomeData struct {
	X int `json:"x"`
	Y uint64
	Z string
}

// Outer holds Inner
type Outer struct {
	InnerStuff Inner `json:"inner"`
}

// Inner is held by Outer
type Inner struct {
	X *int `json:"x,omitempty"`
	Y *int `json:"y"`
}
