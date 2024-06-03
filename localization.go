package main

import (
	"bufio"
	"os"
	"strings"
)

var lang []TranslationPair
var language string = "en"

type TranslationPair struct {
	MsgID  string
	MsgStr string
}

// return object containing the localization, reading .po files, getting all the msgid and msgstr pairs
func generateLocalization() []TranslationPair {

	// read po file
	file, err := os.Open("languages/" + language + ".po")
	if err != nil {
		return nil
	}
	defer file.Close()

	var translations []TranslationPair

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 5 { // if line length is less than 5, continue
			continue
		}
		if line[:5] == "msgid" { // if the line starts with msgid, then we are starting a new translation
			msgid := line[5:]
			msgid = strings.TrimSpace(msgid)
			msgid = strings.Trim(msgid, "\"")
			translations = append(translations, TranslationPair{MsgID: msgid})
		} else if line[:6] == "msgstr" { // if the line starts with msgstr, it's paired with the previous msgid removing the quotes
			msgstr := line[6:]
			msgstr = strings.TrimSpace(msgstr)
			msgstr = strings.Trim(msgstr, "\"")
			translations[len(translations)-1].MsgStr = msgstr
		}
	}

	return translations
}

func getLocalization(msgid string) string {
	// find the msgstr that corresponds to the msgid
	for _, translation := range lang {
		if translation.MsgID == msgid {
			return translation.MsgStr
		}
	}
	return ""
}
