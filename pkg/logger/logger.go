package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Doittikorn/go-e-commerce/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

// ใช้ในการสร้าง logger เพื่อ save ข้อมูลต่างๆ จาก context library นั้นๆ
func New(c *fiber.Ctx, res any) LoggerImpl {
	log := &logger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

// print ข้อมูลที่อยู่ใน struct logger
func (l *logger) Print() LoggerImpl {
	utils.Debug(l)
	return l
}

// TODO: ทำให้การที่อยู่การ save log ไม่ขึ้นอยู่กับฟังก์ชันนี้
//
// save ข้อมูลลง store ที่เรากำหนดไว้
func (l *logger) SaveToStorage() {
	data := utils.OutPut(l) // เอาข้อมูลใน struct logger มาเป็น  byte[]

	fileName := fmt.Sprintf("./assets/logs/logger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

// เอาข้อมูล query param มาใส่ struct logger
func (l *logger) SetQuery(c *fiber.Ctx) {
	var query any
	if err := c.QueryParser(&query); err != nil {
		log.Printf("error parsing body: %v", err)
	}
	l.Query = query
}

// เอาข้อมูล body มาใส่ struct logger
func (l *logger) SetBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("error parsing body: %v", err)
	}
	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body
	}
}

// นำ response data มาใส่ใน struct logger
func (l *logger) SetResponse(res any) {
	l.Response = res
}
