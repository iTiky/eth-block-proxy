package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VersionGitTag    string
	VersionGitCommit string
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s", VersionGitTag, VersionGitCommit)
}

// newVersionCmd creates a new root.version cobra.Command.
func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print app version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(GetVersion())

			return nil
		},
	}

	return cmd
}
