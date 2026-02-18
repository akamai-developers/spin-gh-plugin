package github

type Tools struct {
	Rust   string
	Go     string
	TinyGo string
	Node   string
	Python string
	Spin   string
}

func DefaultTools() Tools {
	return Tools{
		Rust:   "1.92.0",
		Go:     "1.25.7",
		TinyGo: "v0.39.0",
		Python: "3.13.0",
		Node:   "24",
		Spin:   "",
	}
}
