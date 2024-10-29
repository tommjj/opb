package browser

import (
	"fmt"

	"github.com/tommjj/go-opb/internal/config"
	"github.com/tommjj/go-opb/internal/utils"
)

func Reset(conf *config.ConfigStore) error {
	fmt.Println("Set your browser path")
	path := utils.GetInput()
	return conf.SetBrowser(path)
}
