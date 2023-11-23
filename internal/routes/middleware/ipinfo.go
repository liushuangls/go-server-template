package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ip2location/ip2location-go/v9"

	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/routes/common"
)

func SetIPInfo(ipCli *ip2location.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := common.GetRealIP(c)
		ipInfo := &request.IPInfo{IP: ip}
		result, err := ipCli.Get_country_short(ip)
		if err != nil {
			slog.Error("middleware SetIPInfo get client ip err", "ip", ip, "err", err)
		}
		ipInfo.CountryShort = result.Country_short
		c.Set(common.IpInfoKey, ipInfo)
	}
}
