package v1

import (
	"fmt"

	"github.com/charmbracelet/huh"
	fstreev1 "github.com/rcdmrl/go-sandbox/fstree/v1"
	fstreev2 "github.com/rcdmrl/go-sandbox/fstree/v2"
)

type MainForm struct {
	projectName    string
	projectVersion string

	tree1 *fstreev1.ParallelDir
	tree2 *fstreev2.ParallelDir
}

func NewMainForm(tree1 *fstreev1.ParallelDir, tree2 *fstreev2.ParallelDir) *MainForm {
	return &MainForm{
		"",
		"",
		tree1,
		tree2,
	}
}

// Run executes the multi-step flow: pick project, then (if needed) pick a version.
func (f *MainForm) Run() error {
	if err := f.runProjectSelect(); err != nil {
		return err
	}

	switch f.projectName {
	case "fstree":
		return f.runFSTreeVersionSelect()
	case "sayonara":
		return nil
	default:
		return fmt.Errorf("unknown project %q", f.projectName)
	}
}

// Dispatch runs the selected project/version after the user has gone through the TUI options
func (f *MainForm) Dispatch() error {
	switch f.projectName {
	case "fstree":
		switch f.projectVersion {
		case "v1":
			f.tree1.Run()
		case "v2":
			f.tree2.Run()
		default:
			return fmt.Errorf("unknown fs tree version %q", f.projectVersion)
		}
	case "sayonara":
		fmt.Println("You called quits. Cya!")
	default:
		return fmt.Errorf("what's %q dude?", f.projectName)
	}
	return nil
}

// runProjectSelect shows the top-level project chooser
func (f *MainForm) runProjectSelect() error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which project?").
				Options(
					huh.NewOption("FS Tree", "fstree"),
					huh.NewOption("Sayonara", "sayonara"),
				).
				Value(&f.projectName),
		),
	).Run()
}

// runFSTreeVersionSelect shows the for the fs tree project. version chooser.
func (f *MainForm) runFSTreeVersionSelect() error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which fs tree version?").
				Options(
					huh.NewOption("v1", "v1"),
					huh.NewOption("v2", "v2"),
				).
				Value(&f.projectVersion),
		),
	).Run()
}
