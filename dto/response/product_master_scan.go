package dto

type ProductMasterScanResponse struct {
	ID               string  `json:"id"`
	BarcodeWarehouse string  `json:"barcode_warehouse"`
	NameWarehouse    string  `json:"name_warehouse"`
	ItemWarehouse    int     `json:"item_warehouse"`
	Location         string  `json:"location"`
	RackStagingID    *string `json:"rack_staging_id"`
}
