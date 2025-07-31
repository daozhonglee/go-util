// Package api 提供标准化的API响应结构
package api

import (
	"errors"
)

// Success 创建成功响应 (code=0)
func Success(message string) (*Response, error) {
	return NewResponse(message, 0)
}

// Error 创建错误响应
func Error(message string, code int) (*Response, error) {
	return NewResponse(message, code)
}

// SuccessWithData 创建带数据的成功响应
func SuccessWithData(message string, data interface{}) (*DataResponse, error) {
	return NewDataResponse(message, 0, data)
}

// ErrorWithData 创建带数据的错误响应
func ErrorWithData(message string, code int, data interface{}) (*DataResponse, error) {
	return NewDataResponse(message, code, data)
}

// SuccessWithPage 创建分页成功响应
func SuccessWithPage(message string, data interface{}, page Pagination) (*PageResponse, error) {
	return NewPageResponse(message, 0, data, page)
}

// NewResponse 创建基础响应
func NewResponse(message string, code int) (*Response, error) {
	response := &Response{
		Message: message,
		Code:    code,
	}
	if code != 0 {
		return response, errors.New(message)
	}
	return response, nil
}

// NewDataResponse 创建带数据的响应
func NewDataResponse(message string, code int, data interface{}) (*DataResponse, error) {
	response := &DataResponse{
		Message: message,
		Code:    code,
		Data:    data,
	}
	if code != 0 {
		return response, errors.New(message)
	}
	return response, nil
}

// NewPageResponse 创建分页响应
func NewPageResponse(message string, code int, data interface{}, page Pagination) (*PageResponse, error) {
	response := &PageResponse{
		Message: message,
		Code:    code,
		Data:    data,
		Page:    page,
	}
	if code != 0 {
		return response, errors.New(message)
	}
	return response, nil
}

// NewResponseNoError 创建不返回错误的响应
func NewResponseNoError(message string, code int) (*Response, error) {
	return &Response{
		Message: message,
		Code:    code,
	}, nil
}

// Response 基础响应结构
type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// DataResponse 数据响应结构
type DataResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Page    Pagination  `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}
