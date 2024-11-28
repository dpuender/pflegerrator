package service

import (
	"fmt"
	"log"
	"pflegerrator/structs"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func RvGenerator(person structs.Person) string {
	firstletter := strings.ToUpper(string(person.LastName[0]))
	firstletterPosition := getAplphabeticalPosition(rune(firstletter[0]))
	formattedBirthDate := formatBirthDate(person.BirthDate)
	areaNumber := getRandomAreaNumber()
	birthMonthEquivalent := getBirthMonthEquivalent(person.BirthDate, person.Sex)
	tmpRVNr := areaNumber + formattedBirthDate + firstletterPosition + birthMonthEquivalent
	checkDigits := generateCheckDigits(tmpRVNr)
	rvNr := areaNumber + formattedBirthDate + firstletter + birthMonthEquivalent + strconv.Itoa(checkDigits)
	return rvNr
}

func getAplphabeticalPosition(letter rune) string {
	alphabeticalPosition := letter - 'A' + 1
	if alphabeticalPosition < 10 {
		return "0" + strconv.Itoa(int(alphabeticalPosition))
	} else {
		return strconv.Itoa(int(alphabeticalPosition))
	}
}

func formatBirthDate(birthDate string) string {

	parsedDate := formatDate(birthDate)

	return parsedDate.Format("020106")
}

func getBirthMonthEquivalent(birthDate string, sex string) string {

	parsedDate := formatDate(birthDate)

	birthMonthEquivalent := int(parsedDate.Month())

	if sex == "w" {
		birthMonthEquivalent += 50
	}

	if birthMonthEquivalent < 10 {
		return "0" + strconv.Itoa(birthMonthEquivalent)
	} else {
		return strconv.Itoa(birthMonthEquivalent)
	}
}

func getRandomAreaNumber() string {
	var areaNumbers = [31]string{"02", "03", "04", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "23", "24", "25", "26", "28", "29", "38", "78", "39", "79", "80", "81", "82", "89"}

	areaNumber := areaNumbers[rand.Intn(31)]

	return areaNumber
}

func generateCheckDigits(tmpRVNr string) int {
	var checksum int

	checksum += parseInt(tmpRVNr[0]) * 2
	checksum += parseInt(tmpRVNr[1]) * 1
	checksum += parseInt(tmpRVNr[2]) * 2
	checksum += parseInt(tmpRVNr[3]) * 5
	checksum += parseInt(tmpRVNr[4]) * 7
	checksum += parseInt(tmpRVNr[5]) * 1
	checksum += parseInt(tmpRVNr[6]) * 2
	checksum += parseInt(tmpRVNr[7]) * 1
	checksum += parseInt(tmpRVNr[8]) * 2
	checksum += parseInt(tmpRVNr[9]) * 1
	checksum += parseInt(tmpRVNr[10]) * 2
	checksum += parseInt(tmpRVNr[11]) * 1

	return digitSum(checksum) % 10
}

func digitSum(number int) int {

	var digitSum int

	for number > 0 {
		digitSum += digitSum + (number % 10)
		number = number - (number % 10)
		number = number / 10
	}

	return digitSum
}

func parseInt(s byte) int {
	parsedInt, err := strconv.Atoi(string(s))
	if err != nil {
		log.Fatal(err)
	}
	return parsedInt
}

func formatDate(date string) time.Time {
	formats := []string{
		"2006-01-02",
		time.RFC3339,
		"02.01.2006",
	}

	var parsedDate time.Time
	var err error

	for _, format := range formats {
		parsedDate, err = time.Parse(format, date)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Println("Fehler beim Parsen des Datums:", err)
	}

	return parsedDate
}
