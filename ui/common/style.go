package common

import (
	"github.com/charmbracelet/glamour"
	gansi "github.com/charmbracelet/glamour/ansi"
)

// StyleConfig returns the default Glamour style configuration.
func StyleConfig() gansi.StyleConfig {
	noColor := ""
	s := glamour.DarkStyleConfig
	// This fixes an issue with the default style config. For example
	// highlighting empty spaces with red in Dockerfile type.
	s.Document.StylePrimitive.Color = &noColor
	s.CodeBlock.Chroma.Text.Color = &noColor
	s.CodeBlock.Chroma.Name.Color = &noColor
	return s
}
