# Inbound API Reference & Implementation Guide

## 📚 Quick Reference Table

### 1. INBOUND MANUAL Endpoints

| Method | Endpoint                      | Purpose                               | Auth |
| ------ | ----------------------------- | ------------------------------------- | ---- |
| GET    | `/api/inbound/list-masters`   | List semua ProductMaster              | -    |
| GET    | `/api/inbound/list-pendings`  | List semua ProductPending             | -    |
| GET    | `/api/inbound/manual-pending` | List pending items (manual type only) | -    |
| POST   | `/api/inbound/manual`         | Create manual inbound item            | -    |

### 2. INBOUND BULK Endpoints

| Method | Endpoint                        | Purpose                    | Auth |
| ------ | ------------------------------- | -------------------------- | ---- |
| POST   | `/api/inbound/bulk-upload`      | Upload & process bulk file | -    |
| GET    | `/api/inbound/bulk-summary-all` | Get bulk statistics        | -    |

### 3. INBOUND BAST Endpoints

| Method | Endpoint                                                  | Purpose                       | Auth |
| ------ | --------------------------------------------------------- | ----------------------------- | ---- |
| POST   | `/api/inbound/bast-upload`                                | Upload BAST file              | -    |
| GET    | `/api/inbound/bast-summary`                               | Get BAST summary (date range) | -    |
| GET    | `/api/inbound/bast-summary-all`                           | Get all BAST statistics       | -    |
| GET    | `/api/product-documents/bast`                             | List BAST documents           | -    |
| GET    | `/api/inbound/bast-scanner/document/:document_id`         | Get document for scanner      | -    |
| GET    | `/api/inbound/bast-scanner/:document_id/product/:barcode` | Get product by barcode        | -    |
| POST   | `/api/inbound/bast-scanner/:document_id/scan/:barcode`    | Scan & move to master         | -    |
| POST   | `/api/inbound/bast-scanner/:document_id/finish`           | Finish BAST document          | -    |

### 4. INBOUND SKU Endpoints

| Method | Endpoint                                  | Purpose                              | Auth |
| ------ | ----------------------------------------- | ------------------------------------ | ---- |
| GET    | `/api/inbound-sku/summary-all`            | Get SKU statistics                   | -    |
| POST   | `/api/inbound-sku/upload`                 | Upload SKU file                      | -    |
| POST   | `/api/inbound-sku/crosscheck/:pending_id` | Crosscheck item (split good/damaged) | -    |
| POST   | `/api/inbound-sku/finish/:document_id`    | Finish SKU document                  | -    |
| GET    | `/api/inbound-sku/document/:document_id`  | Get SKU document detail              | -    |

---

## 🔧 Implementation Guide

### Inbound Manual - Complete Example

**Controller Implementation:**

```go
// File: controller/inbound_manual_controller.go

func InboundManualHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req models.InboundRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            // Handle validation error
            utils.SendValidationError(c, verrs)
            return
        }

        inboundService := services.NewInboundService(nil)
        _, master, err := inboundService.InboundManual(req, db)
        if err != nil {
            utils.SendError(c, 500, err.Error())
            return
        }

        utils.SendSuccess(c, gin.H{
            "message": "Berhasil migrasi",
            "barcode": master.BarcodeWarehouse,
            "price": master.Price,
            "price_warehouse": master.PriceWarehouse,
        }, "OK", nil, http.StatusOK)
    }
}
```

**Request Body:**

```json
{
  "name": "Laptop ASUS VivoBook 15",
  "item": 1,
  "price": 8000000,
  "category_id": "550e8400-e29b-41d4-a716-446655440000",
  "sticker_id": null,
  "status": "good",
  "note": null
}
```

**Response:**

```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "message": "Berhasil migrasi",
    "barcode": "BC-1705000000000-12345",
    "price": 8000000,
    "price_warehouse": 7200000,
    "name": "Laptop ASUS VivoBook 15",
    "category_name": "Electronics"
  }
}
```

---

