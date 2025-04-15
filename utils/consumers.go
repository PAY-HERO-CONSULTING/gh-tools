package utils

import "strings"

func ExtractConsumer(queue string) string {
	return extractPrefix(queue)
}

func ExtractGroup(topic string) string {
	return extractPrefix(topic)
}

func extractPrefix(value string) string {
	if !strings.Contains(value, ":") {
		return ""
	}

	return strings.Split(value, ":")[0]
}

func ExtractTopic(queue string) string {
	if !strings.Contains(queue, ":") {
		return ""
	}

	return strings.Split(queue, ":")[1]
}
