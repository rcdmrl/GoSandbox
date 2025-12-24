package v1

import (
	"fmt"

	"github.com/charmbracelet/huh"
	containersv1 "github.com/rcdmrl/go-sandbox/containers/v1"
	fstreev1 "github.com/rcdmrl/go-sandbox/fstree/v1"
	fstreev2 "github.com/rcdmrl/go-sandbox/fstree/v2"
)

const (
	ProjFSTree     = "fstree"
	ProjContainers = "containers"
	ProjSayonara   = "sayonara"

	FSTreeV1 = "v1"
	FSTreeV2 = "v2"
)

type MainForm struct {
	projectName    string
	projectVersion string
	// deps
	tree1 *fstreev1.ParallelDir
	tree2 *fstreev2.ParallelDir
	dc1   *containersv1.DockerCompose
}

func NewMainForm(tree1 *fstreev1.ParallelDir, tree2 *fstreev2.ParallelDir, dc1 *containersv1.DockerCompose) *MainForm {
	return &MainForm{"", "", tree1, tree2, dc1}
}

// Run executes the multi-step flow: pick project, then (if needed) pick a version.
func (f *MainForm) Run() error {
	if err := f.runProjectSelect(); err != nil {
		return err
	}

	// In case projects have more screens, like one to choose versions
	switch f.projectName {
	case ProjFSTree:
		return f.runFSTreeVersionSelect()
	default:
		return nil
	}
}

// Dispatch runs the selected project/version after the user has gone through the TUI options
func (f *MainForm) Dispatch() error {
	switch f.projectName {
	case ProjFSTree:
		switch f.projectVersion {
		case FSTreeV1:
			f.tree1.Run()
		case FSTreeV2:
			f.tree2.Run()
		default:
			return fmt.Errorf("unknown fs tree version %q", f.projectVersion)
		}
	case ProjContainers:
		err := f.dc1.Start()
		if err != nil {
			return err
		}
		defer f.dc1.ShutDown()
		return nil
	case ProjSayonara:
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
					huh.NewOption("FS Tree", ProjFSTree),
					huh.NewOption("Containers", ProjContainers),
					huh.NewOption("Sayonara", ProjSayonara),
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
					huh.NewOption("v1", FSTreeV1),
					huh.NewOption("v2", FSTreeV2),
				).
				Value(&f.projectVersion),
		),
	).Run()
}
