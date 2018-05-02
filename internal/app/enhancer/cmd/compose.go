package cmd

import (
	"github.com/knakk/rdf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
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
			"input":  viper.GetString("compose.input"),
			"output": viper.GetString("compose.output"),
		}).Info("configuration:")

		f, err := os.Open(viper.GetString("compose.input"))
		if err != nil {
			log.WithField("error", err).Fatal("error opening the input file:")
		}
		dec := rdf.NewTripleDecoder(f, rdf.NTriples)
		for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
			if err != nil {
				log.WithField("error", err).Fatal("error decoding ntriples:")
			} else {
				log.WithField("triple", triple).Info("triple decoded:")
			}
		}
	},
}
