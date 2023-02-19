package mn

import (
	"fmt"

	mn "github.com/Olin-Hydro/mother-nature/pkg"
	_ "github.com/joho/godotenv/autoload"
)

func GetGardenConfig(store mn.Storage, client mn.HTTPClient, gardenId string) (mn.GardenConfig, error) {
	garden := mn.Garden{}
	config := mn.GardenConfig{}
	req, err := store.CreateGardenReq(gardenId)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	err = mn.DecodeJson(&garden, res.Body)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	req2, err := store.CreateConfigReq(garden.ConfigID)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	fmt.Println(req2.URL)
	res2, err := client.Do(req2)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	err = mn.DecodeJson(&config, res2.Body)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	return config, nil
}

//nolint:unused
func main() {
	conf := mn.LoadConfigFromEnv()
	h, err := mn.NewHydrangea(
		conf.HydrangeaGardenURL,
		conf.HydrangeaRALogURL,
		conf.HydrangeaRAURL,
		conf.HydrangeaSensorLogURL,
		conf.HydrangeaCommandURL,
		conf.HydrangeaConfigURL,
	)
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
	mn.Cache, err = mn.UpdateRACache(mn.Cache, gardenConfig.RAConfigs, conf.GardenId, client, h)
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
