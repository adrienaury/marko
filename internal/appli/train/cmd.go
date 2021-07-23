package train

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type trainFlags struct {
	order uint64
}

func NewCommand(parent string) *cobra.Command {
	flags := trainFlags{
		order: 0,
	}

	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:     "train",
		Short:   "Train a markov chain from a dataset",
		Run:     runTrain,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  %s train --order 2 < dataset.txt > model.json`, parent),
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Uint64Var(&flags.order, "order", flags.order, "train a chain of order N")

	return cmd
}

func runTrain(cmd *cobra.Command, args []string) {
	flags, err := getFlagsTrain(cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read flags")
		os.Exit(1)
	}

	log.Debug().
		Uint64("order", flags.order).
		Msg("Running command train")

	if err := os.Chmod("", os.ModeAppend); err != nil {
		log.Error().Err(err).Msg("Failed to train model")
		os.Exit(1)
	}
}

func getFlagsTrain(cmd *cobra.Command) (*trainFlags, error) {
	order, err := cmd.Flags().GetUint64("order")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &trainFlags{
		order: order,
	}, nil
}
