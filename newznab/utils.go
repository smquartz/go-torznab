package newznab

import (
	"strconv"
	"strings"
)

// stringifyCategories takes a []Category slice, and returns a comma
// delimited list of the category IDs
func stringifyCategories(in []Category) []string {
	var categories []string
	for _, v := range in {
		categories = append(categories, strconv.Itoa(int(v)))
	}
	return []string{strings.Join(categories, ",")}
}
