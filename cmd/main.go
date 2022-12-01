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

func UpdateActuationTimes(actuationTimes pkg.RActuationTimes, raConfigs []pkg.RAConfig, client pkg.HTTPClient, store pkg.Storage) (pkg.RActuationTimes, error) {
	for i := 0; i < len(raConfigs); i++ {
		raConfig := raConfigs[i]
		if _, ok := actuationTimes.Times[raConfig.Id]; ok {
			continue
		}
		req, err := store.CreateRALogsReq(raConfig.Id, "1")
		if err != nil {
			return actuationTimes, fmt.Errorf("InitActuationTimes: %e", err)
		}
		res, err := client.Do(req)
		if err != nil {
			return actuationTimes, fmt.Errorf("InitActuationTimes: %e", err)
		}
		raLog := pkg.RALog{}
		err = pkg.DecodeJson(&raLog, res.Body)
		if err != nil {
			return actuationTimes, fmt.Errorf("InitActuationTimes: %e", err)
		}
		actTime, err := pkg.StrToTime(raLog.CreatedAt)
		if err != nil {
			return actuationTimes, fmt.Errorf("InitActuationTimes: %e", err)
		}
		actuationTimes.Times[raConfig.Id] = actTime
	}
	return actuationTimes, nil
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
	err = SendSchedule(gardenConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	pkg.ActuationTimes, err = UpdateActuationTimes(pkg.ActuationTimes, gardenConfig.ReactiveActuators, client, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO:
	// Check conditions
	// Send commands to gardener if needed
}
