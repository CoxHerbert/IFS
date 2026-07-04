package skills

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Unit   string
}

var dimensionPattern = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(?:cm)?\s*[*x\x{00D7}]\s*(\d+(?:\.\d+)?)\s*(?:cm)?\s*[*x\x{00D7}]\s*(\d+(?:\.\d+)?)\s*(?:cm)?`)

func ParseDimensions(message string) (Dimensions, bool, error) {
	match := dimensionPattern.FindStringSubmatch(strings.TrimSpace(message))
	if len(match) != 4 {
		return Dimensions{}, false, nil
	}

	length, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return Dimensions{}, true, err
	}
	width, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return Dimensions{}, true, err
	}
	height, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		return Dimensions{}, true, err
	}
	if length <= 0 || width <= 0 || height <= 0 {
		return Dimensions{}, true, errors.New("dimensions must be greater than 0")
	}

	return Dimensions{Length: length, Width: width, Height: height, Unit: "cm"}, true, nil
}
