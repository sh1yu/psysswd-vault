package cmd

import (
	"fmt"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configure configurations",
	Long:  `configure configurations`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "view current configure configurations",
	Long:  `view current configure configurations`,
	Args:  cobra.NoArgs,
	Run:   runConfigView,
}

func init() {
	configCmd.AddCommand(configViewCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigView(cmd *cobra.Command, _ []string) {
	vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
	checkError(err)

	configByte, err := yaml.Marshal(vaultConf)
	checkError(err)

	fmt.Println(string(configByte))
}