### Inbound Bulk - Complete Example

**Controller Implementation:**

```go
// File: controller/inbound_bulk_controller.go

func InboundBulkUploadHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        supplier := c.PostForm("supplier")
        typeProduct := c.PostForm("type_product")
        fileType := c.PostForm("type")

        file, header, err := c.Request.FormFile("file")
        if err != nil {
            utils.SendError(c, 400, "File tidak ditemukan")
            return
        }
        defer file.Close()

        req := models.BulkInboundRequest{
            FileName: header.Filename,
            Supplier: supplier,
            TypeProduct: typeProduct,
            Type: fileType,
            Mapping: models.BulkInboundMapping{
                BarcodeHeader: c.PostForm("barcode_header"),
                NameHeader: c.PostForm("name_header"),
                QtyHeader: c.PostForm("qty_header"),
                PriceHeader: c.PostForm("price_header"),
            },
        }

        inserted, skipped, skipDetails := inboundBulkService.InboundBulkProcess(req, db)

        utils.SendSuccess(c, gin.H{
            "inserted": inserted,
            "skipped": skipped,
            "skip_details": skipDetails,
        }, "Bulk inbound selesai", nil, http.StatusOK)
    }
}
```

**Request (Form Data):**

```
POST /api/inbound/bulk-upload
Content-Type: multipart/form-data

supplier: "PT Elektronik Indonesia"
type_product: "reguler"
type: "xlsx"
barcode_header: "barcode"
name_header: "product_name"
qty_header: "quantity"
price_header: "harga"
file: [file.xlsx]
```

**Response:**

```json
{
  "code": 200,
  "message": "Bulk inbound selesai",
  "data": {
    "inserted": 95,
    "skipped": 5,
    "skip_details": [
      "Row 10 skipped: kategori tidak ditemukan di DB: 'Elektronik Premium'",
      "Row 23 skipped: qty/price tidak valid"
    ],
    "filename": "produk_bulk_001.xlsx"
  }
}
```

---

### Inbound BAST - Complete Example

**1. Upload BAST:**

```bash
curl -X POST http://localhost:8080/api/inbound/bast-upload \
  -F "supplier=PT Supplier A" \
  -F "header_barcode=barcode" \
  -F "header_name=product_name" \
  -F "header_item=qty" \
  -F "header_price=price" \
  -F "type=xlsx" \
  -F "file=@bast_file.xlsx"
```

**Response:**

```json
{
  "code": 200,
  "message": "Inbound BAST selesai",
  "data": {
    "inserted": 100,
    "skipped": 0,
    "skip_details": [],
    "filename": "bast_file.xlsx"
  }
}
```

**2. Get Document for Scanner:**

```bash
curl http://localhost:8080/api/inbound/bast-scanner/document/550e8400-e29b-41d4-a716-446655440000
```

**Response:**

```json
{
  "code": 200,
  "data": {
    "document": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "code": "BAST-1705000000000",
      "file_name": "bast_file.xlsx",
      "status": "progress",
      "scanned_count": 20,
      "unscanned_count": 80
    },
    "products": [
      {
        "id": "product-uuid",
        "barcode": "BC-001",
        "name": "Laptop ASUS",
        "item": 10,
        "price": 8000000,
        "status": "discrepancy",
        "date_scanned": null
      }
    ]
  }
}
```

**3. Scan Single Product:**

```bash
curl -X POST http://localhost:8080/api/inbound/bast-scanner/550e8400-e29b-41d4-a716-446655440000/scan/BC-001 \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "cat-001",
    "status": "good",
    "note": "Kondisi baik"
  }'
```

**Response:**

```json
{
  "code": 200,
  "message": "Scan berhasil",
  "data": {
    "message": "Scan berhasil",
    "barcode_warehouse": "BC-001-WH",
    "category_name": "Electronics",
    "price": 8000000,
    "price_warehouse": 7200000,
    "status": "good"
  }
}
```

**4. Finish Document:**

