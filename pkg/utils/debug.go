package utils

import (
	"encoding/json"
	"fmt"
)

// ใช้ในการแสดงข้อมูลในรูปแบบที่สวยงาม
func Debug(data any) {
	bytes, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(bytes))
}

// ใช้ในการแปลงข้อมูลให้เป็น byte[]
func OutPut(data any) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}
