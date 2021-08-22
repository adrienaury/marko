package train

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type trainFlags struct {
	order int
}

func NewCommand(parent string) *cobra.Command {
	flags := trainFlags{
		order: 1,
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
	cmd.Flags().IntVar(&flags.order, "order", flags.order, "train a chain of order N")

	return cmd
}

func runTrain(cmd *cobra.Command, args []string) {
	flags, err := getFlagsTrain(cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read flags")
		os.Exit(1)
	}

	log.Debug().
		Int("order", flags.order).
		Msg("Running command train")

	chain := gomarkov.NewChain(flags.order)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		chain.Add(strings.Split(scanner.Text(), ""))
	}

	if scanner.Err() != nil {
		log.Error().Err(err).Msg("Failed to train model")
		os.Exit(1)
	}

	jsonObj, _ := json.Marshal(chain)
	os.Stdout.Write(jsonObj)
}

func getFlagsTrain(cmd *cobra.Command) (*trainFlags, error) {
	order, err := cmd.Flags().GetInt("order")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &trainFlags{
		order: order,
	}, nil
}
