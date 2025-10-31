package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/alex-vee-sh/veessh/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "veessh",
	Short: "Console connection manager for SSH/SFTP/Telnet and more",
}

// Execute is the entrypoint for the Cobra command tree.
func Execute() error {
	addSubcommands()
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.Version = version.String()
	rootCmd.SetVersionTemplate(versionTemplate)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return rootCmd.ExecuteContext(ctx)
}

func addSubcommands() {
	rootCmd.AddCommand(cmdAdd)
	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdShow)
	rootCmd.AddCommand(cmdConnect)
	rootCmd.AddCommand(cmdRemove)
	rootCmd.AddCommand(cmdPick)
	rootCmd.AddCommand(cmdFavorite)
	rootCmd.AddCommand(cmdExport)
	rootCmd.AddCommand(cmdImport)
	rootCmd.AddCommand(cmdImportSSH)
	rootCmd.AddCommand(cmdCompletion)
	rootCmd.AddCommand(cmdVersion)
}

var flagJSON bool
var flagVersionShort bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "output JSON where supported")
	rootCmd.PersistentFlags().BoolVarP(&flagVersionShort, "version", "v", false, "show version and exit")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if flagVersionShort {
			fmt.Fprint(cmd.OutOrStdout(), renderVersion())
			os.Exit(0)
		}
		return nil
	}
}

func OutputJSON() bool { return flagJSON }

const versionTemplate = `{{.Name}} {{.Version}}

 /\_/\   veessh
( o.o )  {{.Version}}
 > ^ <

`

func renderVersion() string {
	return fmt.Sprintf(`%s %s

 /\\_/\\   veessh
( o.o )  %s
 > ^ <

`, rootCmd.Name(), rootCmd.Version, rootCmd.Version)
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
