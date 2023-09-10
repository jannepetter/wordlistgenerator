package cmd

import (
	"fmt"
	"strings"
	"wlg/fun"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Adds the words to your wordlist. Default file is output.txt.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		outputfileExists := fun.FileExists(outPutfile)
		if !outputfileExists {
			fun.InitFile(outPutfile)
		}
		//check if user provided their list of words, add the adds if user selected adds
		if userInputList != "" {
			mangleLvl := 101
			fun.ReadFileForMangle(userInputList, outPutfile, tempFileAddress, lang, mangleLvl, adds)
		}

		baseFile := lang + "_base.txt"
		basewords := fun.ReadFile("base/" + baseFile)
		fmt.Println(baseFile, lang)
		var categories = fun.GetCategories(&basewords)

		var wordsToSave = make([]string, 0)

		if userInput != "" {
			wordsToSave = append(wordsToSave, strings.Split(userInput, " ")...)
		}

		if names > 0 {
			maleNames := categories["maleNames"]
			femaleNames := categories["femaleNames"]
			wordsToSave = append(wordsToSave, *fun.AppendNames(&maleNames, &femaleNames, names)...)
		}
		if misc {
			wordsToSave = append(wordsToSave, categories["misc"]...)
		}
		if dates != "" {
			wordsToSave = append(wordsToSave, *fun.GetDates(dates)...)
		}
		if count > 0 {
			wordsToSave = append(wordsToSave, fun.CreateCountList(count)...)
		}
		if birthdayYear != "" {
			wordsToSave = append(wordsToSave, *fun.GetBirthYearList(birthdayYear)...)
		}
		if profanity {
			wordsToSave = append(wordsToSave, categories["profanity"]...)
		}
		if adds > 0 {
			wordsToSave = *fun.AddAddings(&wordsToSave, &categories, adds)
		}
		if patterns {
			wordsToSave = append(wordsToSave, categories["patterns"]...)
		}
		if len(wordsToSave) > 0 {
			fun.SaveList(outPutfile, &wordsToSave)
		}
	},
}

func init() {
	rootCmd.AddCommand(baseCmd)
	baseCmd.PersistentFlags().IntVarP(&names, "names", "n", 0, "Use the default names for selected language:\n"+
		"1 - All names\n"+
		"2 - Male names\n"+
		"3 - Female names\n"+
		"4 - All names lower\n"+
		"5 - Male names lower\n"+
		"6 - Female names lower\n")
	baseCmd.PersistentFlags().IntVarP(&adds, "adds", "A", 0, "Endings that will be added to your choice of words:\n"+
		"1 - Common adds\n"+
		"2 - Popular markings\n"+
		"3 - Category 1 adds\n"+
		"4 - Category 2 adds\n")
	baseCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your file to be created")
	baseCmd.PersistentFlags().BoolVarP(&misc, "misc", "m", false, "Use miscellaneous words from the language base list")
	baseCmd.PersistentFlags().BoolVarP(&profanity, "profanity", "p", false, "Use profanity")
	baseCmd.PersistentFlags().BoolVarP(&patterns, "patterns", "P", false, "Use patterns")
	baseCmd.PersistentFlags().StringVarP(&dates, "dates", "d", "", "Add dates to your list. Use e.g -y \"1990 50\" to get all the dates from 1990 to 2040.\n"+
		"\"1990 50 ./\" will create the dates as 1.1.1990 and 1/1/1990. You can use what ever you wish to separete the date.\n"+
		"You can change the format which the dates are presented e.g \"1990 40 . ymd\" will create the dates as 1990.1.1\n")
	baseCmd.PersistentFlags().StringVarP(&userInput, "list", "l", "", "E.g wlg base -l \"one two three\" ")
	baseCmd.PersistentFlags().StringVarP(&userInputList, "List", "L", "", "Address to file of words, if you want to use your own words")
	baseCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "Add count to your list")
	baseCmd.PersistentFlags().StringVarP(&birthdayYear, "year", "y", "", "Use e.g -y \"1990 50\" to get year presentations from 1990 to 2040")
	baseCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
}
