package mn

import (
	"fmt"

	"github.com/Olin-Hydro/mother-nature/pkg"
)

func GetGardenConfig(store pkg.Storage, client pkg.HTTPClient, gardenId string) (pkg.GardenConfig, error) {
	garden := pkg.Garden{}
	req, err := store.CreateGardenReq(gardenId)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	err = pkg.DecodeJson(&garden, res.Body)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	return garden.Config, nil
}

func SendSchedule(gardenConfig pkg.GardenConfig) error {
	schedule, err := pkg.NewSchedule(gardenConfig)
	if err != nil {
		return fmt.Errorf("#SendSchedule: %e", err)
	}
	schedBytes, err := pkg.EncodeSchedule(schedule)
	if err != nil {
		return fmt.Errorf("#SendSchedule: %e", err)
	}
	fmt.Println(schedBytes)
	// TODO: Send the schedule to the gardener
	return nil
}

//nolint:unused
func main() {
	conf := pkg.LoadConfigFromEnv()
	h, err := pkg.NewHydrangea(conf.HydrangeaGardenURL, conf.HydrangeaRALogURL, conf.HydrangeaRAURL, conf.HydrangeaSensorLogURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := pkg.Client
	gardenConfig, err := GetGardenConfig(h, client, conf.GardenId)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = SendSchedule(gardenConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	pkg.Cache, err = pkg.UpdateRACache(pkg.Cache, gardenConfig.RAConfigs, client, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO:
	// Check conditions
	// Send commands to gardener if needed
}
