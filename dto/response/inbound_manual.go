package dto

import "time"

// StatusBreakdown untuk breakdown per status
type StatusBreakdown struct {
	TotalItem  int     `json:"total_item"`
	TotalPrice float64 `json:"total_price"`
	Persentase float64 `json:"persentase"`
}

// InboundDocumentSummary untuk API Summary
type InboundDocumentSummary struct {
	Code        string          `json:"code"`
	FileName    string          `json:"file_name"`
	FileItem    int             `json:"file_item"`
	FilePrice   float64         `json:"file_price"`
	Good        StatusBreakdown `json:"good"`
	Damaged     StatusBreakdown `json:"damaged"`
	Abnormal    StatusBreakdown `json:"abnormal"`
	Non         StatusBreakdown `json:"non"`
	Discrepancy StatusBreakdown `json:"discrepancy"`
	File        StatusBreakdown `json:"file"`
	Status      string          `json:"status"`
}

// InboundProductPendingResponse untuk list barang
type InboundProductPendingResponse struct {
	ID          string     `json:"id"`
	Barcode     string     `json:"barcode"`
	Name        string     `json:"name"`
	Item        int        `json:"item"`
	ItemGood    int        `json:"item_good"`
	ItemDamaged int        `json:"item_damaged"`
	Price       float64    `json:"price"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
	DateScanned *time.Time `json:"date_scanned"`
	IsSKU       bool       `json:"is_sku"`
}
