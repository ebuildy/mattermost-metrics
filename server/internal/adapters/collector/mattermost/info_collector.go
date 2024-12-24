package mattermost

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

// collectInfo gathers and returns core system information such as server version, edition, installation time, and SQL driver name.
func (c *Collector) collectInfo() *domain.MetricsDataInfo {
	apiSystem := c.gateway.PluginAPI.System
	apiConfiguration := c.gateway.PluginAPI.Configuration

	installationTime, _ := apiSystem.GetSystemInstallDate()

	SQLDriverName := ""

	if apiConfiguration.GetConfig() != nil {
		SQLDriverName = *apiConfiguration.GetConfig().SqlSettings.DriverName
	}

	return &domain.MetricsDataInfo{
		MattermostVersion:          apiSystem.GetServerVersion(),
		MattermostEdition:          c.getInfoEdition(),
		MattermostInstallationTime: installationTime,
		SQLDriverName:              SQLDriverName,
	}
}

func (c *Collector) getInfoEdition() string {
	if c.gateway.PluginAPI.System.GetLicense() == nil {
		return "free"
	}

	return "entreprise"
}
