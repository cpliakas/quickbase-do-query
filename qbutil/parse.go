package qbutil

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Split splits a string by "." or "," with space afterwards.
func Split(s string) []string {
	return regexp.MustCompile(`\s*[.,]\s*`).Split(s, -1)
}

// ParseFieldsOption parses a field list option.
func ParseFieldsOption(fieldsStr string) ([]int, error) {
	fieldsStr = strings.TrimSpace(fieldsStr)
	if fieldsStr == "" {
		return []int{}, nil
	}

	parts := Split(fieldsStr)
	fields := make([]int, len(parts))

	for key, part := range parts {
		fid, err := strconv.Atoi(part)
		if err != nil {
			// TODO: Invalid input error instead of generic.
			return []int{}, errors.New("invalid field ID")
		}
		fields[key] = fid
	}

	return fields, nil
}

// ParseSortOption parses a sort option and returns the field list, order list,
// and error respectively.
func ParseSortOption(sortStr string) ([]int, []string, error) {
	sortStr = strings.TrimSpace(sortStr)
	if sortStr == "" {
		return []int{}, []string{}, nil
	}

	parts := Split(sortStr)
	sort := make([]int, len(parts))
	order := make([]string, len(parts))

	re := regexp.MustCompile(`^([0-9]+)\s*(D|A|DESC|ASC)?$`)
	for k, part := range parts {

		match := re.FindStringSubmatch(part)
		if len(match) == 0 {
			// TODO: Invalid input error instead of generic.
			return []int{}, []string{}, errors.New("invalid input")
		}

		fid, err := strconv.Atoi(match[1])
		if err != nil {
			// TODO: Invalid input error instead of generic.
			return []int{}, []string{}, errors.New("invalid field ID")
		}
		sort[k] = fid

		// TODO: Validate whether match[2] exists?
		order[k] = match[2]
		if order[k] == "DESC" || order[k] == "D" {
			order[k] = "D"
		} else {
			order[k] = "A"
		}
	}

	return sort, order, nil
}
