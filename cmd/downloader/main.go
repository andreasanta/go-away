package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	s "strings"
	"time"

	"goaway/internal/models"
	"goaway/internal/utils"

	"path/filepath"

	"github.com/go-git/go-git/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const chunk = 256

var uniqueIpMap map[uint32]bool
var total int = 0

func main() {

	utils.LoadConfig()

	log.Println("Cloning IP Set Repo")

	os.RemoveAll("./tmp/ipsets")
	_, err := git.PlainClone("./tmp/ipsets", false, &git.CloneOptions{

		URL:      os.Getenv("IPSET_GIT_REPO"),
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatalf("Error %s", err)
	}

	// Update timestamps
	log.Println("Updating timestamps...")
	cmd := exec.Command("./tmp/ipsets", "YES_I_AM_SURE_DO_IT_PLEASE")
	cmd.Wait()

	db := prepareDatabase()
	sqdb, _ := db.DB()
	defer sqdb.Close()

	uniqueIpMap = make(map[uint32]bool)
	loadIpsetsFiles(db)

}

func prepareDatabase() (db *gorm.DB) {

	os.Remove(os.Getenv("DB_FILE_PATH"))
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_PATH")), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Error),
	})

	//db.LogMode(true)

	bailOnError(err)

	// Create table if not exists, "Gorm" takes care of it
	db.AutoMigrate(&models.IPRange{})

	// Truncate it
	db.Exec("DELETE FROM ip_ranges")

	return db
}

func loadSingleFile(db *gorm.DB, path string) {

	// Open file and read line by line
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records []models.IPRange

	for scanner.Scan() {

		// Discard comments and empty lines
		line := s.TrimSpace(scanner.Text())

		if s.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Parse CIDR
		var start, end uint32
		if s.Contains(line, "/") {
			start, end = HostRange(line)
		} else {
			ip := net.ParseIP(line)
			start = ip2int(ip)
			end = 0
		}

		finfo, _ := file.Stat()
		r := models.IPRange{
			Start:     start,
			End:       end,
			List:      filepath.Base(file.Name()),
			Timestamp: finfo.ModTime(),
		}

		if len(records) == chunk {

			total += chunk

			// log.Printf("Flushing inserts at total %d", total)
			db.Create(&records)
			records = nil

		} else if !(uniqueIpMap[start]) {
			records = append(records, r)
			uniqueIpMap[start] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func ip2int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func loadIpsetsFiles(db *gorm.DB) {

	// Read all the ipsets file in the tmp subdir
	matches, _ := filepath.Glob("./tmp/ipsets/*.*set")

	for _, file := range matches {
		initialTime := time.Now()
		log.Printf("Loading file %s", file)
		loadSingleFile(db, file)
		log.Printf("Loaded in %dms", time.Now().Sub(initialTime).Milliseconds())
	}

	log.Printf("Total processed records %d", total)

}

func bailOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Below courtesy of https://gist.github.com/kotakanbe/d3059af990252ba89a82

func HostRange(cidr string) (start uint32, end uint32) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}

	ones, _ := ipnet.Mask.Size()
	ip_nums := uint32(math.Trunc(math.Pow(2, float64(32-ones))))
	//log.Printf("Found CIDR %s => Total IPs %d", ip, ip_nums)

	start = ip2int(ip)
	end = start + ip_nums

	return start, end
}
