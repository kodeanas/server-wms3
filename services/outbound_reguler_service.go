package services

import (
	"fmt"
	"time"
	"wms/models"
	"wms/repositories"
)

type OutboundRegulerService interface {
	GetBuyers() interface{}
	GetBuyerClassInfo(id string) interface{}
	ScanProduct(ctx interface{}) interface{}
	AddProduct(ctx interface{}) interface{}
	DeleteProduct(id string) interface{}
	AddDiscount(ctx interface{}) interface{}
	UpdateTax(ctx interface{}) interface{}
	UpdateBox(ctx interface{}) interface{}
	CompleteOrder(ctx interface{}) interface{}
	GetOrderDetail(orderID string) interface{}
}

type outboundRegulerService struct {
	buyerRepo         repositories.BuyerRepository
	classRepo         repositories.ClassRepository
	orderRepo         repositories.OrderRepository
	productOrderRepo  repositories.ProductOrderRepository
	discountOrderRepo repositories.DiscountOrderRepository
	categoryRepo      repositories.CategoryRepository
	productMasterRepo repositories.ProductMasterRepository
}

func NewOutboundRegulerService(
	buyerRepo repositories.BuyerRepository,
	classRepo repositories.ClassRepository,
	orderRepo repositories.OrderRepository,
	productOrderRepo repositories.ProductOrderRepository,
	discountOrderRepo repositories.DiscountOrderRepository,
	categoryRepo repositories.CategoryRepository,
	productMasterRepo repositories.ProductMasterRepository,
) OutboundRegulerService {
	return &outboundRegulerService{
		buyerRepo:         buyerRepo,
		classRepo:         classRepo,
		orderRepo:         orderRepo,
		productOrderRepo:  productOrderRepo,
		discountOrderRepo: discountOrderRepo,
		categoryRepo:      categoryRepo,
		productMasterRepo: productMasterRepo,
	}
}

func (s *outboundRegulerService) GetBuyers() interface{} {
	buyers, err := s.buyerRepo.List()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return buyers
}
func (s *outboundRegulerService) GetBuyerClassInfo(id string) interface{} {
	// Ambil buyer
	buyer, err := s.buyerRepo.GetByID(id)
	if err != nil {
		return map[string]interface{}{"error": "buyer not found"}
	}
	// Ambil class buyer
	class, err := s.classRepo.GetByID(buyer.ClassID)
	if err != nil {
		return map[string]interface{}{"error": "class not found"}
	}

	// Siapkan response
	resp := map[string]interface{}{
		"buyer_id":       buyer.ID.String(),
		"buyer_name":     buyer.Name,
		"class_name":     class.Name,
		"class_discount": float64(class.Disc),
	}

	// Cek penurunan class jika buyer pernah transaksi
	lastOrder, err := s.orderRepo.GetLastDoneByBuyer(buyer.ID.String())
	transaksiSelesai, _ := s.orderRepo.CountDoneByBuyer(buyer.ID.String())
	if err == nil && lastOrder != nil {
		daysSince := int((NowDate().Sub(lastOrder.UpdatedAt)).Hours() / 24)
		if daysSince > class.Week*7 && class.Iteration > 1 {
			// Turunkan class buyer satu tingkat
			prevClass, err := s.classRepo.GetPrevByIteration(class.Iteration)
			if err == nil && prevClass != nil {
				// Update class buyer (opsional: update DB juga jika mau)
				resp["class_name"] = prevClass.Name
				resp["class_discount"] = float64(prevClass.Disc)
				resp["class_demotion_note"] = fmt.Sprintf("Class turun ke %s karena tidak transaksi lebih dari %d minggu", prevClass.Name, class.Week)
			}
		}
	}

	// Cek next class
	nextClass, err := s.classRepo.GetNextByIteration(class.Iteration)
	if err == nil && nextClass != nil {
		kurang := nextClass.MinOrder - transaksiSelesai
		if kurang <= 0 {
			kurang = 0
		}
		resp["next_class_note"] = "Next class " + nextClass.Name + ": perlu " + itoa(kurang) + " transaksi lagi setelah transaksi ini selesai."
	} else {
		resp["next_class_note"] = "Sudah di class tertinggi."
	}
	return resp
}

