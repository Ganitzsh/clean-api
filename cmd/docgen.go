package cmd

import (
	"fmt"

	"github.com/ganitzsh/f3-te/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"
	"github.com/spf13/cobra"
)

var docgenCmd = &cobra.Command{
	Use:   "docgen",
	Short: "Output the routes as Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(docgen.MarkdownRoutesDoc((api.Routes()).(*chi.Mux), docgen.MarkdownOpts{
			ProjectPath: "github.com/ganitzsh/f3-te",
			Intro:       "Generated doc for Payment API",
		}))
	},
}
