package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func findAnchoredString(text []byte, anchor string) string {
	textStrings := strings.Split(string(text), "\n")
	anchoredCount := 0
	result := ""
	for _, str := range textStrings {
		if strings.Contains(str, anchor) {
			anchoredCount++
			result = str
		}
	}
	if anchoredCount > 1 {
		fmt.Println("WARN: HTML context has more than one anchored string!")
	}
	return result
}

func getDigits(str string) string {
	result := ""
	for _, ch := range str {
		if unicode.IsDigit(ch) {
			result += string(ch)
		}
	}
	return result
}

func getTargetNumberFromString(str string, anchorBefore, anchorAfter string) int64 {
	iBefore := strings.Index(str, anchorBefore) + len(anchorBefore)
	iAfter := iBefore + strings.Index(str[iBefore:], anchorAfter)
	result, err := strconv.ParseInt(getDigits(str[iBefore:iAfter]), 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func parseVacantHH(HTMLContext []byte) int64 {
	targetStringAnchor := "Сегодня на сайте"
	targetString := findAnchoredString(HTMLContext, targetStringAnchor)
	vacant := getTargetNumberFromString(targetString, "Сегодня на сайте", "ваканс")
	return vacant
}

func parseResumeHH(HTMLContext []byte) int64 {
	targetStringAnchor := "Сегодня на сайте"
	targetString := findAnchoredString(HTMLContext, targetStringAnchor)
	resume := getTargetNumberFromString(targetString, "ваканс", "резюме")
	return resume
}

func parseWorkersATOL(HTMLContext []byte) int64 {
	// заглушка, возвращаем -1 чтобы отключить автоматическое собирание данных
	// метод для парсинга сейчас реализовать не получается, т.к. сайт АТОЛА не доступен из внешней сетки
	return -1
}
