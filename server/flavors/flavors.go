package flavors

import _ "embed"

type Enum int

const (
	GitHub Enum = iota
	Pico
)

//go:embed css/github.css
var githubCSS string

//go:embed css/pico.css
var picoCSS string

func (enum Enum) GetCss() string {
	switch enum {
	case GitHub:
		return githubCSS
	case Pico:
		return picoCSS
	default:
		panic("Missing Enum")
	}
}
