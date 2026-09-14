package github

type Tools struct {
	Rust       string
	RustTarget string
	Go         string
	Node       string
	Python     string
	Spin       string
}

func DefaultTools() Tools {
	return Tools{
		Rust:       "1.98.1",
		RustTarget: "wasm32-wasip2",
		Go:         "1.26.6",
		Python:     "3.13.0",
		Node:       "26",
		Spin:       "", // empty string is latest stable
	}
}
