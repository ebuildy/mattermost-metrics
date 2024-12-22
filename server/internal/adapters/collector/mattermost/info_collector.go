package mattermost

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

func (c *Collector) collectInfo() *domain.MetricsDataInfo {
	apiSystem := c.gateway.PluginAPI.System
	apiConfiguration := c.gateway.PluginAPI.Configuration

	installationTime, _ := apiSystem.GetSystemInstallDate()

	return &domain.MetricsDataInfo{
		MattermostVersion:          apiSystem.GetServerVersion(),
		MattermostEdition:          c.getInfoEdition(),
		MattermostInstallationTime: installationTime,
		SQLDriverName:              *apiConfiguration.GetConfig().SqlSettings.DriverName,
	}
}

func (c *Collector) getInfoEdition() string {
	if c.gateway.PluginAPI.System.GetLicense() == nil {
		return "free"
	}

	return "entreprise"
}
