package http

import (
	"goaway/internal/models"
	"goaway/internal/utils"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var IpList map[uint32]models.IPRange
var RangeIpList utils.RangeSet

func LoadDatabase() (int, int64) {

	log.Println("Loading database, will take roughly 1 minute")
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

	log.Printf("Loaded single IP DB rows %d", len(records))
	IpList = records

	var rangeRecords utils.RangeSet
	// Now we load ranges in rangeset, the fun part!
	db.Table("ip_ranges").Where("end <> ?", 0).FindInBatches(&results, 10000, func(tx *gorm.DB, batch int) error {
		for _, result := range results {

			var md map[string]string
			md = make(map[string]string)
			md["List"] = result.List
			md["Timestamp"] = strconv.Itoa(int(result.Timestamp.Unix()))

			r := utils.Range{
				Low:      result.Start,
				High:     result.End,
				Metadata: md,
			}
			rangeRecords.AddRange(r)
		}
		return nil
	})

	RangeIpList = rangeRecords

	log.Printf("Loaded range records rows %d ", len(RangeIpList.Ranges))

	totalRecords := len(records) + len(rangeRecords.Ranges)
	totalTime := time.Now().Sub(startTime).Milliseconds()

	log.Printf("Loaded single All DB rows %d in %d ms", totalRecords, totalTime)

	return totalRecords, totalTime
}

func FindSuspiciousIp(ip uint32) (models.IPRange, bool) {

	// First check ranges
	r := RangeIpList.Contains(ip)
	if r != nil {
		log.Printf("Got range with metadata %s, %s", r.Metadata["Timestamp"], r.Metadata["List"])
		integerStamp, _ := strconv.Atoi(r.Metadata["Timestamp"])
		stamp := time.Unix(int64(integerStamp), 0)
		return models.IPRange{
			Start:     r.Low,
			End:       r.High,
			List:      r.Metadata["List"],
			Timestamp: stamp,
		}, true
	}

	// If not check ips
	val, ok := IpList[ip]
	return val, ok

}

func RunServer() {

	utils.LoadConfig()

	LoadDatabase()

	r := gin.Default()

	// Maximize parallellism
	runtime.GOMAXPROCS(runtime.NumCPU())

	v1 := r.Group("api/v1")
	{
		v1.GET("/ip/:ip", GetIp)
		v1.GET("/reload", GetReload)
	}

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
