package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID()string{
  return uuid.New().String()
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// returns 0 if the string isn't a string
func StringToInt(str string)(int){
  val,err := strconv.Atoi(str)
  if err != nil{
    log.Println(err)
    return 0
  }
  return val
}

// theres a timestamp for this but I couldn't find it
var now = time.Now()
var CurrentTime = now.Format("2006-01-02 15:04:05")

func IntToString(val int) string{  return strconv.Itoa(val) }
// learn go generics and avoid such issues
//check if int is in array

func ArrayContainsInt(arr []int,element int) bool{
  for _,e := range arr {
    if e == element{
      return true
    }
  }
  return false
}

func ArrayContainsString(arr []string, element string) bool{
  for _,e := range arr {
    if e == element{
      return true
    }
  }
  return false
}

//Check if a string is empty returns True if string is a string
func CheckifStringIsEmpty(data string) bool{
  if len(strings.TrimSpace(data)) == 0{
    return false
  }
  if len(data) == 0{
    return false
  }
  return true
}

func TrueRand(len int) string{
  bytes := make([]byte,len)
  for i := 0; i < len; i++{
    bytes[i] = byte(randInt(97,122))
  }
  if !CheckifStringIsEmpty(string(bytes)){
    TrueRand(len)
  }
  return string(bytes)
}

func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}

func RandString(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("QWERTYUIOPLKJHGFDSAZXCVBNM123456789qwertyuioplkjhgfdsazxcvbnm")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  id = strings.ToUpper(id)
  if !CheckifStringIsEmpty(id){
    RandString(length)
  }
  return id
}

//Retrns a random string with numbers and letters (caps on)
func RandNoLetter(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("QWERTYUIOPLKJHGFDSAZXCVBNM123456789")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  id = strings.ToUpper(id)
  if !CheckifStringIsEmpty(id){
    RandNoLetter(length)
  }
  return id
}

//Returns A Random letters
func RandLetters(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBBNM")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  if !CheckifStringIsEmpty(id){
    RandLetters(length)
  }
  return id
}

//Returns a random number in string format
func RandNo(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("1234567890")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  if !CheckifStringIsEmpty(id){
    RandNo(length)
  }
  return id
}

func generateRefPrefix(timestamp time.Time) string {

	year := timestamp.Year()
	month := timestamp.Month()
	day := timestamp.Day()

	// Determine the year character ('A' for the current year)
	yearCharacter := 'A' + rune((year-2023)%26)

	// Determine the month character ('A' for January, 'B' for February, ..., 'L' for December)
	monthCharacter := 'A' + rune(month-1)

	// Determine the day character ('1' to '9' for the first 9 days, 'A' to 'Z' for the rest)
	var dayCharacter rune
	if day >= 1 && day <= 9 {
		dayCharacter = '0' + rune(day)
	} else {
		dayCharacter = 'A' + rune(day-10)
	}

	// Combine the characters into a string
	result := string([]rune{dayCharacter, monthCharacter, yearCharacter})

	return result
}

func formatTimeComponent(t int) string {

  formatted := fmt.Sprintf("%02d", t)

  return formatted
}

func generateTimeString() string {

  now := time.Now()
  hour := formatTimeComponent(now.Hour())
  minute := formatTimeComponent(now.Minute())
  second := formatTimeComponent(now.Second())
  millisecond := formatTimeComponent(now.Nanosecond() / 1000000)

  return hour + minute + second + millisecond
}

func GenerateTransactionRef() string {
  prefix := generateRefPrefix(time.Now())
  rest := generateTimeString()

  ref := fmt.Sprintf("%s%s", prefix, rest)

  return ref
}
