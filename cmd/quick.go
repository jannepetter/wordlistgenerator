package cmd

import (
	"wlg/fun"

	"github.com/spf13/cobra"
)

var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Use ready made strategies to form your list.",
	Long:  `Use ready made strategies to form your list.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputfileExists := fun.FileExists(outPutfile)
		if !outputfileExists {
			fun.InitFile(outPutfile)
		}
		baseFile := lang + "_base.txt"
		basewords := fun.ReadFile("base/" + baseFile)
		var categories = fun.GetCategories(&basewords)
		fun.UseStrategy(&categories, strategy, outPutfile, lang)
	},
}

func init() {
	rootCmd.AddCommand(quickCmd)
	quickCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your file to be created")
	quickCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
	quickCmd.PersistentFlags().IntVarP(&strategy, "strategy", "s", 0, "Ready made strategies using language base file.\n\n"+
		"0 - Dates from the past 50 years in different formats and patterns.\n"+
		"1 - Names with lower and uppercase, marks, and common addings.\n"+
		"2 - Names with lower and uppercase with birth year (past 50 years). E.g john08, john2008 + additional addings.\n"+
		"3 - Names with lower and uppercase with birth years and marks. + additional addings\n"+
		"4 - Miscellaneous words from base language file combined with marks and common addings\n"+
		"5 - Miscellaneous words from base language file combined with birth year (50 years) + additional addings.\n"+
		"6 - Miscellaneous words from base language file combined with years, marks and category 1 & 2 adds.\n"+
		"7 - Profanity from base file with marks and common adds.\n"+
		"8 - Profanity from base file with birth year (50 years) + additional addings.\n"+
		"9 - Do all the previous at one go.\n")
}
