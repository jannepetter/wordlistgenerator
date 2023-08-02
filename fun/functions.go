package fun

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "sync"
)

// var linesWrittenInFile int
// var expectedLines int
// var linesWrittenMutex sync.Mutex

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(address string) []string {
	fileAddress, err := os.Open(address)
	Check(err)
	fileScanner := bufio.NewScanner(fileAddress)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	fileAddress.Close()

	return fileLines
}

func SaveList(address string, wordlist []string) {

	// if linesWrittenInFile > 350000000 {
	// 	return
	// }
	f, err := os.OpenFile(address, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err)
	defer f.Close()
	writer := bufio.NewWriterSize(f, 128000)

	var chunksize = 200000
	for i := 0; i < len(wordlist); i += chunksize {
		end := i + chunksize
		if end > len(wordlist) {
			end = len(wordlist)
		}
		writeFile(wordlist[i:end], writer)
	}
}

func writeFile(wordlist []string, writer *bufio.Writer) {

	for i := 0; i < len(wordlist); i++ {
		_, err := writer.WriteString(wordlist[i] + "\n")
		Check(err)
	}
	writer.Flush()

	// if progress bar needed. Maybe dont need if only one thread writing in??
	// linesWrittenMutex.Lock()
	// linesWrittenInFile += len(wordlist)
	// linesWrittenMutex.Unlock()
}

func getGategoryMode(word string) int {
	wordlist := strings.Split(word[6:], " ")
	mode, err := strconv.Atoi(wordlist[0])
	Check(err)
	return mode
}

func GetCategories(basewords *[]string) map[string][]string {

	var categories = map[string][]string{
		"maleNames":   make([]string, 0),
		"femaleNames": make([]string, 0),
		"misc":        make([]string, 0),
		"marks":       make([]string, 0),
		"maleAdds":    make([]string, 0),
		"femaleAdds":  make([]string, 0),
		"profanity":   make([]string, 0),
		"commonAdds":  make([]string, 0),
	}

	categoryMode := 0
	for _, word := range *basewords {

		if strings.HasPrefix(word, "######") {
			mode := getGategoryMode(word)
			categoryMode = mode
			continue
		}

		switch categoryMode {
		case 1:
			categories["maleNames"] = append(categories["maleNames"], word)
		case 2:
			categories["femaleNames"] = append(categories["femaleNames"], word)
		case 3:
			categories["misc"] = append(categories["misc"], word)
		case 4:
			categories["commonAdds"] = append(categories["commonAdds"], word)
		case 5:
			categories["marks"] = append(categories["marks"], word)
		case 6:
			categories["maleAdds"] = append(categories["maleAdds"], word)
		case 7:
			categories["femaleAdds"] = append(categories["femaleAdds"], word)
		case 8:
			categories["profanity"] = append(categories["profanity"], word)
		default:
			continue
		}
	}

	return categories
}
func combineDate(day string, month string, year string, char string, format string) string {

	word := ""
	switch format {
	case "ymd":
		word = year + char + month + char + day
	case "ydm":
		word = year + char + day + char + month
	case "dym":
		word = day + char + year + char + month
	case "myd":
		word = month + char + year + char + day
	case "mdy":
		word = month + char + day + char + year
	default:
		word = day + char + month + char + year
	}
	return word
}
func modDates(day string, month string, year string, charlist []string, format string) []string {
	var datelist []string
	var needMore int = 0

	for _, char := range charlist {
		datelist = append(datelist, combineDate(day, month, year, char, format))
	}

	if len(day) == 1 {
		day = "0" + day
		needMore++
	}
	if len(month) == 1 {
		month = "0" + month
		needMore++
	}
	if needMore > 0 {
		for _, char := range charlist {
			datelist = append(datelist, combineDate(day, month, year, char, format))
		}
	}
	return datelist
}
func GetDates(dateString string) []string {
	values := strings.Split(dateString, " ")
	startingYear, err := strconv.Atoi(values[0])
	Check(err)
	count, err := strconv.Atoi(values[1])
	Check(err)
	separators := ""

	format := "dmy"
	if len(values) == 3 && strings.Contains(values[2], "y") {
		format = values[2]
	} else if len(values) == 3 {
		separators = values[2]
	}

	if len(values) > 3 {
		separators = values[2]
		format = values[3]
	}

	var datelist []string
	var days int = 1
	var month int = 1
	var year int = startingYear
	charList := strings.Split(separators, "")

	if len(charList) == 0 {
		charList = append(charList, "")
	}

	for {
		dayStr := strconv.Itoa(days)
		monthStr := strconv.Itoa(month)
		yearStr := strconv.Itoa(year)
		wordList := modDates(dayStr, monthStr, yearStr, charList, format)

		datelist = append(datelist, wordList...)
		days++

		if days == 32 {
			days = 1
			month++
			if month == 13 {
				month = 1
				year++
			}
		}

		if year > startingYear+count {
			break
		}
	}

	return datelist
}

