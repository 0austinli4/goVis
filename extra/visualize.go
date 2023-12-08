package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func visualize(goid uint64, w fyne.Window) {
	// Create a label to display the Goroutine ID
	goidLabel := widget.NewLabel(fmt.Sprintf("Goroutine ID: %d", goid))
	// Create a box to hold the label
	box := container.NewVBox(
		widget.NewLabel("goid"),
		goidLabel,
	)
	w.SetContent(box)
}
