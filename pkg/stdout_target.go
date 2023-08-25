package pkg

import (
	"fmt"
)

// Print the CSVSource struct to stdout as Redis commands
func (c *CSVSource) Print(prefix string) {
	for ui, row := range c.Rows {
		for i, value := range row {
			keyname := fmt.Sprintf("%s%s-%v", prefix, c.Header[i], ui)
			fmt.Printf("SET %s %s\n", keyname, value)
		}
	}
}
