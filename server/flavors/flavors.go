/*
Package flavors provides an enumeration of builtin CSS themes
and a method to retrieve the css.

Source of Themes:

  - GitHub
    https://github.com/sindresorhus/github-markdown-css
  - Pico
    picocss.com (red classless theme)
*/
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
