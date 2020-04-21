package syllabus

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/chez-shanpu/syllabus-bot/pkg/db"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
)

type BoardInfo struct {
	gorm.Model
	CheckSum  string
	Title     string
	Date      string
	ClassName string
	Content   string
}

func InitBoardDB(i db.DBInfo) error {
	m := BoardInfo{}
	//err := db.InitDB(i, m)
	d, err := db.OpenDB(i)
	if err != nil {
		return err
	}
	defer d.Close()

	// Migrate the schema
	d.AutoMigrate(&m)
	return nil
}

func CreateBoardRecord(i db.DBInfo, bi BoardInfo) error {
	//tableName := "board_infos"
	//err := db.CreateRecord(i, bi, tableName)
	d, err := db.OpenDB(i)
	if err != nil {
		return err
	}
	defer d.Close()

	d.Create(&bi)
	log.Printf("[INFO] Create Record: %v", bi)
	return nil
}

func GetBoardRecord(i db.DBInfo, checkSum string) (*BoardInfo, error) {
	var bi BoardInfo

	d, err := db.OpenDB(i)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	d.Where("check_sum = ?", checkSum).First(&bi)

	return &bi, nil
}

func GetBoardInfoList() *[]BoardInfo {
	var bs []BoardInfo
	str := ""
	b := BoardInfo{}

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("syllabus.naist.jp"),
	)

	c.OnHTML("table.tbl01", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			str += e.ChildText("td")
			str += "|"
			tmp := e.ChildText("td.w15pr")
			if tmp != "" {
				b.Date = tmp
			}
		})
		h := md5.Sum([]byte(str))
		b.CheckSum = hex.EncodeToString(h[:])
		slice := strings.Split(str, "|")
		str = ""
		if slice[0] == "" {
			return
		}
		b.Title = strings.Replace(slice[0], b.Date, "", 1)
		b.ClassName = slice[1]
		b.Content = slice[2]
		bs = append(bs, b)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://syllabus.naist.jp/informations/preview_list")

	return &bs
}