```bash
curl -X POST http://localhost:8080/api/inbound/bast-scanner/550e8400-e29b-41d4-a716-446655440000/finish
```

---

### Inbound SKU - Complete Example

**1. Upload SKU:**

```bash
curl -X POST http://localhost:8080/api/inbound-sku/upload \
  -F "supplier=PT SKU Supplier" \
  -F "file=@sku_file.xlsx"
```

**Response:**

```json
{
  "code": 200,
  "message": "Inbound BAST selesai",
  "data": {
    "inserted": 50,
    "skipped": 0,
    "skip_details": [],
    "filename": "sku_file.xlsx"
  }
}
```

**2. Get Document:**

```bash
curl http://localhost:8080/api/inbound-sku/document/doc-uuid
```

**Response:**

```json
{
  "code": 200,
  "data": {
    "document": {
      "id": "doc-uuid",
      "file_name": "sku_file.xlsx",
      "type": "sku",
      "status": "pending",
      "supplier": "PT SKU Supplier"
    },
    "product_pending": [
      {
        "id": "pending-uuid",
        "barcode": "SKU-001",
        "name": "Product A",
        "item": 10,
        "price": 50000,
        "item_good": null,
        "item_damaged": null,
        "status": "good"
      }
    ]
  }
}
```

**3. Crosscheck Item:**

```bash
curl -X POST http://localhost:8080/api/inbound-sku/crosscheck/pending-uuid \
  -H "Content-Type: application/json" \
  -d '{
    "item_good": 8,
    "item_damaged": 2
  }'
```

**Response:**

```json
{
  "code": 200,
  "message": "Crosscheck updated"
}
```

**4. Finish SKU:**

```bash
curl -X POST http://localhost:8080/api/inbound-sku/finish/doc-uuid
```

**Response:**

```json
{
  "code": 200,
  "message": "Inbound SKU finished"
}
```

---

## 📋 Data Validation Rules

### Manual Inbound Validation

```go
// Required fields
- name: string, non-empty
- item: int > 0
- price: float64 > 0
- status: string in ["good", "abnormal", "damaged", "non"]

// Conditional
- if status != "good": note recommended
- if price >= 100000: category_id recommended (for pricing strategy)
- if price < 100000: sticker_id optional (auto-match if not provided)
```

### Bulk Inbound Validation

```go
// File Validation
- Supported types: .csv, .xlsx, .xls
- Min 1 row, Max 10000 rows
- Required columns: barcode, name, qty, price

// Data Validation per row
- barcode: non-empty
- name: non-empty
- qty: integer > 0
- price: numeric > 0
- For reguler: price >= 100000, kategori harus exist
- For sticker: price < 100000, sticker harus exist dengan range match
```

### BAST Validation

```go
// File same as Bulk
// Header mapping required
- header_barcode: column name in file
- header_name: column name in file
- header_item: column name in file
- header_price: column name in file

// Scan validation
- barcode: must exist in ProductPending (not yet scanned)
- category_id: must exist in Categories (if reguler product)
- status: must be one of ["good", "abnormal", "damaged", "non"]
```

### SKU Validation

```go
// File Validation (Excel)
// Expected structure:
// Col 0: barcode
// Col 1: name
// Col 2: ??? (ignored)
// Col 3: price
// Col 4: item (initial quantity)

// Crosscheck Validation
- item_good: integer >= 0
- item_damaged: integer >= 0
- item_good + item_damaged <= pending.item
```

---

## 🗂️ File Structure Reference

### Key Files Location

