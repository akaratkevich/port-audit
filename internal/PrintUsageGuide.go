package internal

import (
	"github.com/pterm/pterm"
)

// PrintUsageGuide asks the user if they want to see the usage guide and displays it if they choose to do so.
func PrintUsageGuide(guide string) bool {
	//// Ask the user if they want to see the usage guide
	//showGuide, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Y").Show("Do you want to see the usage guide? (Y/N)")
	//showGuide = strings.ToUpper(strings.TrimSpace(showGuide))
	//
	//if showGuide != "Y" {
	//	return false
	//}

	// Display the usage guide
	paddedBox := pterm.DefaultBox.WithLeftPadding(2).WithRightPadding(2).WithTopPadding(1).WithBottomPadding(1)
	title := pterm.LightRed("Usage Guide")
	box := paddedBox.WithTitle(title).WithTitleTopLeft().Sprint(pterm.NewStyle(pterm.FgLightWhite, pterm.Italic).Sprint(guide))

	pterm.DefaultPanel.WithPanels([][]pterm.Panel{
		{{box}},
	}).Render()

	return true
}
