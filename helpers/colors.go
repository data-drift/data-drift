package helpers

import (
	"crypto/md5"
	"fmt"
)

func GetColorFromString(s string) string {
	// Hash the input string using MD5
	hash := md5.Sum([]byte(s))

	// Convert the first 3 bytes of the hash to a 24-bit integer
	r := int(hash[0]) << 16
	g := int(hash[1]) << 8
	b := int(hash[2])

	// Combine the RGB values into a single integer
	color := r | g | b

	// Convert the integer to a CSS color string
	return fmt.Sprintf("#%06x", color)
}
