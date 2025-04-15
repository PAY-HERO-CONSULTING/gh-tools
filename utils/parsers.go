package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/custom_types"
	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
	"github.com/nyaruka/phonenumbers"
)

var (
	nonDigitRe        = regexp.MustCompile("[^0-9eE+]")
	kenyanRe          = regexp.MustCompile("^(0)?(1|7)[0-9]{8}$")
	kenyanCountryCode = "254"
)

func ParsePhoneNumber(phoneNumberStr string) (custom_types.PhoneNumber, error) {
	phoneNumberStr = nonDigitRe.ReplaceAllString(phoneNumberStr, "")

	phoneNumberStr = sanitizeScientificNotation(phoneNumberStr)

	if len(phoneNumberStr) < 7 {
		return custom_types.PhoneNumber{}, apperrs.NewErrorWithType(
			fmt.Errorf("number (%v) is less than 7 digits", phoneNumberStr),
			apperrs.BadRequest,
		).LogInfoMessage("gh-tools", "number is less than 7 digits", fmt.Errorf("number (%v) is less than 7 digits", phoneNumberStr))
	}

	if len(phoneNumberStr) > 15 {
		return custom_types.PhoneNumber{}, apperrs.NewErrorWithType(
			fmt.Errorf("number (%v) is greater than 15 digits", phoneNumberStr),
			apperrs.BadRequest,
		).LogInfoMessage("gh-tools", "number is greater than 15 digits", fmt.Errorf("number (%v) is greater than 15 digits", phoneNumberStr))
	}

	countryCode, _ := GetDialCode(phoneNumberStr)
	var err error

	if countryCode == "" || countryCode == kenyanCountryCode || hasKenyanPrefix(phoneNumberStr) {
		countryCode, phoneNumberStr, err = cleanKenyanNumber(phoneNumberStr)
		if err != nil {
			return custom_types.PhoneNumber{}, err
		}
	}

	num, err := phonenumbers.Parse(phoneNumberStr, countryDialCodesMap[countryCode])
	if err != nil {
		return custom_types.PhoneNumber{}, apperrs.NewErrorWithType(
			fmt.Errorf("invalid phone number"),
			apperrs.BadRequest,
		).LogInfoMessage("gh-tools", "invalid phone number", fmt.Errorf("(%v) is not a valid number", phoneNumberStr))
	}

	ok := phonenumbers.IsValidNumber(num)
	if !ok {
		return custom_types.PhoneNumber{}, apperrs.NewErrorWithType(
			fmt.Errorf("invalid phone number"),
			apperrs.BadRequest,
		)
	}

	countryCode = strconv.FormatInt(int64(num.GetCountryCode()), 10)

	number := strconv.FormatInt(int64(num.GetNationalNumber()), 10)

	if num.GetItalianLeadingZero() {
		number = "0" + number
	}

	countryCode = "+" + countryCode
	number = strings.TrimPrefix(number, countryCode)

	phoneNumber := custom_types.PhoneNumber{
		DialCode: null.NullValue(countryCode),
		Number:   null.NullValue(number),
	}

	return phoneNumber, nil
}

func cleanKenyanNumber(phoneNumber string) (string, string, error) {
	newPhoneNumber := phoneNumber

	if len(newPhoneNumber) < 5 {
		return "", phoneNumber, apperrs.NewBadRequest(
			fmt.Sprintf("number (%v) is less than 5 digits", newPhoneNumber),
		)
	}

	if newPhoneNumber[0:1] == "+" {
		newPhoneNumber = newPhoneNumber[1:]
	}

	if newPhoneNumber[0:3] == "254" {
		newPhoneNumber = newPhoneNumber[3:]
	}

	if newPhoneNumber[0:1] == "0" {
		newPhoneNumber = newPhoneNumber[1:]
	}

	if !kenyanRe.MatchString(newPhoneNumber) {
		return "", phoneNumber, apperrs.NewBadRequest(
			fmt.Sprintf("number (%v) is not valid kenyan number", newPhoneNumber),
		).LogInfoMessage("gh-tools", "invalid kenyan number", fmt.Errorf("number (%v) is not valid kenyan number", newPhoneNumber))
	}

	newPhoneNumber = "+254" + newPhoneNumber

	return kenyanCountryCode, newPhoneNumber, nil
}

func GetCountryCode(phoneNumber custom_types.PhoneNumber) (string, error) {
	dialCode := null.ValueFromNull[string](phoneNumber.DialCode)
	countryCode, ok := countryDialCodesMap[dialCode[1:]]
	if !ok {
		return "", apperrs.NewError(errors.New("invalid country code"))
	}

	return countryCode, nil
}

func GetDialCode(phoneNumber string) (string, error) {
	valid := false
	countryCode := ""

	if phoneNumber[0:1] == "+" {
		phoneNumber = phoneNumber[1:]
	}

	for i := 0; i < 3; i++ {
		countryCode = phoneNumber[0 : i+1]
		if _, ok := countryDialCodesMap[countryCode]; !ok {
			continue
		}

		valid = true
		break
	}

	if !valid {
		return "", apperrs.NewErrorWithType(
			errors.New("invalid country code"),
			apperrs.BadRequest,
		).LogInfoMessage("gh-tools", "invalid country code", errors.New("entered country code is invalid"))
	}

	return countryCode, nil
}

func sanitizeScientificNotation(phoneNumber string) string {
	phoneNumber = strings.ToLower(phoneNumber)
	index := strings.Index(phoneNumber, "e")
	if index > -1 {
		phoneNumber = phoneNumber[:strings.Index(phoneNumber, "e")]
	}

	return phoneNumber
}

func hasKenyanPrefix(phoneNumberStr string) bool {
	localPrefixes := "124567"
	return strings.Contains(localPrefixes, phoneNumberStr[0:1]) && len(phoneNumberStr) < 10
}