```
wms-v3/
├── controller/
│   ├── inbound_manual_controller.go      ← Manual handlers
│   ├── inbound_bulk_controller.go        ← Bulk handlers
│   ├── inbound_bast_controller.go        ← BAST handlers
│   └── inbound_sku_controller.go         ← SKU handlers
│
├── services/
│   ├── inbound_service.go                ← Manual service logic
│   ├── inbound_bulk_service.go           ← Bulk service logic
│   ├── inbound_bast_service.go           ← BAST service logic
│   └── inbound_sku_service.go            ← SKU service logic
│
├── models/
│   ├── product_document.go               ← Document model
│   ├── product_pending.go                ← Pending model
│   ├── product_master.go                 ← Master model
│   ├── product_repair.go                 ← Repair model (SKU)
│   ├── category.go                       ← Category model
│   └── sticker.go                        ← Sticker model
│
├── repositories/
│   ├── product_document_repository.go    ← Document queries
│   ├── product_pending_repository.go     ← Pending queries
│   ├── product_master_repository.go      ← Master queries
│   └── ...
│
├── routes/
│   └── routes.go                         ← API route definitions
│
└── utils/
    └── barcode.go                        ← File parsing utilities
```

---

## 🔍 Common Errors & Troubleshooting

### Error: "Kategori tidak ditemukan"

- **Cause**: Category name in file doesn't match database
- **Solution**:
  - Check category name spelling in database
  - Use exact case from database
  - Ensure category exists in Categories table

### Error: "Sticker range tidak match"

- **Cause**: Product price outside sticker min_price/max_price range
- **Solution**:
  - Check sticker min/max range in database
  - Add new sticker with correct price range

### Error: "File tidak valid"

- **Cause**: File format or structure incorrect
- **Solution**:
  - Verify file format (.csv/.xlsx/.xls)
  - Check header names match mapping
  - Ensure data starts from row 2 (row 1 = header)

### Error: "Barcode sudah exist"

- **Cause**: Duplicate barcode in file or database
- **Solution**:
  - Remove duplicate entries in file
  - Use unique barcode per product

### Error: "Column kurang lengkap"

- **Cause**: Missing required columns in file
- **Solution**:
  - Verify all required columns present
  - Check column order matches mapping

---

## 📊 Database Query Examples

### Get Inbound Statistics

```sql
-- Total inbound by type
SELECT type, COUNT(*) as count, SUM(file_item) as total_items
FROM product_documents
WHERE type IN ('manual', 'bulk', 'bast', 'sku')
GROUP BY type;

-- Document status breakdown
SELECT type, status, COUNT(*) as count
FROM product_documents
GROUP BY type, status;

-- Pending by status
SELECT status, COUNT(*) as count
FROM product_pendings
GROUP BY status;

-- Product location distribution
SELECT location, COUNT(*) as count
FROM product_masters
GROUP BY location;
```

### Find Specific Items

```sql
-- Find product by barcode across all inbound
SELECT
  pd.id as pending_id,
  pm.id as master_id,
  pp.barcode,
  pp.name,
  doc.type,
  doc.status,
  pp.status
FROM product_pendings pp
LEFT JOIN product_masters pm ON pp.id = pm.product_pending_id
LEFT JOIN product_documents doc ON pp.document_id = doc.id
WHERE pp.barcode = 'BC-001';

-- Find all items in specific document
SELECT
  pp.barcode,
  pp.name,
  pp.item,
  pp.price,
  pm.location,
  pp.status
FROM product_pendings pp
LEFT JOIN product_masters pm ON pp.product_pending_id = pm.id
WHERE pp.document_id = 'doc-uuid'
ORDER BY pp.created_at;
```

---

## 🚀 Performance Tips

1. **Bulk Upload**:
   - Use XLSX instead of CSV for faster parsing
   - Limit file to 5000 rows for optimal performance
   - Pre-validate data before upload

2. **Scanner Process**:
   - Cache category and sticker lists in memory
   - Use batch updates if scanning multiple items
   - Keep document in "progress" status until all items scanned

3. **Database**:
   - Index on (document_id, barcode) for faster lookups
   - Archive old documents to separate table
   - Regular VACUUM ANALYZE for PostgreSQL

---

## 📚 Related Documentation

- See `INBOUND_PROCESS_DOCUMENTATION.md` for detailed process flows
- See `INBOUND_DIAGRAMS.md` for visual flow diagrams
- See `routes/routes.go` for complete route definitions
