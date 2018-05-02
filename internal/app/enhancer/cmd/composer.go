package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(composeCmd)

	composeCmd.Flags().StringP("input", "i", "", "Input file")
	composeCmd.Flags().StringP("output", "o", "", "Output file")

	viper.BindPFlag("compose.input", composeCmd.Flags().Lookup("input"))
	viper.BindPFlag("compose.output", composeCmd.Flags().Lookup("output"))
}

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Applies a composed query to a data set",
	Long:  `Applies a compound query to a data set`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"input": viper.GetString("compose.input"),
			"output": viper.GetString("compose.output"),
		}).Info("configuration:")
	},
}
