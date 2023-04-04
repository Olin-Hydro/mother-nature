package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	mn "github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
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
	} else if res.StatusCode != http.StatusOK {
		return config, fmt.Errorf("#getConfig: Garden response did not return status ok: %s", res.Status)
	}
	err = mn.DecodeJson(&garden, res.Body)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	if garden.Id == "" {
		return config, fmt.Errorf("#getConfig: Could not find garden with id: %s", gardenId)
	}
	req2, err := store.CreateConfigReq(garden.ConfigID)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	res2, err := client.Do(req2)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	} else if res2.StatusCode != http.StatusOK {
		return config, fmt.Errorf("#getConfig: Config response did not return status ok: %s", res2.Status)
	}
	err = mn.DecodeJson(&config, res2.Body)
	if err != nil {
		return config, fmt.Errorf("#getConfig: %e", err)
	}
	if config.Id == "" {
		return config, fmt.Errorf("#getConfig: Could not find config with id %s", garden.ConfigID)
	}
	return config, nil
}

func SendCommands(store mn.Storage, client mn.HTTPClient, commands []mn.Command) error {
	req, err := store.CreateCommandReq(commands)
	if err != nil {
		return fmt.Errorf("#SendCommands: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("#SendCommands: %e", err)
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("#SendCommands: Error sending commands to hydrangea: %v", res)
	}
	return nil
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) {
	conf := mn.LoadConfigFromEnv()
	h, err := mn.NewHydrangea(
		conf.HydrangeaGardenURL,
		conf.HydrangeaRALogURL,
		conf.HydrangeaRAURL,
		conf.HydrangeaSensorLogURL,
		conf.HydrangeaCommandURL,
		conf.HydrangeaConfigURL,
		conf.ApiKey,
	)
	if err != nil {
		log.Error(err)
		return
	}
	client := mn.Client
	gardenConfig, err := GetGardenConfig(h, client, conf.GardenId)
	if err != nil {
		log.Error(err)
		return
	}
	mn.Cache, err = mn.UpdateRACache(mn.Cache, gardenConfig.RAConfigs, conf.GardenId, client, h)
	if err != nil {
		log.Error(err)
		return
	}
	commands, err := mn.CreateRACommands(gardenConfig)
	if err != nil {
		log.Error(err)
		return
	}
	if len(commands) > 0 {
		cmdStr, _ := json.Marshal(commands)
		log.Info(string(cmdStr))
		err := SendCommands(h, client, commands)
		if err != nil {
			log.Error(err)
		}
		log.Info("Sent Commands")
	} else {
		log.Info("No commands to send")
	}
}
