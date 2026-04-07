package models

type BulkInboundMapping struct {
	BarcodeHeader string `json:"barcode_header" binding:"required"`
	NameHeader    string `json:"name_header" binding:"required"`
	QtyHeader     string `json:"qty_header" binding:"required"`
	PriceHeader   string `json:"price_header" binding:"required"`
}

type BulkInboundRequest struct {
	FileName    string             `json:"file_name"`
	Supplier    string             `json:"supplier" binding:"required"`
	TypeProduct string             `json:"type_product" binding:"required,oneof=reguler sticker"`
	Type        string             `json:"type" binding:"required,oneof=csv xlsx xls"`
	Mapping     BulkInboundMapping `json:"mapping" binding:"required"`
	Rows        [][]string         `json:"rows" binding:"required"`
	Headers     []string           `json:"headers" binding:"required"`
}
