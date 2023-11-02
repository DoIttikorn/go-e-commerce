package entities

import (
	"github.com/Doittikorn/go-e-commerce/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type ResponseImpl interface {
	Success(code int, data any) ResponseImpl
	Error(code int, traceId, msg string) ResponseImpl
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) ResponseImpl {
	return &Response{
		Context: c,
	}
}

// ใช้ในการ save StatusCode และ Data ที่ได้จากการทำงานของ handler ใส่ใน struct Response เพื่อให้สามารถเรียกใช้ได้ทุกที่
func (r *Response) Success(code int, data any) ResponseImpl {
	r.StatusCode = code
	r.Data = data
	logger.New(r.Context, &r.Data).Print().SaveToStorage()
	return r
}

// ใช้ในการ save StatusCode, traceId และ msg ที่ได้จากการทำงานของ handler ใส่ใน struct Response เพื่อให้สามารถเรียกใช้ได้ทุกที่
func (r *Response) Error(code int, traceId, msg string) ResponseImpl {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Msg:     msg,
	}

	logger.New(r.Context, &r.ErrorRes).Print().SaveToStorage()

	r.IsError = true
	return r
}

// ใช้ในการส่งข้อมูลจาก struct Response ไปยัง client
func (r *Response) Res() error {

	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes
		}
		return &r.Data
	}())

}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
