package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type item struct {
	ID         string
	IsSelected bool
	Details    string
}

type items []*item

func (i items) String() string {
	var s string = ""
	for idx, item := range i {
		s += item.ID
		if idx < len(i)-1 {
			s += ","
		}
	}
	return s
}

func (i *items) NOTSlice() []string {
	var s []string
	for _, item := range *i {
		if !item.IsSelected {
			s = append(s, item.ID)
		}
	}
	return s
}

// selectItems() prompts user to select one or more items in the given slice
func selectItems(selectedPos int, allItems items, returnOpposite bool) (items, error) {
	// Always prepend a "Done" item to the slice if it doesn't already exist.
	const doneID = "Done"
	if len(allItems) > 0 && allItems[0].ID != doneID {
		var items = []*item{
			{
				ID: doneID,
			},
		}
		allItems = append(items, allItems...)
	}

	templates := &promptui.SelectTemplates{
		Details: `{{ .ID }}{{if .Details}}: {{ .Details }}{{end}}`,
		Label: `{{if .IsSelected}}
                    ✔
                {{end}} {{ .ID }} - label`,
		Active:   `{{if .IsSelected}}{{ "✔" | green }} {{end}}{{ .ID | cyan }}`,
		Inactive: `{{if .IsSelected}}{{ "✔" | green }} {{end}}{{ .ID | white }}`,
	}

	prompt := promptui.Select{
		Label:        "Module Selector",
		Items:        allItems,
		Templates:    templates,
		Size:         10,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenItem := allItems[selectionIdx]

	if chosenItem.ID != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return selectItems(selectionIdx, allItems, returnOpposite)
	}

	// If the user selected the "Done" item, return
	// all selected items or the opposite if returnOpposite is true
	var selectedItems []*item
	for _, i := range allItems {
		if (returnOpposite && !i.IsSelected) || (!returnOpposite && i.IsSelected) {
			if i.ID == doneID {
				continue
			}

			selectedItems = append(selectedItems, i)
		}

	}
	return selectedItems, nil
}
