package utils

import "strings"

func IsWeb3Currency(currency string) bool {
	switch strings.ToLower(currency) {
	case "ckes":
		return true
	case "cusd":
		return true
	case "usdc":
		return true
	case "ceur":
		return true
	case "lsk":
		return true
	case "celo":
		return true
	default:
		return false
	}
}

func GetWeb3ProviderFromToken(token string) string {

	switch strings.ToLower(token) {
	case "ckes":
		return "celo"
	case "cusd":
		return "celo"
	case "usdc":
		return "celo"
	case "usdt":
		return "lisk"
	case "ceur":
		return "celo"
	case "lsk":
		return "lisk"
	case "celo":
		return "celo"
	default:
		return "lisk"
	}
}
