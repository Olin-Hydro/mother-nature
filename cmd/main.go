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

//nolint:unused
func main() {
	conf := pkg.LoadConfigFromEnv()
	h, err := pkg.NewHydrangea(conf.HydrangeaGardenURL, conf.HydrangeaRALogURL, conf.HydrangeaSensorLogURL)
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
	schedule, err := pkg.NewSchedule(gardenConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	schedBytes, err := pkg.EncodeSchedule(schedule)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(schedBytes)
	// TODO:
	// Send schedule to gardener
	// TODO:
	// Check conditions
	// Send commands to gardener if needed
}
