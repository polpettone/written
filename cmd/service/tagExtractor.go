package service

import "strings"

const TAG_PREFIX = "#"

func ExtractTags(content string) []string {
	fields := strings.Fields(content)

	tags := []string{}
	for _, f := range fields {
		if len(f) > 1 && strings.HasPrefix(f, TAG_PREFIX) {
			tags = append(tags, f)
		}
	}
	return tags
}
