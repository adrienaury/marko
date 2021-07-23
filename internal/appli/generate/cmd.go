package generate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type flags struct {
	limit uint64
}

func NewCommand(parent string) *cobra.Command {
	flags := flags{
		limit: 1,
	}

	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate values from a Markov chain",
		Run:     run,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  %s generate --limit 3 < model.json`, parent),
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Uint64Var(&flags.limit, "limit", flags.limit, "limit the number of results")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	flags, err := getFlags(cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read flags")
		os.Exit(1)
	}

	chain, err := load(os.Stdin)
	if err != nil {
		log.Error().Err(err).Msg("Failed to train model")
		os.Exit(1)
	}

	var i uint64
	for i = 0; i < flags.limit; i++ {
		generate(chain)
	}

	log.Debug().
		Uint64("limit", flags.limit).
		Msg("Running command generate")
}

func getFlags(cmd *cobra.Command) (*flags, error) {
	limit, err := cmd.Flags().GetUint64("limit")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &flags{
		limit: limit,
	}, nil
}

func load(in *os.File) (*gomarkov.Chain, error) {
	var chain gomarkov.Chain

	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}

	return &chain, nil
}

func generate(chain *gomarkov.Chain) {
	order := chain.Order
	tokens := make([]string, 0)
	for i := 0; i < order; i++ {
		tokens = append(tokens, gomarkov.StartToken)
	}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - order):])
		tokens = append(tokens, next)
	}
	fmt.Println(strings.Join(tokens[order:len(tokens)-1], ""))
}
