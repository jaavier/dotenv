package dotenv

import (
	"io/ioutil"
	"os"
	"strings"
)

func Load(file string) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	} else {
		var variables = (strings.Split(string(fileContent), "\n"))
		for _, variable := range variables {
			var key, value string
			var content = strings.SplitAfter(variable, "=")
			if len(content) == 2 {
				key = strings.Replace(content[0], "=", "", -1)
				value = content[1]
			} else if len(content) > 2 {
				key = strings.Replace(content[0], "=", "", -1)
				value = strings.Join(content[1:], "")
			}
			if len(key) > 0 && len(value) > 0 {
				os.Setenv(key, value)
			}
		}
		return nil
	}
}
