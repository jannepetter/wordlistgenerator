package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	adds          int
	names         int
	profanity     bool
	patterns      bool
	misc          bool
	lang          string
	outPutfile    string
	dates         string
	userInput     string
	userInputList string
	count         int
	birthdayYear  string
	strategy      int
)

// go run wordlistgenerator create -h
// rootCmd represents the base command when called without any subcommands
var tempFileAddress = "temporary_add_file.txt"

func GetVariables() map[string]interface{} {
	flags := make(map[string]interface{})
	flags["lang"] = lang
	return flags
}

var rootCmd = &cobra.Command{
	Use:   "wlg",
	Short: "A brief description of your application",
	Long:  ``,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
