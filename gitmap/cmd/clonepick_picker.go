package cmd

// clonepick_picker.go: thin cmd-side glue that drives the bubbletea
// picker exposed by clonepick.RunPicker. Split out of clonepick.go
// so the dispatcher file stays under the 200-line cap and the picker
// failure paths (cancel = exit 130, empty pick = exit 2) live in one
// auditable place.

import (
	"errors"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v16/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v16/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v16/gitmap/constants"
)

// maybeRunClonePickPicker launches the --ask picker when requested,
// replaces plan.Paths with the user's selection, and translates
// ErrPickerCancelled into the spec'd exit-130 path. No-op when ask
// is false so the non-interactive flow stays a single straight-line
// call from runClonePick.
func maybeRunClonePickPicker(plan clonepick.Plan, ask bool) clonepick.Plan {
	if !ask {
		return plan
	}
	picked, err := clonepick.RunPicker(plan)
	if err != nil {
		handleClonePickPickerError(plan, err)
	}
	if len(picked) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgClonePickMissingPaths)
		maybeExitOnCmdFaithfulMismatch()
		os.Exit(2)
	}
	plan.Paths = picked
	plan.UsedAsk = true

	return plan
}

// handleClonePickPickerError centralises the cancel-vs-fatal split.
// Kept separate so maybeRunClonePickPicker stays under the 15-line
// function cap.
func handleClonePickPickerError(plan clonepick.Plan, err error) {
	if errors.Is(err, clonepick.ErrPickerCancelled) {
		fmt.Fprintln(os.Stderr, constants.MsgClonePickUserCancelled)
		maybeExitOnCmdFaithfulMismatch()
		os.Exit(130)
	}
	cliexit.Fail(constants.CmdClonePick, "picker", plan.RepoUrl, err, 1)
}
