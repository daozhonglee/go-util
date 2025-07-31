package api

import "testing"

func TestSuccess(t *testing.T) {
	response, err := Success("操作成功")
	if err != nil {
		t.Errorf("Expected no error for success, got %v", err)
	}
	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}
	if response.Message != "操作成功" {
		t.Errorf("Expected '操作成功', got %s", response.Message)
	}
}

func TestError(t *testing.T) {
	response, err := Error("操作失败", 1001)
	if err == nil {
		t.Error("Expected error for error response")
	}
	if response.Code != 1001 {
		t.Errorf("Expected code 1001, got %d", response.Code)
	}
}

func TestSuccessWithData(t *testing.T) {
	data := map[string]string{"key": "value"}
	response, err := SuccessWithData("成功", data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestSuccessWithPage(t *testing.T) {
	data := []string{"item1", "item2"}
	page := Pagination{Offset: 0, Limit: 10, Total: 2}
	response, err := SuccessWithPage("成功", data, page)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Page.Total != 2 {
		t.Errorf("Expected total 2, got %d", response.Page.Total)
	}
}
