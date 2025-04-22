package main

import (
	"fmt"

	"github.com/johnstontu/rss_feed/internal/config"
)

func main() {
	cfg, err := config.Read(".gatorconfig.json")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%+v\n", cfg)

}
