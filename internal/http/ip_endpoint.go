package http

import (
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status    string `json:"status"`
	List      string `json:"description,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func ip2int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func GetIp(c *gin.Context) {

	log.Println("Inside Get IP %s", c.Param("ip"))
	ip := net.ParseIP(c.Param("ip"))

	if ip == nil {
		c.AbortWithError(400, errors.New("IP Not valid"))
		return
	}
	log.Printf("DWORD IP conversion %d", ip2int(ip))

	susIp, ok := FindSuspiciousIp(ip2int(ip))

	if !ok {
		c.JSON(200, Response{
			Status: "ok",
		})
	} else {
		c.JSON(200, Response{
			Status:    "ko",
			List:      susIp.List,
			Timestamp: susIp.Timestamp.Unix(),
		})
	}

}
