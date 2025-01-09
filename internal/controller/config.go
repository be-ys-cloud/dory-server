package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

func Config(w http.ResponseWriter, _ *http.Request) {
	d, _ := json.Marshal(configurationToConfigurationDTO(configuration.Configuration).Features)
	_, _ = w.Write(d)
}

func configurationToConfigurationDTO(config structures.Configuration) structures.ConfigurationDTO {
	return structures.ConfigurationDTO{
		Features: structures.FeaturesDTO{
			DisableUnlock:                   config.Features.DisableUnlock,
			DisablePasswordUpdate:           config.Features.DisablePasswordUpdate,
			DisablePasswordReinitialization: config.Features.DisablePasswordReinitialization,
			DisableTOTP:                     config.Features.DisableTOTP,
		}}
}
