package fun

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var CHUNKSIZE int = 100000
var SAVESIZE int = 1000000

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

func SaveList(address string, wordlist *[]string) {

	f, err := os.OpenFile(address, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err)
	defer f.Close()
	writer := bufio.NewWriterSize(f, CHUNKSIZE)

	for i := 0; i < len(*wordlist); i += CHUNKSIZE {
		end := i + CHUNKSIZE
		if end > len(*wordlist) {
			end = len(*wordlist)
		}
		writeFile((*wordlist)[i:end], writer)
	}
}

func writeFile(wordlist []string, writer *bufio.Writer) {

	for i := 0; i < len(wordlist); i++ {
		word := []byte(wordlist[i] + "\n")
		_, err := writer.Write(word)
		Check(err)
	}
	writer.Flush()
}

func getGategoryMode(word string) int {
	wordlist := strings.Split(word[6:], " ")
	mode, err := strconv.Atoi(wordlist[0])
	Check(err)
	return mode
}

func GetCategories(basewords *[]string) map[string][]string {

	var categories = map[string][]string{
		"maleNames":      make([]string, 0),
		"femaleNames":    make([]string, 0),
		"misc":           make([]string, 0),
		"marks":          make([]string, 0),
		"adds_category1": make([]string, 0),
		"adds_category2": make([]string, 0),
		"profanity":      make([]string, 0),
		"commonAdds":     make([]string, 0),
		"patterns":       make([]string, 0),
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
			categories["adds_category1"] = append(categories["adds_category1"], word)
		case 7:
			categories["adds_category2"] = append(categories["adds_category2"], word)
		case 8:
			categories["profanity"] = append(categories["profanity"], word)
		case 9:
			categories["patterns"] = append(categories["patterns"], word)
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
func GetDates(dateString string) *[]string {
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

	return &datelist
}

// add the words after each word on the output file
func AddWordsToOutputFile(wordlist *[]string, address string, outputfile string) {

	file, err := os.Open(outputfile)
	Check(err)

	defer file.Close()

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(file)

	// Check if output file is empty -> just save the wordlist
	var wordCollector []string
	if !scanner.Scan() {
		SaveList(address, wordlist)
	} else {
		firstLine := scanner.Text()
		for _, word := range *wordlist {
			newLine := firstLine + word // Modify the line by adding the additional word
			wordCollector = append(wordCollector, newLine)
			if len(wordCollector) > SAVESIZE {
				SaveList(address, &wordCollector)
				wordCollector = wordCollector[:0]
			}
		}
	}
	// continue for the rest of the lines
	for scanner.Scan() {
		line := scanner.Text()
		Check(err)
		for _, word := range *wordlist {
			newLine := line + word // Modify the line by adding the additional word
			wordCollector = append(wordCollector, newLine)
			if len(wordCollector) > SAVESIZE {
				SaveList(address, &wordCollector)
				wordCollector = wordCollector[:0]
			}
		}
	}
	if len(wordCollector) > 0 {
		SaveList(address, &wordCollector)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	// Close the input and temporary output files
	file.Close()

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

func GetBirthYearList(birthdayYear string) *[]string {
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

	return &wordlist
}

func SaveAsLeetSpeak(replacements map[string][]string, words *[]string, address string) {
	var wordlist []string
	for _, word := range *words {
		wordlist = append(wordlist, removeDuplicates(makeCombinations(word, replacements))...)
		if len(wordlist) > SAVESIZE {
			// save list in chunks before memory runs out
			SaveList(address, &wordlist)
			wordlist = wordlist[:0] // clear the list
		}
	}
	// save the rest
	if len(wordlist) > 0 {
		SaveList(address, &wordlist)
	}
}

func ReadFileForMangle(address string, outputFile string, tempOutputFile string, lang string, mangleLvl int, addLvl int) error {
	if address == "" {
		return errors.New("you need to define which list you want to be mangled with -L \"yourlist.txt\"")
	}
	fileHandle, err := os.Open(address)
	Check(err)
	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
		if len(fileLines) > SAVESIZE {
			mangleAndSave(&fileLines, address, outputFile, tempOutputFile, lang, mangleLvl, addLvl)
			fileLines = fileLines[:0]
		}
	}
	fileHandle.Close()
	// save the rest
	mangleAndSave(&fileLines, address, outputFile, tempOutputFile, lang, mangleLvl, addLvl)

	return nil
}

func mangleAndSave(fileLines *[]string, address string, outputfile string, tempOutputFile string, lang string, strategy int, addLvl int) {

	switch strategy {
	case 1:
		SaveList(tempOutputFile, UpperAndLowerCaseList(fileLines))
	case 2:
		SaveList(tempOutputFile, LightMangle(fileLines, outputfile))
	case 3:
		baseFile := lang + "_base.txt"
		basewords := ReadFile("base/" + baseFile)
		var categories = GetCategories(&basewords)
		combined := AddAddings(fileLines, &categories, 2)
		SaveList(tempOutputFile, combined)
		combined = AddAddings(fileLines, &categories, 1)
		SaveList(tempOutputFile, combined)
	case 4:
		replacements := GetLeetSpeakReplacements(lang)
		SaveAsLeetSpeak(replacements, fileLines, tempOutputFile)

	// cases above 100 are helpers to cope if user provided list is too big to be done in memory
	case 101:
		baseFile := lang + "_base.txt"
		basewords := ReadFile("base/" + baseFile)
		var categories = GetCategories(&basewords)
		fileLines = AddAddings(fileLines, &categories, addLvl)
		SaveList(outputfile, fileLines)
	case 102:
		// add command, read and add to outputfile
		fileHandle, err := os.Open(outputfile)
		Check(err)
		fileScanner := bufio.NewScanner(fileHandle)
		fileScanner.Split(bufio.ScanLines)
		var newFileLines []string
		linesExist := false
		for fileScanner.Scan() {
			newFileLines = append(newFileLines, fileScanner.Text())
			newFileLines = *AppendWords(&newFileLines, fileLines)
			linesExist = true
			SaveList(tempOutputFile, &newFileLines)
			newFileLines = newFileLines[:0]
		}
		fileHandle.Close()
		if !linesExist {
			SaveList(tempOutputFile, fileLines)
		}
	}

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
	"b": {"8", "|3"},
	"e": {"3", "€"},
	"d": {"|)"},
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
	"A": {"@", "4"},
	"B": {"8", "|3"},
	"E": {"3", "€"},
	"D": {"|)"},
	"I": {"!", "1"},
	"L": {"!_", "|_"},
	"O": {"0"},
	"S": {"$", "5"},
	"T": {"+"},
	"W": {"vv"},
	"Z": {"2"},
	"Ä": {"a"},
	"Ö": {"o"},
}

func GetLeetSpeakReplacements(lang string) map[string][]string {
	if lang == "fi" {
		return finnish_replacements
	}
	return finnish_replacements
}

func AddAddings(words *[]string, categories *map[string][]string, adds int) *[]string {
	var wordlist []string
	categoryMap := *categories
	var addingList []string
	switch adds {
	case 1:
		addingList = categoryMap["commonAdds"]
	case 2:
		addingList = categoryMap["marks"]
	case 3:
		addingList = categoryMap["adds_category1"]
	case 4:
		addingList = categoryMap["adds_category2"]
	}

	if len(*words) > 0 && len(addingList) > 0 {
		for _, word := range *words {
			for _, item := range addingList {
				wordlist = append(wordlist, word+item)
			}
		}

	}

	if len(wordlist) > 0 {
		return &wordlist
	} else if len(*words) > 0 {
		return words
	} else {
		return &addingList
	}
}

func AppendNames(maleNames *[]string, femaleNames *[]string, nameChoice int) *[]string {
	var wordlist []string
	index := 0
	toLowerCase := true
	switch nameChoice {
	case 1:
		wordlist = append(wordlist, *maleNames...)
		wordlist = append(wordlist, *femaleNames...)
	case 2:
		wordlist = append(wordlist, *maleNames...)
	case 3:
		wordlist = append(wordlist, *femaleNames...)
	case 4:
		wordlist = append(wordlist, ChangeWordListCasing(maleNames, index, toLowerCase)...)
		wordlist = append(wordlist, ChangeWordListCasing(femaleNames, index, toLowerCase)...)
	case 5:
		wordlist = append(wordlist, ChangeWordListCasing(maleNames, index, toLowerCase)...)
	case 6:
		wordlist = append(wordlist, ChangeWordListCasing(femaleNames, index, toLowerCase)...)
	}
	return &wordlist
}

func upperCaseWordChar(word string, index int) string {
	wordLength := len(word)
	if wordLength == 0 || wordLength-1 < index {
		return word
	}
	return word[:index] + strings.ToUpper(string(word[index])) + word[index+1:]
}

func lowerCaseWordChar(word string, index int) string {
	wordLength := len(word)
	if wordLength == 0 || wordLength-1 < index {
		return word
	}
	return word[:index] + strings.ToLower(string(word[index])) + word[index+1:]
}

func ChangeWordListCasing(words *[]string, index int, toLower bool) []string {
	var wordList []string

	if toLower {
		for _, word := range *words {
			wordList = append(wordList, lowerCaseWordChar(word, index))
		}
	} else {
		for _, word := range *words {
			wordList = append(wordList, upperCaseWordChar(word, index))
		}
	}
	return wordList
}

func AppendWords(words *[]string, additions *[]string) *[]string {
	var wordlist []string
	for _, word := range *words {
		for _, item := range *additions {
			wordlist = append(wordlist, word+item)
		}
	}
	return &wordlist
}

func AppendSlices(words *[]string, slices ...*[]string) {
	for _, slice := range slices {
		*words = append(*words, *slice...)
	}
}

func Combinator(categories *map[string][]string, choice int) *[]string {
	var wordlist []string

	switch choice {
	case 1:
		wordlist = (*categories)["maleNames"]
		wordlist = *UpperAndLowerCaseList(&wordlist)
	case 2:
		wordlist = (*categories)["femaleNames"]
		wordlist = *UpperAndLowerCaseList(&wordlist)
	case 3:
		wordlist = append((*categories)["maleNames"], (*categories)["femaleNames"]...)
		wordlist = *UpperAndLowerCaseList(&wordlist)
	case 4:
		wordlist = (*categories)["misc"]
		wordlist = *UpperAndLowerCaseList(&wordlist)
	case 5:
		profanity := (*categories)["profanity"]
		wordlist = *UpperAndLowerCaseList(&profanity)
	case 6:
		// malenames with category1 adds
		wordlist = *Combinator(categories, 1)
		adds := (*categories)["adds_category1"]
		wordlist = *AppendWords(&wordlist, &adds)
	case 7:
		// femalenames with category2 adds
		wordlist = *Combinator(categories, 2)
		adds := (*categories)["adds_category2"]
		wordlist = *AppendWords(&wordlist, &adds)
	case 8:
		// malenames with category1 adds
		// femalenames with category2 adds
		wordlist = *Combinator(categories, 6)
		wordlist = append(wordlist, *Combinator(categories, 7)...)
	case 9:
		adds := (*categories)["adds_category1"]
		wordlist = append(adds, (*categories)["adds_category2"]...)
	case 10:
		marks := (*categories)["marks"]
		commonAdds := (*categories)["commonAdds"]
		AppendSlices(&wordlist, &marks, &commonAdds)
	case 11:
		adds := Combinator(categories, 9)
		wordlist = append(*adds, *Combinator(categories, 10)...)
	}
	return &wordlist
}

func UpperAndLowerCaseList(words *[]string) *[]string {
	var wordList []string

	for _, word := range *words {
		word1 := strings.ToLower(string(word))
		word2 := strings.ToUpper(string(word[0])) + word[1:]
		wordList = append(wordList, word1)
		wordList = append(wordList, word2)
	}
	return &wordList
}

func UseStrategy(categories *map[string][]string, strategy int, address string, lang string) {

	year := time.Now().Year() - 50
	dateString := strconv.Itoa(year) + " 50"
	var words []string

	switch strategy {
	case 0:
		dates := *GetDates(dateString)
		dateMod := dateString + " ./"
		dates = append(dates, *GetDates(dateMod)...)
		words = append(dates, (*categories)["patterns"]...)
		SaveList(address, &words)
	case 1:
		allNames := Combinator(categories, 3)
		SaveList(address, allNames)

		marksAndCommonAdds := Combinator(categories, 10)
		words = *AppendWords(allNames, marksAndCommonAdds)
		SaveList(address, &words)
	case 2:
		allNames := Combinator(categories, 3)
		yearList := GetBirthYearList(dateString)
		words = *AppendWords(allNames, yearList)
		SaveList(address, &words)

		namesWithGenderAdds := Combinator(categories, 8)
		SaveList(address, namesWithGenderAdds)
	case 3:
		yearList := GetBirthYearList(dateString)
		marks := (*categories)["marks"]
		yearsWithMarks := AppendWords(yearList, &marks)
		allNames := Combinator(categories, 3)
		namesWithYears := AppendWords(allNames, yearsWithMarks)
		SaveList(address, namesWithYears)

		namesWithGenderadds := Combinator(categories, 8)
		namesWithMarks := AppendWords(namesWithGenderadds, &marks)
		SaveList(address, namesWithMarks)

	case 4:
		misc := Combinator(categories, 4)
		SaveList(address, misc)

		adds := Combinator(categories, 10)
		words = *AppendWords(misc, adds)
		SaveList(address, &words)
	case 5:
		misc := Combinator(categories, 4)
		category1and2Adds := Combinator(categories, 9)
		miscWithAdds := AppendWords(misc, category1and2Adds)
		SaveList(address, miscWithAdds)

		yearList := GetBirthYearList(dateString)
		withYears := AppendWords(misc, yearList)
		SaveList(address, withYears)
	case 6:
		misc := Combinator(categories, 4)
		yearList := GetBirthYearList(dateString)
		marks := (*categories)["marks"]
		yearsWithMarks := AppendWords(yearList, &marks)
		miscWithYearAdds := AppendWords(misc, yearsWithMarks)
		SaveList(address, miscWithYearAdds)

		category1and2Adds := Combinator(categories, 9)
		addsWithMarks := AppendWords(category1and2Adds, &marks)
		miscWithAdds := AppendWords(misc, addsWithMarks)
		SaveList(address, miscWithAdds)

	case 7:
		profanity := Combinator(categories, 5)
		SaveList(address, profanity)

		adds := Combinator(categories, 10)
		profanityWithAdds := AppendWords(profanity, adds)
		SaveList(address, profanityWithAdds)
	case 8:
		profanity := Combinator(categories, 5)
		yearList := GetBirthYearList(dateString)
		category1and2Adds := Combinator(categories, 9)
		AppendSlices(&words, category1and2Adds, yearList)
		profanityWithAdds := AppendWords(profanity, &words)
		SaveList(address, profanityWithAdds)
	case 9:
		UseStrategy(categories, 0, address, lang)
		UseStrategy(categories, 1, address, lang)
		UseStrategy(categories, 2, address, lang)
		UseStrategy(categories, 3, address, lang)
		UseStrategy(categories, 4, address, lang)
		UseStrategy(categories, 5, address, lang)
		UseStrategy(categories, 6, address, lang)
		UseStrategy(categories, 7, address, lang)
		UseStrategy(categories, 8, address, lang)
	}
}

func LightMangle(words *[]string, outputfile string) *[]string {
	var wordlist []string
	for _, word := range *words {
		wordlist = append(wordlist, DuplicateWord(string(word), true))
		wordlist = append(wordlist, UpperCaseWord(string(word)))
		wordlist = append(wordlist, ReverseWord(string(word), true))
		wordlist = append(wordlist, string(word))
		if unicode.IsUpper(rune(word[0])) {
			newWord := lowerCaseWordChar(word, 0)
			wordlist = append(wordlist, newWord)
		} else {
			newWord := upperCaseWordChar(word, 0)
			wordlist = append(wordlist, newWord)
		}
		if len(word) <= 4 {
			newWord := word + word + word
			wordlist = append(wordlist, newWord)
		}
	}
	return &wordlist
}

func UpperCaseWord(word string) string {
	return strings.ToUpper(word)
}

func ReverseWord(word string, separate bool) string {
	runes := []rune(word)
	numberPart := ""
	if separate {
		wordPart, number := separateNumbersAndWords(word)
		numberPart = number
		runes = []rune(wordPart)
	}
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	if separate {
		return string(runes) + numberPart
	}
	return string(runes)
}
func DuplicateWord(word string, separate bool) string {

	if separate {
		wordPart, number := separateNumbersAndWords(word)
		var wordBuilder strings.Builder
		wordBuilder.WriteString(wordPart)
		wordBuilder.WriteString(wordPart)
		if len(wordPart) == 0 {
			wordBuilder.WriteString(number)
		}
		wordBuilder.WriteString(number)
		return wordBuilder.String()
	}
	return word + word
}

func separateNumbersAndWords(input string) (string, string) {
	var wordBuilder strings.Builder
	var remainingBuilder strings.Builder
	firstNumber := false
	for _, char := range input {
		isLetterOrSymbol := !firstNumber && (unicode.IsLetter(char) || unicode.IsPunct(char))
		if isLetterOrSymbol {
			wordBuilder.WriteRune(char)
		} else {
			firstNumber = true
			remainingBuilder.WriteRune(char)
		}
	}

	return wordBuilder.String(), remainingBuilder.String()
}

func InitFile(address string) {
	tempFile, err := os.Create(address)
	Check(err)
	tempFile.Close()
}

func RenameTempFile(address string, newname string) {
	err := os.Rename(address, newname)
	Check(err)
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
