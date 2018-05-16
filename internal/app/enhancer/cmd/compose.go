package cmd

import (
	"github.com/johynpapin/yuuri/internal/pkg/enhancer"
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
	composeCmd.Flags().StringArrayP("query", "q", nil, "The composed query")

	viper.BindPFlag("compose.input", composeCmd.Flags().Lookup("input"))
	viper.BindPFlag("compose.output", composeCmd.Flags().Lookup("output"))
	viper.BindPFlag("compose.query", composeCmd.Flags().Lookup("query"))
}

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Applies a composed query to a data set",
	Long:  `Applies a compound query to a data set`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"input":  viper.GetString("compose.input"),
			"output": viper.GetString("compose.output"),
			"query":  viper.GetStringSlice("compose.query"),
		}).Info("configuration:")

		inputFile, err := os.Open(viper.GetString("compose.input"))
		if err != nil {
			log.WithField("error", err).Fatal("error opening the input file:")
		}
		defer inputFile.Close()

		outputFile, err := os.Create(viper.GetString("compose.output"))
		if err != nil {
			log.WithField("error", err).Fatal("error opening the output file:")
		}
		defer outputFile.Close()

		dec := rdf.NewTripleDecoder(inputFile, rdf.NTriples)
		enc := rdf.NewTripleEncoder(outputFile, rdf.NTriples)

		agrovocEnhancer := enhancer.NewAgrovocEnhancer()
		agrovocEnhancer.Set("download_output", "agrovoc.nt")

		agrovocOutputFile, err := os.Create("agrovoc.nt")
		if err != nil {
			log.WithField("error", err).Fatal("error creating the agrovoc output file:")
		}
		agrovocOutputFile.Close()

		for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
			if err != nil {
				log.WithField("error", err).Fatal("error decoding ntriples:")
			} else {
				triples, err := agrovocEnhancer.Next(triple)
				if err != nil {
					log.WithFields(log.Fields{
						"error":  err,
						"triple": triple,
					}).Fatal("error enhancing this triple:")
				}

				err = encodeTriples(triples, enc)
				if err != nil {
					log.WithFields(log.Fields{
						"error":   err,
						"triples": triples,
					}).Fatal("error encoding triples:")
				}
			}
		}

		triples, err := agrovocEnhancer.End()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("error ending the enhancer:")
		}

		err = encodeTriples(triples, enc)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"triples": triples,
			}).Fatal("error encoding triples:")
		}

		log.Info("done.")
	},
}

func encodeTriples(triples []rdf.Triple, enc *rdf.TripleEncoder) error {
	for _, triple := range triples {
		err := enc.Encode(triple)
		if err != nil {
			return err
		}
	}

	return nil
}
