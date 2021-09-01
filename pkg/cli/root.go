package cli

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chore",
		Short: "Elegant SSH tasks runner",
		Long:  "Elegant SSH tasks runner",
	}

	cmd.AddCommand(
		NewInitCommand(),
		NewRunCommand(),
	)

	return cmd
}
