package completion

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func GenZshCompletion(cmd *cobra.Command, w io.Writer) error {
	if _, err := fmt.Fprintln(w, "autoload -U compinit; compinit"); err != nil {
		return err
	}
	if err := cmd.Root().GenZshCompletion(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "compdef _clockify-cli clockify-cli"); err != nil {
		return err
	}
	return nil
}
