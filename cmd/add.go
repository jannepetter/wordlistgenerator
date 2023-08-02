package cmd

import (
	"fmt"
	"strings"
	"wordlistgenerator/fun"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the words at the end of each word in your outputfile. Default file is output.txt.",
	Long:  `You can use this to build your wordlist from ground up.`,
	Run: func(cmd *cobra.Command, args []string) {
		baseFile := lang + "_base.txt"
		basewords := fun.ReadFile("base/" + baseFile)
		var categories = fun.GetCategories(&basewords)

		fmt.Println("flags :", names, dates, misc, outPutfile)

		var wordsToAdd = make([]string, 0)
		// check if user provided their list of words
		if userInputList != "" {
			wordsToAdd = append(wordsToAdd, fun.ReadFile(userInputList)...)
		}

		if userInput != "" {
			wordsToAdd = append(wordsToAdd, strings.Split(userInput, " ")...)
		}

		if names {
			wordsToAdd = append(wordsToAdd, categories["femaleNames"]...)
			wordsToAdd = append(wordsToAdd, categories["maleNames"]...)
		}
		if misc {
			wordsToAdd = append(wordsToAdd, categories["misc"]...)
		}

		if dates != "" {
			wordsToAdd = append(wordsToAdd, fun.GetDates(dates)...)
		}
		if count > 0 {
			wordsToAdd = append(wordsToAdd, fun.CreateCountList(count)...)
		}

		if birthdayYear != "" {
			wordsToAdd = append(wordsToAdd, fun.GetBirthYearList(birthdayYear)...)
		}
		if profanity {
			wordsToAdd = append(wordsToAdd, categories["profanity"]...)
		}
		if adds > 0 {
			wordsToAdd = fun.AddAddings(&wordsToAdd, &categories, adds)
		}

		if len(wordsToAdd) > 0 {
			fun.AddWordsToOutputFile(&wordsToAdd, outPutfile)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Use the default names for selected language")
	addCmd.PersistentFlags().IntVarP(&adds, "adds", "A", 0, "Endings that will be added to your choice of words and then added to the ouputfile:\n"+
		"1 - Common adds\n"+
		"2 - Popular markings\n"+
		"3 - Male adds\n"+
		"4 - Female adds\n"+
		"If you choose only the addings, the addings only will be added to the output list")
	addCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your file to be created")
	addCmd.PersistentFlags().BoolVarP(&misc, "misc", "m", false, "Use miscellaneous words from the language base list")
	addCmd.PersistentFlags().BoolVarP(&profanity, "profanity", "p", false, "Use profanity")
	addCmd.PersistentFlags().StringVarP(&dates, "dates", "d", "", "Add dates to your list. Use e.g -y \"1990 50\" to get all the dates from 1990 to 2040.\n"+
		"\"1990 50 ./\" will create the dates as 1.1.1990 and 1/1/1990. You can use what ever you wish to separete the date.\n"+
		"You can change the format which the dates are presented e.g \"1990 40 . ymd\" will create the dates as 1990.1.1")
	addCmd.PersistentFlags().StringVarP(&userInput, "list", "l", "", "E.g wordlistgenerator base -l \"one two three\" ")
	addCmd.PersistentFlags().StringVarP(&userInputList, "List", "L", "", "Address to file of words, if you want to use your own words")
	addCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "Add count at the end of your words.")
	addCmd.PersistentFlags().StringVarP(&birthdayYear, "year", "y", "", "Use e.g -y \"1990 50\" to get year presentations from 1990 to 2040")
	addCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
}
