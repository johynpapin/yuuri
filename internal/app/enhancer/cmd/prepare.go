package cmd

import (
	"github.com/johynpapin/yuuri/internal/pkg/monoprix"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(prepareCmd)

	prepareCmd.Flags().StringP("output", "o", "", "Output file")
	prepareCmd.Flags().StringP("assets", "a", "assets", "Assets folder path")
	prepareCmd.Flags().StringP("torip", "", "127.0.0.1", "IP of the TOR proxy")
	prepareCmd.Flags().IntP("torport", "", 9050, "PORT of the TOR proxy")
	prepareCmd.Flags().IntP("torcontrolport", "", 9051, "PORT of TOR control")
	prepareCmd.Flags().StringP("torcontrolpassword", "", "", "Password of TOR control")
	prepareCmd.Flags().IntP("orchestratorport", "", 4242, "Port of the orchestrator")
	prepareCmd.Flags().IntP("workers", "", 10, "Number of workers")

	viper.BindPFlag("prepare.output", prepareCmd.Flags().Lookup("output"))
	viper.BindPFlag("prepare.assets", prepareCmd.Flags().Lookup("assets"))
	viper.BindPFlag("prepare.torip", prepareCmd.Flags().Lookup("torip"))
	viper.BindPFlag("prepare.torport", prepareCmd.Flags().Lookup("torport"))
	viper.BindPFlag("prepare.torcontrolport", prepareCmd.Flags().Lookup("torcontrolport"))
	viper.BindPFlag("prepare.torcontrolpassword", prepareCmd.Flags().Lookup("torcontrolpassword"))
	viper.BindPFlag("prepare.orchestratorport", prepareCmd.Flags().Lookup("orchestratorport"))
	viper.BindPFlag("prepare.workers", prepareCmd.Flags().Lookup("workers"))
}

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Applies a composed query to a data set",
	Long:  `Applies a compound query to a data set`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"output": viper.GetString("prepare.output"),
		}).Info("configuration:")

		monoprixScraper := monoprix.NewMonoprixScraper(monoprix.Config{
			OutputPath:         viper.GetString("prepare.output"),
			AssetsPath:         viper.GetString("prepare.assets"),
			TORIP:              viper.GetString("prepare.torip"),
			TORPort:            viper.GetInt("prepare.torport"),
			TORControlPort:     viper.GetInt("prepare.torcontrolport"),
			TORControlPassword: viper.GetString("prepare.torcontrolpassword"),
			OrchestratorPort:   viper.GetInt("prepare.orchestratorport"),
			Workers:            viper.GetInt("prepare.workers"),
		})

		err := monoprixScraper.Start()
		if err != nil {
			log.WithField("error", err).Fatal("error when scraping monoprix:")
		}

		log.Info("done.")
	},
}
