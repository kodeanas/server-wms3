package services

import (
	"fmt"
	"time"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
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
	ListOrders() interface{}
	DeleteAllDiscountsByOrderID(orderID string) interface{}
}

type outboundRegulerService struct {
	buyerRepo         repositories.BuyerRepository
	classRepo         repositories.ClassRepository
	orderRepo         repositories.OrderRepository
	productOrderRepo  repositories.ProductOrderRepository
	discountOrderRepo repositories.DiscountOrderRepository
	categoryRepo      repositories.CategoryRepository
	productMasterRepo repositories.ProductMasterRepository
	taxRepo           *repositories.TaxRepository
}

func NewOutboundRegulerService(
	buyerRepo repositories.BuyerRepository,
	classRepo repositories.ClassRepository,
	orderRepo repositories.OrderRepository,
	productOrderRepo repositories.ProductOrderRepository,
	discountOrderRepo repositories.DiscountOrderRepository,
	categoryRepo repositories.CategoryRepository,
	productMasterRepo repositories.ProductMasterRepository,
	taxRepo *repositories.TaxRepository,
) OutboundRegulerService {
	return &outboundRegulerService{
		buyerRepo:         buyerRepo,
		classRepo:         classRepo,
		orderRepo:         orderRepo,
		productOrderRepo:  productOrderRepo,
		discountOrderRepo: discountOrderRepo,
		categoryRepo:      categoryRepo,
		productMasterRepo: productMasterRepo,
		taxRepo:           taxRepo,
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
		// Ambil buyer dan class_id
		buyer, err := s.buyerRepo.GetByID(buyerID)
		if err != nil || buyer == nil {
			return map[string]interface{}{"error": "Buyer tidak ditemukan"}
		}
		buyerUUID, err := uuid.Parse(buyerID)
		if err != nil {
			return map[string]interface{}{"error": "BuyerID tidak valid"}
		}
		classUUID, err := uuid.Parse(buyer.ClassID)
		if err != nil {
			return map[string]interface{}{"error": "ClassID buyer tidak valid"}
		}
		order = &models.Order{
			UserID:  nil,
			BuyerID: buyerUUID,
			ClassID: classUUID,
			Type:    "regular",
			Status:  "progress",
			Code:    fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
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

	// Validasi: produk dengan barcode sama sudah ada di order?
	products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
	for _, p := range products {
		if p.ProductID == product.ID.String() {
			return map[string]interface{}{"error": "Produk sudah pernah discan di order ini"}
		}
	}
	// Tambah product order
	po := &models.ProductOrder{
		OrderID:        order.ID.String(),
		ProductID:      product.ID.String(),
		Name:           product.Name,
		Price:          product.Price,
		PriceWarehouse: product.PriceWarehouse,
		Discount:       discount,
	}
	if err := s.productOrderRepo.Create(po); err != nil {
		return map[string]interface{}{"error": "Gagal menambah produk ke order"}
	}

	// Kalkulasi total harga produk dan grand total (pakai price_warehouse)
	products, _ = s.productOrderRepo.ListByOrderID(order.ID.String())
	totalProduk := 0.0
	for _, p := range products {
		totalProduk += (p.PriceWarehouse - p.Discount)
	}
	// Ambil diskon class buyer (persentase)
	classDiscount := 0.0
	if order.ClassID != uuid.Nil {
		class, err := s.classRepo.GetByID(order.ClassID.String())
		if err == nil {
			classDiscount = float64(class.Disc)
		}
	}
	diskonClass := totalProduk * classDiscount / 100.0
	totalSetelahDiskonClass := totalProduk - diskonClass
	// Hitung diskon (voucher, dsb)
	discounts, _ := s.discountOrderRepo.ListByOrderID(order.ID.String())
	totalDiskon := 0.0
	for _, d := range discounts {
		if d.IsNominal {
			totalDiskon += d.Value
		} else {
			totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
		}
	}
	// Hitung tax value
	taxValue := 0.0
	if order.IsTax {
		taxValue = order.Tax * totalSetelahDiskonClass / 100.0
	}
	// Hitung total box
	totalBoxValue := float64(order.TotalBox) * order.PriceBox
	// Grand total
	order.TotalPrice = totalProduk
	order.GrandTotal = totalSetelahDiskonClass - totalDiskon - taxValue + totalBoxValue
	_ = s.orderRepo.Update(order)

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

	// Cek apakah id benar-benar ada di tabel product_orders
	po, err := s.productOrderRepo.GetByID(id)
	if err != nil || po == nil {
		return map[string]interface{}{"error": "Produk order tidak ditemukan"}
	}
	err = s.productOrderRepo.Delete(id)
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
	value, _ := req["value"].(float64)
	usage, _ := req["usage"].(int)
	// Jika type voucher, usage = value
	if dtype == "voucher" {
		usage = int(value)
	}
	do := &models.DiscountOrder{
		OrderID:   orderID,
		Type:      dtype,
		Name:      name,
		IsNominal: isNominal,
		Value:     value,
		Usage:     usage,
	}
	if err := s.discountOrderRepo.Create(do); err != nil {
		return map[string]interface{}{"error": "Gagal menambah diskon"}
	}
	// Setelah menambah diskon, update grand total order
	order, err := s.orderRepo.GetByID(orderID)
	if err == nil && order != nil {
		// Hitung ulang grand total
		products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
		totalProduk := 0.0
		for _, p := range products {
			totalProduk += (p.PriceWarehouse - p.Discount)
		}
		classDiscount := 0.0
		if order.ClassID != uuid.Nil {
			class, err := s.classRepo.GetByID(order.ClassID.String())
			if err == nil {
				classDiscount = float64(class.Disc)
			}
		}
		diskonClass := totalProduk * classDiscount / 100.0
		totalSetelahDiskonClass := totalProduk - diskonClass
		discounts, _ := s.discountOrderRepo.ListByOrderID(order.ID.String())
		totalDiskon := 0.0
		for _, d := range discounts {
			if d.IsNominal {
				totalDiskon += d.Value
			} else {
				totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
			}
		}
		taxValue := 0.0
		if order.IsTax {
			taxValue = order.Tax * totalSetelahDiskonClass / 100.0
		}
		totalBoxValue := float64(order.TotalBox) * order.PriceBox
		order.GrandTotal = totalSetelahDiskonClass - totalDiskon - taxValue + totalBoxValue
		_ = s.orderRepo.Update(order)
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
	taxID, _ := req["tax_id"].(string)
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	order.Status = "progress"
	order.IsTax = isTax
	var grandTotal float64
	if isTax && taxID != "" {
		tax, err := s.taxRepo.FindByID(taxID)
		if err != nil || tax == nil {
			return map[string]interface{}{"error": "Tax tidak ditemukan"}
		}
		order.Tax = float64(tax.Tax)
		// Hitung totalProduk setelah diskon class
		products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
		totalProduk := 0.0
		for _, p := range products {
			totalProduk += (p.PriceWarehouse - p.Discount)
		}
		classDiscount := 0.0
		if order.ClassID != uuid.Nil {
			class, err := s.classRepo.GetByID(order.ClassID.String())
			if err == nil {
				classDiscount = float64(class.Disc)
			}
		}
		diskonClass := totalProduk * classDiscount / 100.0
		totalSetelahDiskonClass := totalProduk - diskonClass
		// Hitung diskon
		discounts, _ := s.discountOrderRepo.ListByOrderID(order.ID.String())
		totalDiskon := 0.0
		for _, d := range discounts {
			if d.IsNominal {
				totalDiskon += d.Value
			} else {
				totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
			}
		}
		// Hitung tax value
		order.TaxValue = order.Tax * totalSetelahDiskonClass / 100.0
		// Hitung total box
		totalBoxValue := float64(order.TotalBox) * order.PriceBox
		// Grand total: dikurangi tax, dikurangi voucher/diskon, ditambah box
		grandTotal = totalSetelahDiskonClass - totalDiskon - order.TaxValue + totalBoxValue
		order.GrandTotal = grandTotal
	} else {
		order.Tax = 0
		order.TaxValue = 0
		// Hitung totalProduk setelah diskon class
		products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
		totalProduk := 0.0
		for _, p := range products {
			totalProduk += (p.PriceWarehouse - p.Discount)
		}
		classDiscount := 0.0
		if order.ClassID != uuid.Nil {
			class, err := s.classRepo.GetByID(order.ClassID.String())
			if err == nil {
				classDiscount = float64(class.Disc)
			}
		}
		diskonClass := totalProduk * classDiscount / 100.0
		totalSetelahDiskonClass := totalProduk - diskonClass
		// Hitung diskon
		discounts, _ := s.discountOrderRepo.ListByOrderID(order.ID.String())
		totalDiskon := 0.0
		for _, d := range discounts {
			if d.IsNominal {
				totalDiskon += d.Value
			} else {
				totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
			}
		}
		// Hitung total box
		totalBoxValue := float64(order.TotalBox) * order.PriceBox
		// Grand total tanpa tax
		grandTotal = totalSetelahDiskonClass - totalDiskon + totalBoxValue
		order.GrandTotal = grandTotal
	}
	if err := s.orderRepo.Update(order); err != nil {
		return map[string]interface{}{"error": "Gagal update tax"}
	}
	return map[string]interface{}{"success": true, "tax_value": order.TaxValue, "grand_total": order.GrandTotal}
}
func (s *outboundRegulerService) UpdateBox(ctx interface{}) interface{} {
	req, ok := ctx.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "invalid request"}
	}
	orderID, _ := req["order_id"].(string)
	// Robust parsing for total_box
	var totalBox int
	if v, ok := req["total_box"].(int); ok {
		totalBox = v
	} else if v, ok := req["total_box"].(float64); ok {
		totalBox = int(v)
	} else {
		totalBox = 0
	}
	// Robust parsing for price_box
	var priceBox float64
	if v, ok := req["price_box"].(float64); ok {
		priceBox = v
	} else if v, ok := req["price_box"].(int); ok {
		priceBox = float64(v)
	} else {
		priceBox = 0
	}
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	order.TotalBox = totalBox
	order.PriceBox = priceBox

	// --- Hitung ulang grandTotal ---
	products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
	totalProduk := 0.0
	for _, p := range products {
		totalProduk += (p.PriceWarehouse - p.Discount)
	}
	classDiscount := 0.0
	if order.ClassID != uuid.Nil {
		class, err := s.classRepo.GetByID(order.ClassID.String())
		if err == nil {
			classDiscount = float64(class.Disc)
		}
	}
	diskonClass := totalProduk * classDiscount / 100.0
	totalSetelahDiskonClass := totalProduk - diskonClass
	// Hitung diskon
	discounts, _ := s.discountOrderRepo.ListByOrderID(order.ID.String())
	totalDiskon := 0.0
	for _, d := range discounts {
		if d.IsNominal {
			totalDiskon += d.Value
		} else {
			totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
		}
	}
	// Hitung tax value
	taxValue := 0.0
	if order.IsTax {
		taxValue = order.Tax * totalSetelahDiskonClass / 100.0
	}
	// Hitung total box
	totalBoxValue := float64(order.TotalBox) * order.PriceBox
	// Grand total
	order.GrandTotal = totalSetelahDiskonClass - totalDiskon - taxValue + totalBoxValue

	// DEBUG LOG
	fmt.Println("[DEBUG][UpdateBox] Before update:", "order_id=", order.ID, "total_box=", order.TotalBox, "price_box=", order.PriceBox, "grand_total=", order.GrandTotal)

	err = s.orderRepo.Update(order)
	fmt.Println("[DEBUG][UpdateBox] After update error:", err)

	if err != nil {
		return map[string]interface{}{"error": "Gagal update box"}
	}
	return map[string]interface{}{"success": true, "grand_total": order.GrandTotal, "total_box": order.TotalBox, "price_box": order.PriceBox}
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
	// Hitung total harga produk (pakai price_warehouse)
	products, _ := s.productOrderRepo.ListByOrderID(orderID)
	totalProduk := 0.0
	for _, p := range products {
		totalProduk += (p.PriceWarehouse - p.Discount)
	}
	// Ambil diskon class buyer (persentase)
	classDiscount := 0.0
	if order.ClassID != uuid.Nil {
		class, err := s.classRepo.GetByID(order.ClassID.String())
		if err == nil {
			classDiscount = float64(class.Disc)
		}
	}
	diskonClass := totalProduk * classDiscount / 100.0
	totalSetelahDiskonClass := totalProduk - diskonClass
	// Hitung diskon (voucher, dsb)
	discounts, _ := s.discountOrderRepo.ListByOrderID(orderID)
	totalDiskon := 0.0
	for _, d := range discounts {
		if d.IsNominal {
			totalDiskon += d.Value
		} else {
			totalDiskon += float64(d.Usage) * totalSetelahDiskonClass / 100.0
		}
	}
	// Hitung tax value
	taxValue := 0.0
	if order.IsTax {
		taxValue = order.Tax * totalSetelahDiskonClass / 100.0
	}
	// Hitung total box
	totalBoxValue := float64(order.TotalBox) * order.PriceBox
	// Grand total: dikurangi tax, dikurangi voucher/diskon, ditambah box
	grandTotal := totalSetelahDiskonClass - totalDiskon - taxValue + totalBoxValue
	order.GrandTotal = grandTotal
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
func (s *outboundRegulerService) ListOrders() interface{} {
	orders := []models.Order{}
	err := s.orderRepo.ListAll(&orders)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return orders
}

// Delete all discounts/vouchers for an order and recalculate grand total
func (s *outboundRegulerService) DeleteAllDiscountsByOrderID(orderID string) interface{} {
	// Get all discounts for the order
	discounts, err := s.discountOrderRepo.ListByOrderID(orderID)
	if err != nil {
		return map[string]interface{}{"error": "Gagal mengambil daftar voucher/discount"}
	}
	for _, d := range discounts {
		_ = s.discountOrderRepo.Delete(d.ID.String())
	}
	// Recalculate grand total
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil || order == nil {
		return map[string]interface{}{"error": "Order tidak ditemukan"}
	}
	products, _ := s.productOrderRepo.ListByOrderID(order.ID.String())
	totalProduk := 0.0
	for _, p := range products {
		totalProduk += (p.PriceWarehouse - p.Discount)
	}
	classDiscount := 0.0
	if order.ClassID != uuid.Nil {
		class, err := s.classRepo.GetByID(order.ClassID.String())
		if err == nil {
			classDiscount = float64(class.Disc)
		}
	}
	diskonClass := totalProduk * classDiscount / 100.0
	totalSetelahDiskonClass := totalProduk - diskonClass
	// No discounts left
	totalDiskon := 0.0
	taxValue := 0.0
	if order.IsTax {
		taxValue = order.Tax * totalSetelahDiskonClass / 100.0
	}
	totalBoxValue := float64(order.TotalBox) * order.PriceBox
	order.GrandTotal = totalSetelahDiskonClass - totalDiskon - taxValue + totalBoxValue
	_ = s.orderRepo.Update(order)
	return map[string]interface{}{"success": true, "grand_total": order.GrandTotal}
}
