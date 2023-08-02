package cmd

import (
	"fmt"
	"strings"
	"wordlistgenerator/fun"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Adds the words to your wordlist",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseFile := lang + "_base.txt"
		basewords := fun.ReadFile("base/" + baseFile)

		var categories = fun.GetCategories(&basewords)

		fmt.Println("flags :", names, dates, misc, outPutfile)
		var wordsToSave = make([]string, 0)

		// check if user provided their list of words
		if userInputList != "" {
			wordsToSave = append(wordsToSave, fun.ReadFile(userInputList)...)
		}

		if userInput != "" {
			wordsToSave = append(wordsToSave, strings.Split(userInput, " ")...)
		}

		if names {
			wordsToSave = append(wordsToSave, categories["femaleNames"]...)
			wordsToSave = append(wordsToSave, categories["maleNames"]...)
		}
		if misc {
			wordsToSave = append(wordsToSave, categories["misc"]...)
		}
		if dates != "" {
			wordsToSave = append(wordsToSave, fun.GetDates(dates)...)
		}
		if count > 0 {
			wordsToSave = append(wordsToSave, fun.CreateCountList(count)...)
		}
		if birthdayYear != "" {
			wordsToSave = append(wordsToSave, fun.GetBirthYearList(birthdayYear)...)
		}
		if profanity {
			wordsToSave = append(wordsToSave, categories["profanity"]...)
		}
		if adds > 0 {
			wordsToSave = fun.AddAddings(&wordsToSave, &categories, adds)
		}
		if leetSpeak && len(wordsToSave) > 0 {
			replacements := fun.GetLeetSpeakReplacements(lang)
			fun.SaveAsLeetSpeak(replacements, &wordsToSave, outPutfile)
		} else if len(wordsToSave) > 0 {
			fun.SaveList(outPutfile, wordsToSave)
		}
	},
}

func init() {
	rootCmd.AddCommand(baseCmd)
	baseCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Use the default names for selected language")
	baseCmd.PersistentFlags().IntVarP(&adds, "adds", "A", 0, "Endings that will be added to your choice of words:\n"+
		"1 - Common adds\n"+
		"2 - Popular markings\n"+
		"3 - Male adds\n"+
		"4 - Female adds\n"+
		"If you choose only the addings, the addings only will be added to the output list.\n")
	baseCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your file to be created")
	baseCmd.PersistentFlags().BoolVarP(&misc, "misc", "m", false, "Use miscellaneous words from the language base list")
	baseCmd.PersistentFlags().BoolVarP(&profanity, "profanity", "p", false, "Use profanity")
	baseCmd.PersistentFlags().StringVarP(&dates, "dates", "d", "", "Add dates to your list. Use e.g -y \"1990 50\" to get all the dates from 1990 to 2040.\n"+
		"\"1990 50 ./\" will create the dates as 1.1.1990 and 1/1/1990. You can use what ever you wish to separete the date.\n"+
		"You can change the format which the dates are presented e.g \"1990 40 . ymd\" will create the dates as 1990.1.1\n")
	baseCmd.PersistentFlags().StringVarP(&userInput, "list", "l", "", "E.g wordlistgenerator base -l \"one two three\" ")
	baseCmd.PersistentFlags().StringVarP(&userInputList, "List", "L", "", "Address to file of words, if you want to use your own words")
	baseCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "Add count to your list")
	baseCmd.PersistentFlags().StringVarP(&birthdayYear, "year", "y", "", "Use e.g -y \"1990 50\" to get year presentations from 1990 to 2040")
	baseCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
	baseCmd.PersistentFlags().BoolVarP(&leetSpeak, "leetspeak", "E", false, "Add the selected words in as leetspeak.\n"+
		"Try it at first with smaller lists. The list gets really big pretty fast\n")

}
