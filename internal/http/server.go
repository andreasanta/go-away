package http

import (
	"goaway/internal/models"
	"goaway/internal/utils"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var IpList map[uint32]models.IPRange

func LoadDatabase() (int, int64) {

	log.Println("Loading database, will take a long while")
	startTime := time.Now()

	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_PATH")), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	records := make(map[uint32]models.IPRange)

	if err != nil {
		log.Fatal(err)
	}

	var results []models.IPRange
	db.Table("ip_ranges").Where("end = ?", 0).FindInBatches(&results, 10000, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			// batch processing found records
			records[result.Start] = result
		}
		return nil
	})

	log.Printf("Loaded DB rows %d in %dms", len(records))
	IpList = records

	totalRecords := len(records)
	totalTime := time.Now().Sub(startTime).Milliseconds()

	return totalRecords, totalTime
}

func FindSuspiciousIp(ip uint32) (models.IPRange, bool) {

	log.Printf("Loaded DB rows %d", len(IpList))

	val, ok := IpList[ip]
	return val, ok
}

func RunServer() {

	utils.LoadConfig()

	LoadDatabase()

	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.GET("/ip/:ip", GetIp)
		v1.GET("/reload", GetReload)
	}

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
