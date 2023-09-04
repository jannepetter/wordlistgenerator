package cmd

import (
	"fmt"
	"time"
	"wlg/fun"

	"github.com/spf13/cobra"
)

var mangleCmd = &cobra.Command{
	Use:   "mangle",
	Short: "Mangle your wordlist",
	Long: `Provide your wordlist that you want to be mangled and select a mangling strategy.
	E.g wlg mangle -L "mylist.txt" -s 0 -o "myoutputfile.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		if strategy < 100 {
			err := fun.ReadFileForMangle(userInputList, lang, outPutfile, strategy)
			fun.Check(err)
		}
		startTime := time.Now()
		for i := 0; i < 1000000; i++ {
			fun.ReverseWord("jåöpäs;10uiui--", true)
			// fmt.Println(word)
		}
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)

		fmt.Printf("Elapsed time: %s\n", elapsedTime)

	},
}

func init() {
	rootCmd.AddCommand(mangleCmd)
	mangleCmd.PersistentFlags().StringVarP(&outPutfile, "output", "o", "output.txt", "The address where you want your mangled file to be created.")
	mangleCmd.PersistentFlags().StringVarP(&userInputList, "List", "L", "", "Address to file of words you want mangled.")
	mangleCmd.PersistentFlags().StringVarP(&lang, "lang", "", "fi", "Language")
	mangleCmd.PersistentFlags().IntVarP(&strategy, "strategy", "s", 0, "Mangles your provided wordlist to the output file.\n"+
		"Yourfile and outputfile need to be different.\n\n"+
		"1 - Uppercase and lowercase your list.\n"+
		"2 - Light mangle. Upper and lowercase word, duplicate, reverse, uppercase entire word. Short words are tripled.\n"+
		"3 - Mangle your wordlist with leetspeak. WARNING! this can get huge very quickly.\n")
}