func NowDate() (now time.Time) {
	now = time.Now()
	return
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func (s *outboundRegulerService) ScanProduct(ctx interface{}) interface{} {
	// ctx diasumsikan *echo.Context, parsing request
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	barcode, _ := req["barcode"].(string)
	orderID, _ := req["order_id"].(string)
	buyerID, _ := req["buyer_id"].(string)

	// Validasi produk
	product, err := s.productMasterRepo.FindByBarcodeWarehouse(barcode)
	if err != nil || product == nil || product.Location != "display" {
		return map[string]interface{}{"error": "Produk tidak ditemukan atau bukan location display"}
	}

	// Jika order_id kosong, buat order baru
	var order *models.Order
	if orderID == "" {
		order = &models.Order{
			BuyerID: buyerID,
			Type:    "regular",
			Status:  "progress",
		}
		if err := s.orderRepo.Create(order); err != nil {
			return map[string]interface{}{"error": "Gagal membuat order"}
		}
	} else {
		order, err = s.orderRepo.GetByID(orderID)
		if err != nil {
			return map[string]interface{}{"error": "Order tidak ditemukan"}
		}
	}

	// Ambil diskon dari category
	discount := 0.0
	if product.CategoryID != nil {
		cat, err := s.categoryRepo.GetByID(*product.CategoryID)
		if err == nil && cat.Discount != nil {
			discount = float64(*cat.Discount)
		}
	}

	// Tambah product order
	po := &models.ProductOrder{
		OrderID:   order.ID.String(),
		ProductID: product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		Discount:  discount,
	}
	if err := s.productOrderRepo.Create(po); err != nil {
		return map[string]interface{}{"error": "Gagal menambah produk ke order"}
	}

	return map[string]interface{}{"order_id": order.ID.String(), "product_order_id": po.ID.String()}
}
func (s *outboundRegulerService) AddProduct(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	productID, _ := req["product_id"].(string)
	qty, ok := req["qty"].(int)
	if !ok || qty < 1 {
		qty = 1
	}
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	product, err := s.productMasterRepo.FindByID(productID)
	if err != nil {
		return map[string]interface{}{"error": "Produk tidak ditemukan"}
	}
	discount := 0.0
	if product.CategoryID != nil {
		cat, err := s.categoryRepo.GetByID(*product.CategoryID)
		if err == nil && cat.Discount != nil {
			discount = float64(*cat.Discount)
		}
	}
	var lastID string
	for i := 0; i < qty; i++ {
		po := &models.ProductOrder{
			OrderID:   order.ID.String(),
			ProductID: product.ID.String(),
			Name:      product.Name,
			Price:     product.Price,
			Discount:  discount,
		}
		if err := s.productOrderRepo.Create(po); err != nil {
			return map[string]interface{}{"error": "Gagal menambah produk ke order"}
		}
		lastID = po.ID.String()
	}
	return map[string]interface{}{"product_order_id": lastID}
}
func (s *outboundRegulerService) DeleteProduct(id string) interface{} {
	err := s.productOrderRepo.Delete(id)
	if err != nil {
		return map[string]interface{}{"error": "Gagal menghapus produk"}
	}
	return map[string]interface{}{"success": true}
}
func (s *outboundRegulerService) AddDiscount(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	dtype, _ := req["type"].(string)
	name, _ := req["name"].(string)
	isNominal, _ := req["is_nominal"].(bool)
	valueNominal, _ := req["value_nominal"].(float64)
	valuePercent, _ := req["value_percentage"].(int)
	do := &models.DiscountOrder{
		OrderID:         orderID,
		Type:            dtype,
		Name:            name,
		IsNominal:       isNominal,
		ValueNominal:    valueNominal,
		ValuePercentage: valuePercent,
	}
	if err := s.discountOrderRepo.Create(do); err != nil {
		return map[string]interface{}{"error": "Gagal menambah diskon"}
	}
	return map[string]interface{}{"discount_order_id": do.ID.String()}
}
func (s *outboundRegulerService) UpdateTax(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	isTax, _ := req["is_tax"].(bool)
	tax, _ := req["tax"].(float64)
	taxValue, _ := req["tax_value"].(float64)
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	order.Status = "progress"
	order.TotalPPN = tax
	order.GrandTotalAfter = taxValue
	order.IsTax = isTax
	if err := s.orderRepo.Update(order); err != nil {
		return map[string]interface{}{"error": "Gagal update tax"}
	}
	return map[string]interface{}{"success": true}
}
func (s *outboundRegulerService) UpdateBox(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	totalBox, _ := req["total_box"].(int)
	priceBox, _ := req["price_box"].(float64)
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	order.QuantityCartonBox = totalBox
	order.CartonBoxPrice = priceBox
	if err := s.orderRepo.Update(order); err != nil {
		return map[string]interface{}{"error": "Gagal update box"}
	}
	return map[string]interface{}{"success": true}
}
func (s *outboundRegulerService) CompleteOrder(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	// Hitung total harga produk
	products, _ := s.productOrderRepo.ListByOrderID(orderID)
	totalProduk := 0.0
	for _, p := range products {
		totalProduk += (p.Price - p.Discount)
	}
	// Hitung total box
	totalBox := float64(order.QuantityCartonBox) * order.CartonBoxPrice
	// Hitung diskon
	discounts, _ := s.discountOrderRepo.ListByOrderID(orderID)
	totalDiskon := 0.0
	for _, d := range discounts {
		if d.IsNominal {
			totalDiskon += d.ValueNominal
		} else {
			totalDiskon += float64(d.ValuePercentage) * totalProduk / 100.0
		}
	}
	// Hitung tax
	totalTax := order.TotalPPN
	// Grand total
	grandTotal := totalProduk + totalBox + totalTax - totalDiskon
	order.GrandTotalAfter = grandTotal
	order.Status = "done"
	if err := s.orderRepo.Update(order); err != nil {
		return map[string]interface{}{"error": "Gagal update order"}
	}
	return map[string]interface{}{"grand_total": grandTotal}
}
func (s *outboundRegulerService) GetOrderDetail(orderID string) interface{} {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	products, _ := s.productOrderRepo.ListByOrderID(orderID)
	discounts, _ := s.discountOrderRepo.ListByOrderID(orderID)
	return map[string]interface{}{
		"order":     order,
		"products":  products,
		"discounts": discounts,
	}
}
