package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var (
	data string

	binCmd = &cobra.Command{
		Use:     "bin [flags]",
		Short:   "Save arbitrary binary data",
		PreRunE: preRun,
		RunE:    doPushBinary,
	}
)

func init() {
	binCmd.Flags().StringVarP(
		&data,
		"binary-data",
		"b",
		"",
		"Binary data to save",
	)

	binCmd.MarkFlagRequired("data")
}

func doPushBinary(cmd *cobra.Command, _args []string) error {
	id, err := clientApp.Usecases.Secrets.PushBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		[]byte(data),
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
