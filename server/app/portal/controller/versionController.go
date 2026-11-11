package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/setting"

	"github.com/gin-gonic/gin"
)

type VersionInfoResp struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	StartTime string `json:"startTime"`
}

func GetVersion(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(VersionInfoResp{
		Name:      setting.Conf.Name,
		Version:   setting.Conf.Version,
		StartTime: setting.Conf.StartTime,
	})
}
