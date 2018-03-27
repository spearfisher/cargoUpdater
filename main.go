package main

import (
	"fmt"
	"time"

	"github.com/spearfisher/cargoUpdater/della"
	"github.com/spearfisher/cargoUpdater/utils"
)

func main() {
	defer utils.Logger.Println("Application stopped")

	dellaClient := della.NewDellaClient(utils.AppConfig.Login, utils.AppConfig.Password)
	periodicJob(dellaClient)

	period := time.Duration(utils.AppConfig.Period) * time.Minute
	for t := range time.NewTicker(period).C {
		currentHout := t.Hour()
		if (currentHout >= utils.AppConfig.Start) && (currentHout < utils.AppConfig.Stop) {
			periodicJob(dellaClient)
		}
	}
}

func periodicJob(dellaClient *della.Client) {
	cargoData, err := dellaClient.GetList()
	if err != nil {
		utils.Logger.Println(err)
	}

	if len(cargoData.Ids) > 0 {
		fmt.Println(fmt.Sprintf("Updating next entyties: %s", cargoData.Ids))
		dellaClient.RefreshCargos(cargoData)
	}
}
