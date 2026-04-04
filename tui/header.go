package tui

import "github.com/charmbracelet/lipgloss"

// Header is the app's header component.
type Header struct {
	content string
	style   lipgloss.Style
}

// NewHeader creates the header of the app.
func NewHeader() *Header {
	s := lipgloss.NewStyle().
		Bold(true).
		Padding(1).
		AlignHorizontal(lipgloss.Center).
		Width(viewWidth)

	return &Header{
		content: "\U0001f345 Pimpomodoro Timer\n",
		style:   s,
	}
}

// View renders the header.
func (h *Header) View() string {
	return h.style.Render(h.content)
}
