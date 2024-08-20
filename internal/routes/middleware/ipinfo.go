package middleware

import (
	"cmp"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/maxminddb-golang"

	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/routes/common"
)

func SetIPInfo(db *maxminddb.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := common.GetRealIP(c)
		ipInfo := &request.IPInfo{IP: ip}
		result, err := getIPCountryShort(db, ip)
		if err != nil {
			slog.Error("middleware SetIPInfo get client ip err", "ip", ip, "err", err)
		}
		ipInfo.CountryShort = cmp.Or(result, "-")
		c.Set(common.IpInfoKey, ipInfo)
	}
}

func getIPCountryShort(db *maxminddb.Reader, clientIP string) (string, error) {
	ip := net.ParseIP(clientIP)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err := db.Lookup(ip, &record)
	if err != nil {
		return "", err
	}

	return record.Country.ISOCode, nil
}
