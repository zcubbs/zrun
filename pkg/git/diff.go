package git

import (
	"github.com/zcubbs/zrun/pkg/bash"
	"strings"
)

func FileHasChanges(repoPath, file, lastCommit, currentCommit string) (bool, error) {
	output, err := bash.ExecuteCmdWithOutput(
		"git", "-C", repoPath, "diff", "--name-only", lastCommit, currentCommit)
	if err != nil {
		return false, err
	}
	return strings.Contains(output, file), nil
}
