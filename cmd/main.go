package mn

import (
	"fmt"

	mn "github.com/Olin-Hydro/mother-nature/pkg"
)

func GetGardenConfig(store mn.Storage, client mn.HTTPClient, gardenId string) (mn.GardenConfig, error) {
	garden := mn.Garden{}
	req, err := store.CreateGardenReq(gardenId)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	err = mn.DecodeJson(&garden, res.Body)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	return garden.Config, nil
}

//nolint:unused
func main() {
	conf := mn.LoadConfigFromEnv()
	h, err := mn.NewHydrangea(conf.HydrangeaGardenURL, conf.HydrangeaRALogURL, conf.HydrangeaRAURL, conf.HydrangeaSensorLogURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := mn.Client
	gardenConfig, err := GetGardenConfig(h, client, conf.GardenId)
	if err != nil {
		fmt.Println(err)
		return
	}
	mn.Cache, err = mn.UpdateRACache(mn.Cache, gardenConfig.RAConfigs, client, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	commands, err := mn.CreateRACommands(gardenConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(commands)
	// TODO:
	// Send commands to hydrangea
}
