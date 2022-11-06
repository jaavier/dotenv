package example

import (
	"fmt"
	"os"

	"github.com/jaavier/dotenv"
)

func main() {
	if err := dotenv.Load(".env"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(os.Getenv("YOUR_SECRET_KEY_HERE"))
	}
}
