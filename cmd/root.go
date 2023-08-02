package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	adds          int
	names         bool
	profanity     bool
	misc          bool
	lang          string
	outPutfile    string
	dates         string
	userInput     string
	userInputList string
	count         int
	birthdayYear  string
	leetSpeak     bool
)

// go run wordlistgenerator create -h
// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "wordlistgenerator",
	Short: "A brief description of your application",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// rootCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your file to be created")
	// rootCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
	// rootCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Use the default names for selected language")
	// rootCmd.PersistentFlags().BoolVarP(&profanity, "profanity", "p", false, "Use profanity")
	// rootCmd.PersistentFlags().BoolVarP(&misc, "misc", "m", false, "Use miscellaneous words.. or something")
	// rootCmd.PersistentFlags().BoolVarP(&gender, "gender", "g", false, "Use gender thingies")
	// rootCmd.PersistentFlags().StringVarP(&dates, "dates", "d", "", "Use dates")
	// rootCmd.PersistentFlags().StringVarP(&userInput, "list", "l", "", "A string list of words you like to use. Separated by space")
	// rootCmd.PersistentFlags().StringVarP(&userInputList, "List", "L", "", "Address to file of words, if you want to use your own words")
	// rootCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "add count 0 - count at the end of your word.")
	// rootCmd.PersistentFlags().BoolVarP(&leetSpeak, "leetspeak", "E", false, "Add the words in as leetspeak.")

	// rootCmd.PersistentFlags().MarkHidden("profanity")
	// rootCmd.PersistentFlags().MarkHidden("names")
	// rootCmd.PersistentFlags().MarkHidden("lang")
	// rootCmd.PersistentFlags().MarkHidden("misc")
	// rootCmd.PersistentFlags().MarkHidden("dates")
	// rootCmd.PersistentFlags().MarkHidden("gender")
	// rootCmd.PersistentFlags().MarkHidden("mixwords")
	// rootCmd.PersistentFlags().MarkHidden("output")
	// rootCmd.PersistentFlags().MarkHidden("list")
	// rootCmd.PersistentFlags().MarkHidden("List")
	// rootCmd.PersistentFlags().MarkHidden("count")
	// rootCmd.PersistentFlags().MarkHidden("leetspeak")
}