// add the words after each word on the output file
func AddWordsToOutputFile(wordlist *[]string, address string) {
	inputFile := address
	tempFile := "temporary_add_file.txt"
	// Open the input file for reading
	file, err := os.Open(address)
	Check(err)

	defer file.Close()

	// Create a temporary output file for writing
	tempOutputFile, err := os.Create(tempFile)
	Check(err)

	defer tempOutputFile.Close()

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(file)

	// Check if output file is empty -> just save the wordlist
	if !scanner.Scan() {
		SaveList(tempFile, *wordlist)
	}
	// else {
	// 	firstLine := scanner.Text()
	// 	for _, word := range *wordlist {
	// 		newLine := firstLine + word // Modify the line by adding the additional word

	// 		// Write the modified line to the temporary output file
	// 		_, err := fmt.Fprintln(tempOutputFile, newLine)
	// 		Check(err)

	// 	}
	// }
	// continue for the rest of the lines
	var wordCollector []string
	for scanner.Scan() {
		line := scanner.Text()
		Check(err)
		for _, word := range *wordlist {
			newLine := line + word // Modify the line by adding the additional word
			wordCollector = append(wordCollector, newLine)
			// Write the modified line to the temporary output file
			// _, err := fmt.Fprintln(tempOutputFile, newLine)
			// Check(err)
			if len(wordCollector) > 5000000 {
				SaveList(tempFile, wordCollector)
				wordCollector = wordCollector[:0]
			}

		}
	}
	if len(wordCollector) > 0 {
		SaveList(tempFile, wordCollector)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	// Close the input and temporary output files
	file.Close()
	tempOutputFile.Close()

	// Replace the input file with the temporary output file
	err = os.Rename(tempFile, inputFile)
	Check(err)
	fmt.Println("Words added to each line successfully!")
}

func CreateCountList(count int) []string {
	var wordlist []string
	var leading = len(strconv.Itoa(count))

	for i := 0; i <= count; i++ {

		var countStr = strconv.Itoa(i)

		for len(countStr) < leading {
			wordlist = append(wordlist, countStr)
			countStr = "0" + countStr
		}
		wordlist = append(wordlist, countStr)
	}
	return wordlist
}

func GetBirthYearList(birthdayYear string) []string {
	values := strings.Split(birthdayYear, " ")
	startingYear, err := strconv.Atoi(values[0])
	Check(err)
	count, err := strconv.Atoi(values[1])
	Check(err)
	var wordlist []string

	for i := 0; i <= count; i++ {
		var year = startingYear + i
		var yearStr = strconv.Itoa(startingYear + i)
		var shorterYearStr = strconv.Itoa(year % 100)
		if len(shorterYearStr) == 1 {
			shorterYearStr = "0" + shorterYearStr
		}
		wordlist = append(wordlist, yearStr, shorterYearStr)
	}

	return wordlist
}

func SaveAsLeetSpeak(replacements map[string][]string, words *[]string, address string) {
	var wordlist []string
	for _, word := range *words {
		wordlist = append(wordlist, removeDuplicates(makeCombinations(word, replacements))...)
		if len(wordlist) > 5000000 {
			// save list in chunks before memory runs out
			SaveList(address, wordlist)
			wordlist = wordlist[:0] // clear the list
		}
	}
	// save the rest
	SaveList(address, wordlist)
}
func removeDuplicates(words []string) []string {
	var mapped = make(map[string]string)
	var values []string
	for i := 0; i < len(words); i++ {
		mapped[words[i]] = words[i]
	}
	for _, v := range mapped {
		values = append(values, v)
	}
	return values
}

func makeCombinations(word string, replacements map[string][]string) []string {
	wordlist := []string{word}
	count := 0
	for i := 0; i < len(wordlist); i++ {
		var charArray []string = strings.Split(wordlist[i], "")
		for j := 0; j < len(charArray); j++ {
			values, ok := replacements[charArray[j]]
			if ok {
				copyCharArray := make([]string, len(charArray))
				copy(copyCharArray, charArray)
				for _, val := range values {
					copyCharArray[j] = val
					var newWord string = strings.Join(copyCharArray, "")
					wordlist = append(wordlist, newWord)
					count++
				}
			}
		}
		if i > len(wordlist) && count > 0 {
			i = 0
			count = 0
		}
		if len(wordlist) > 1000 { // should not happen, break out if does.
			break
		}
	}

	return wordlist
}

var finnish_replacements = map[string][]string{
	"a": {"@", "4"},
	"b": {"8"},
	"e": {"3", "€"},
	"i": {"!", "1"},
	"l": {"!", "1"},
	"o": {"0"},
	"s": {"$", "5"},
	"t": {"+"},
	"w": {"vv"},
	"q": {"9"},
	"z": {"2"},
	"ä": {"a"},
	"ö": {"o"},
}

func GetLeetSpeakReplacements(lang string) map[string][]string {
	if lang == "fi" {
		return finnish_replacements
	}
	return finnish_replacements
}

func AddAddings(words *[]string, categories *map[string][]string, adds int) []string {
	var wordlist []string
	categoryMap := *categories
	var addingList []string
	switch adds {
	case 1:
		addingList = categoryMap["commonAdds"]
	case 2:
		addingList = categoryMap["marks"]
	case 3:
		addingList = categoryMap["maleAdds"]
	case 4:
		addingList = categoryMap["femaleAdds"]
	}

	if len(*words) > 0 && len(addingList) > 0 {
		for _, word := range *words {
			for _, item := range addingList {
				wordlist = append(wordlist, word+item)
			}
		}
	} else if len(*words) == 0 && len(addingList) > 0 {
		wordlist = append(wordlist, addingList...)
	} else {
		return *words
	}
	return wordlist
}
