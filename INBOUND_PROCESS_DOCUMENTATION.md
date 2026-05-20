# Dokumentasi Proses Inbound WMS-V3

Dokumen ini menjelaskan 4 jenis proses data masuk (inbound) di sistem WMS: Manual, Bulk, BAST, dan SKU.

---

## 📋 Daftar Isi

1. [Inbound Manual](#inbound-manual)
2. [Inbound Bulk](#inbound-bulk)
3. [Inbound BAST](#inbound-bast)
4. [Inbound SKU](#inbound-sku)

---

## 1. Inbound Manual

### 📌 Definisi

Inbound Manual adalah proses input data barang secara individual/manual melalui form tanpa memerlukan upload file.

### 🎯 Tujuan

- Input barang satu per satu (retail/dropship)
- Fleksibel dalam kategori dan sticker
- Cepat untuk kebutuhan urgent

### 📊 Alur Data

```
Input Form (Manual)
    ↓
Validasi Input
    ↓
Generate Barcode
    ↓
Tentukan Kategori/Sticker (berdasarkan harga)
    ↓
Hitung Price Warehouse (berdasarkan diskon kategori/sticker)
    ↓
CREATE → ProductPending
    ↓
CREATE → ProductMaster
    ↓
Set Location (staging_reguler atau staging_sticker)
```

### 🔌 API Endpoint

```
POST /api/inbound/manual
```

### 📥 Input Request

```json
{
  "name": "string (required)", // Nama produk
  "item": "number (required)", // Qty item
  "price": "number (required)", // Harga original
  "category_id": "string (optional)", // ID kategori (untuk reguler, harga >= 100rb)
  "sticker_id": "string (optional)", // ID sticker (untuk sticker, harga < 100rb)
  "status": "string (required)", // good|abnormal|damaged|non
  "note": "string (optional)" // Catatan jika status != good
}
```

### 📤 Output Response

```json
{
  "message": "Berhasil migrasi",
  "barcode": "BC-xxxxxx-xxxxx", // Barcode warehouse (generated)
  "price": 50000, // Harga original
  "price_warehouse": 45000, // Harga setelah diskon (jika ada)
  "name": "Produk A",
  "category_name": "Electronics"
}
```

### 🧮 Logic Penentuan Kategori & Sticker

| Kondisi           | Type       | Location        | Aksi                                             |
| ----------------- | ---------- | --------------- | ------------------------------------------------ |
| Harga ≥ 100.000   | categories | staging_reguler | Gunakan kategori (dapat diskon dari kategori)    |
| Harga < 100.000   | sticker    | staging_sticker | Gunakan sticker (dapat fixed price dari sticker) |
| Status = damaged  | -          | damaged         | Direct ke lokasi damaged                         |
| Status = abnormal | -          | abnormal        | Direct ke lokasi abnormal                        |
| Status = non      | -          | (default)       | Direct ke lokasi default                         |

### 💾 Database Operations

**Table: product_pendings**

- `id`: UUID (auto)
- `document_id`: Reference ke INBOUND_MANUAL document
- `barcode`: Generated
- `name`: Input
- `item`: Input
- `price`: Input
- `status`: Input (good/abnormal/damaged/non)
- `note`: Input jika status != good

**Table: product_masters**

- `id`: UUID (auto)
- `document_id`: INBOUND_MANUAL
- `barcode`: Same as pending
- `barcode_warehouse`: Generated
- `name`: Input
- `name_warehouse`: "Manual"
- `item`: Input
- `price`: Input
- `price_warehouse`: Calculated (price × (1 - discount/100))
- `category_id`: From input atau dari sticker range
- `sticker_id`: From input atau dari price range
- `product_pending_id`: Reference ke ProductPending
- `location`: Determined by status & category logic

### 🛠️ Implementation Files

- **Controller**: `controller/inbound_manual_controller.go` → `InboundManualHandler()`
- **Service**: `services/inbound_service.go` → `InboundManual()`
- **Handler**: Fungsi handler database

### 📝 Contoh Use Case

**Scenario 1: Input produk elektronik reguler**

```
Input: name="Smartphone A", item=1, price=150000, category_id="cat-123", status="good"
Output:
  - ProductPending: barcode=BC-xxxx, status=good
  - ProductMaster: price_warehouse=150000 (jika kategori tidak ada diskon) atau lebih rendah (jika kategori ada diskon)
  - Location: staging_reguler
```

**Scenario 2: Input produk sticker (harga murah)**

```
Input: name="Tissue", item=5, price=25000, sticker_id="sticker-123", status="good"
Output:
  - ProductPending: barcode=BC-xxxx, status=good
  - ProductMaster: price_warehouse=15000 (fixed price dari sticker), sticker_id=sticker-123
  - Location: staging_sticker
```

**Scenario 3: Input produk rusak**

```
Input: name="Printer", item=1, price=800000, category_id="cat-456", status="damaged", note="Rusak saat kirim"
Output:
  - ProductPending: barcode=BC-xxxx, status=damaged, note="Rusak saat kirim"
  - ProductMaster: location=damaged
```

---

## 2. Inbound Bulk

### 📌 Definisi

Inbound Bulk adalah proses input data barang dalam jumlah besar melalui upload file (CSV/XLSX/XLS) dengan mapping header.

### 🎯 Tujuan

- Import data produk dalam jumlah besar sekaligus
- Otomatis matching dengan kategori/sticker di database
- Menghitung total item dan harga dari file

### 📊 Alur Data

```
Upload File (CSV/XLSX/XLS)
    ↓
Validasi Format File
    ↓
Parse File & Extract Rows
    ↓
Hitung Total Item & Price dari File
    ↓
CREATE → ProductDocument (type="bulk")
    ↓
Loop setiap row:
  ├─ Validate data completeness
  ├─ Match category/sticker (by name atau price range)
  ├─ Hitung price_warehouse
  ├─ CREATE → ProductPending
  └─ CREATE → ProductMaster
```

### 🔌 API Endpoint

```
POST /api/inbound/bulk-upload
Content-Type: multipart/form-data
```

### 📥 Input Request (Form Data)

```
file:               [File] - CSV/XLSX/XLS file
supplier:           string - Nama supplier
type_product:       string - "reguler" atau "sticker"
type:               string - "csv", "xlsx", atau "xls"
barcode_header:     string - Nama kolom barcode (default: "barcode")
name_header:        string - Nama kolom nama (default: "name")
qty_header:         string - Nama kolom qty (default: "qty")
price_header:       string - Nama kolom harga (default: "price")
category_header:    string - (optional, untuk type_product=reguler)
```

### 📤 Output Response

```json
{
  "inserted": 150, // Jumlah barang berhasil
  "skipped": 5, // Jumlah barang skip
  "skip_details": [
    // Detail alasan skip
    "Row skipped: qty/price tidak valid",
    "Row skipped: kategori tidak ditemukan"
  ],
  "filename": "produk_bulk_001.xlsx"
}
```

### 📝 Format File Example

**Type: reguler (harga >= 100.000)**

```
barcode    | name          | category      | qty | price
-----------+---------------+---------------+-----+--------
BC-001     | Laptop ASUS    | Electronics   | 10  | 8000000
BC-002     | Monitor LG     | Electronics   | 15  | 2500000
BC-003     | Keyboard Mech  | Accessories   | 25  | 500000
```

**Type: sticker (harga < 100.000)**

```
barcode    | name          | qty | price
-----------+---------------+-----+-------
ST-001     | Tissue Pack   | 100 | 25000
ST-002     | Soap Bottle   | 50  | 35000
ST-003     | Pen Set       | 200 | 15000
```

### 🧮 Validation Rules

**Untuk type_product = "reguler":**

- Harga minimal: 100.000
- Kategori harus ada di database (case-insensitive match)
- Jika kategori tidak match → SKIP

**Untuk type_product = "sticker":**

- Harga maksimal: < 100.000
- Sticker harus ada di database dengan range harga match
- Jika sticker tidak match → SKIP

### 💾 Database Operations

**Table: product_documents**

- `id`: UUID (auto)
- `code`: BULK-{timestamp}
- `file_name`: Dari input
- `file_item`: Sum dari qty di file
- `file_price`: Sum dari price × qty di file
- `type`: "bulk"
- `type_product`: "reguler" atau "sticker"
- `header_barcode`, `header_name`, `header_item`, `header_price`: Mapping headers
- `supplier`: Input
- `status`: "progress"

**Table: product_pendings** (per row)

- `document_id`: Reference ke bulk document
- `barcode`: Dari file
- `name`: Dari file
- `item`: qty dari file
- `price`: price dari file
- `status`: "good" (default untuk bulk)

**Table: product_masters** (per row)

- `document_id`: Reference ke bulk document
- `barcode`: Sama dengan pending
- `category_id`: Matched dari database
- `sticker_id`: Matched dari database
- `price_warehouse`: Calculated berdasarkan diskon kategori atau fixed price sticker
- `location`: "staging_reguler" atau "staging_sticker"

### 🛠️ Implementation Files

- **Controller**: `controller/inbound_bulk_controller.go` → `InboundBulkUploadHandler()`
- **Service**: `services/inbound_bulk_service.go` → `InboundBulkProcess()`
- **Utils**: `utils/barcode.go` → Parse file utilities

### 📝 Contoh Use Case

**Scenario: Upload 100 produk elektronik**

```
Input:
  - File: produk.xlsx (100 rows)
  - supplier: "PT Elektronik Indonesia"
  - type_product: "reguler"
  - mapping: barcode, name, category, qty, price

Output:
  - ProductDocument: created dengan file_item=200, file_price=50000000
  - ProductPending: 95 berhasil, 5 skipped (kategori tidak match)
  - ProductMaster: 95 berhasil dengan location=staging_reguler
```

---

## 3. Inbound BAST

### 📌 Definisi

BAST (Berita Acara Serah Terima) adalah proses inbound dengan upload file kemudian scanner per-item untuk validasi dan status quality control.

### 🎯 Tujuan

- Upload data barang dari penerima/supplier
- Per-item QC validation via scanner/barcode
- Tentukan status produk: good, abnormal, damaged, atau non
- Automatic assignment ke kategori

### 📊 Alur Data

```
Upload File BAST (CSV/XLSX/XLS)
    ↓
Validasi Format & Header Mapping
    ↓
CREATE → ProductDocument (type="bast", status="progress")
    ↓
Loop setiap row → CREATE ProductPending (status="discrepancy")
    ↓
[Optional] GET Document Info & Summary (scanned vs unscanned count)
    ↓
[Per Item] Get Pending by Barcode
    ↓
[Per Item] POST Scan & Move to ProductMaster
    ├─ Input: category_id, status (good/abnormal/damaged/non), note
    ├─ Validate kategori untuk reguler product
    ├─ Calculate price_warehouse dari kategori/sticker
    ├─ CREATE → ProductMaster
    ├─ UPDATE ProductPending → DateScanned
    └─ Response: Master data baru
    ↓
[Finish] Mark Document as Done
```

### 🔌 API Endpoints

#### 1. Upload BAST File

```
POST /api/inbound/bast-upload
Content-Type: multipart/form-data
```

#### 2. Get Document Info (untuk scanner)

```
GET /api/inbound/bast-scanner/document/:document_id
```

#### 3. Get Product by Barcode

```
GET /api/inbound/bast-scanner/:document_id/product/:barcode
```

#### 4. Scan & Move to Master

```
POST /api/inbound/bast-scanner/:document_id/scan/:barcode
```

#### 5. Finish Document

```
POST /api/inbound/bast-scanner/:document_id/finish
```

### 📥 Input Request

**Upload BAST:**

```
file:               [File] - CSV/XLSX/XLS file
supplier:           string - Nama supplier
header_barcode:     string - Nama kolom barcode
header_name:        string - Nama kolom nama
header_item:        string - Nama kolom quantity
header_price:       string - Nama kolom harga
type:               string - "csv", "xlsx", atau "xls"
```

**Scan & Move:**

```json
{
  "category_id": "string (required untuk reguler)", // Kategori produk
  "status": "string (required)", // good|abnormal|damaged|non
  "note": "string (optional)" // Catatan/alasan
}
```

### 📤 Output Response

**Upload:**

```json
{
  "inserted": 95,
  "skipped": 5,
  "skip_details": ["Row skipped: ..."],
  "filename": "bast_001.xlsx"
}
```

**Get Document:**

```json
{
  "document": {
    "id": "doc-uuid",
    "code": "BAST-xxxxxx",
    "file_name": "bast_001.xlsx",
    "status": "progress",
    "scanned_count": 20,
    "unscanned_count": 75
  },
  "products": [
    {
      "id": "pending-uuid",
      "barcode": "BC-001",
      "name": "Laptop",
      "item": 5,
      "price": 8000000,
      "status": "discrepancy",
      "date_scanned": null
    }
  ]
}
```

**Scan & Move:**

```json
{
  "message": "Scan berhasil",
  "barcode_warehouse": "BC-001-WH",
  "category_name": "Electronics",
  "price": 8000000,
  "price_warehouse": 7200000,
  "status": "good"
}
```

### 🧮 Key Logic

**File Upload:**

- Semua produk di-create sebagai ProductPending dengan status="discrepancy"
- Belum ada ProductMaster saat upload
- Document status="progress"

**Scanner Process:**

- Scanner harus membaca barcode dari ProductPending
- Input kategori & status dari user/QC staff
- Otomatis hitung price_warehouse dari kategori/sticker logic
- Create ProductMaster dengan data hasil scan
- Update ProductPending dengan DateScanned

**Status Mapping:**

```
good       → ProductMaster location = staging_reguler/staging_sticker
abnormal   → ProductMaster location = abnormal
damaged    → ProductMaster location = damaged
non        → ProductMaster location = non (atau non-sellable)
```

### 💾 Database Operations

**Table: product_documents**

- `type`: "bast"
- `status`: "progress" (saat upload) → "done" (setelah finish)
- `code`: BAST-{timestamp}
- `header_barcode`, `header_name`, `header_item`, `header_price`: Header mapping
- `file_item`: Total item dari file
- `file_price`: Total harga dari file

**Table: product_pendings**

- `document_id`: Reference ke BAST document
- `status`: "discrepancy" (awal) → akan tetap discrepancy sampai di-scan
- `date_scanned`: NULL sampai di-scan, kemudian terisi timestamp

**Table: product_masters** (created saat scan)

- `document_id`: Reference ke BAST document
- `status`: Dari input user (good/abnormal/damaged/non)
- `product_pending_id`: Reference ke ProductPending yang di-scan
- `category_id`: Input dari scanner
- `price_warehouse`: Calculated

### 🛠️ Implementation Files

- **Controller**: `controller/inbound_bast_controller.go`
- **Service**: `services/inbound_bast_service.go`
- **Handlers**: Multiple handlers untuk upload, get document, scan

### 📝 Contoh Use Case

**Scenario: QC produk elektronik via BAST**

1. **Upload file BAST:**

   ```
   File: bast_electronics_2024.xlsx
   - 100 rows produk dari supplier
   - Semua masuk ke ProductPending dengan status="discrepancy"
   ```

2. **Scanner digunakan untuk setiap produk:**

   ```
   Scan barcode BC-001
   ├─ Display: Laptop ASUS, 5pcs, harga 8000000
   ├─ Staff select: category=Electronics, status=good
   └─ Result: ProductMaster created, location=staging_reguler

   Scan barcode BC-002
   ├─ Display: Monitor LG, 3pcs
   ├─ Staff select: category=Electronics, status=damaged
   └─ Result: ProductMaster created, location=damaged
   ```

3. **Finish document:**
   ```
   Mark document status=done
   Semua unscan product akan di-handle (discard/retry)
   ```

---

## 4. Inbound SKU

### 📌 Definisi

Inbound SKU adalah proses inbound khusus untuk stok keeping unit dengan upload file dan crosscheck per-item untuk split good/damaged quantity.

### 🎯 Tujuan

- Upload data SKU dengan initial quantity
- Crosscheck per-item untuk quality validation
- Split quantity ke good dan damaged items
- Move ke ProductMaster sesuai split quantity

### 📊 Alur Data

```
Upload Excel File SKU
    ↓
Validasi Format File
    ↓
CREATE → ProductDocument (type="sku", status="pending")
    ↓
Loop setiap row → CREATE ProductPending (IsSKU=true, status="good")
    ↓
[Per Item] CrossCheck Pending
  ├─ Input: item_good, item_damaged
  ├─ Validate: item_good + item_damaged <= total item
  ├─ UPDATE ProductPending → ItemGood, ItemDamaged
    ↓
[Finish] Mark Document as Done
  ├─ Loop ProductPending:
  │   ├─ Jika ItemGood > 0 → CREATE ProductMaster (location="staging_sku")
  │   └─ Jika ItemDamaged > 0 → CREATE ProductRepair
  │
  └─ UPDATE Document → status="done"
```

### 🔌 API Endpoints

#### 1. Upload Excel

```
POST /api/inbound-sku/upload
Content-Type: multipart/form-data
```

#### 2. Get Document Info

```
GET /api/inbound-sku/document/:document_id
```

#### 3. Crosscheck Pending

```
POST /api/inbound-sku/crosscheck/:pending_id
```

#### 4. Finish Inbound SKU

```
POST /api/inbound-sku/finish/:document_id
```

### 📥 Input Request

**Upload:**

```
file:       [File] - Excel file
supplier:   string - Nama supplier
```

**Crosscheck:**

```json
{
  "item_good": 8, // Jumlah item baik
  "item_damaged": 2 // Jumlah item rusak
}
```

### 📤 Output Response

**Upload:**

```json
{
  "inserted": 50,
  "skipped": 2,
  "skip_details": ["Row 1: kolom kurang lengkap"],
  "filename": "sku_001.xlsx"
}
```

**Get Document:**

```json
{
  "document": {
    "id": "doc-uuid",
    "file_name": "sku_001.xlsx",
    "type": "sku",
    "status": "pending",
    "supplier": "Supplier A"
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
```

**Crosscheck:**

```json
{
  "message": "Crosscheck updated"
}
```

**Finish:**

```json
{
  "message": "Inbound SKU finished"
}
```

### 📝 Format File Example

**Excel File Structure:**

```
Column 0: barcode    | Column 1: name         | Column 2: ??? | Column 3: price | Column 4: item
-----------+-----------+--------+-------+------
SKU-001    | Product A | ???    | 50000 | 10
SKU-002    | Product B | ???    | 75000 | 15
SKU-003    | Product C | ???    | 30000 | 20
```

### 🧮 Key Logic

**Upload Process:**

- Parse Excel → Extract rows
- Kolom: [barcode(0), name(1), ???(2), price(3), item(4)]
- Create ProductPending: IsSKU=true, status="good", Item dari file
- ItemGood dan ItemDamaged masih NULL (akan di-isi saat crosscheck)

**Crosscheck Process:**

- Input item_good dan item_damaged dari staff
- Validate: item_good + item_damaged <= total item
- Update ProductPending dengan nilai tersebut

**Finish Process:**

- Loop semua ProductPending
- Jika item_good > 0 → Create ProductMaster:
  - ItemWarehouse = item_good
  - Location = "staging_sku"
  - IsSKU = true
- Jika item_damaged > 0 → Create ProductRepair:
  - ItemDamaged = item_damaged
  - ProductPendingID reference
- Update ProductDocument status = "done"

### 💾 Database Operations

**Table: product_documents**

- `type`: "sku"
- `status`: "pending" (awal) → "done" (setelah finish)
- `supplier`: Input

**Table: product_pendings**

- `is_sku`: true
- `item`: Initial quantity dari file
- `item_good`: NULL → updated saat crosscheck
- `item_damaged`: NULL → updated saat crosscheck

**Table: product_masters** (created saat finish)

- `is_sku`: true
- `item_warehouse`: Dari item_good
- `location`: "staging_sku"
- `product_pending_id`: Reference ke ProductPending

**Table: product_repairs** (created saat finish jika item_damaged > 0)

- `item_damaged`: Dari item_damaged
- `product_pending_id`: Reference ke ProductPending

### 🛠️ Implementation Files

- **Controller**: `controller/inbound_sku_controller.go`
- **Service**: `services/inbound_sku_service.go`
- **Models**: `models/product_repair.go`

### 📝 Contoh Use Case

**Scenario: Import 50 SKU dengan QC split good/damaged**

1. **Upload file:**

   ```
   File: sku_batch_001.xlsx
   - 50 rows SKU
   - Setiap row punya barcode, name, price, item
   - Create 50 ProductPending dengan IsSKU=true
   ```

2. **QC Staff melakukan Crosscheck per item:**

   ```
   SKU-001: 10 items received
   ├─ Input: item_good=8, item_damaged=2
   └─ Update ProductPending: ItemGood=8, ItemDamaged=2

   SKU-002: 15 items received
   ├─ Input: item_good=15, item_damaged=0
   └─ Update ProductPending: ItemGood=15, ItemDamaged=0

   SKU-003: 20 items received
   ├─ Input: item_good=19, item_damaged=1
   └─ Update ProductPending: ItemGood=19, ItemDamaged=1
   ```

3. **Finish Document:**

   ```
   System process:
   - SKU-001: Create ProductMaster (ItemWarehouse=8) + ProductRepair (ItemDamaged=2)
   - SKU-002: Create ProductMaster (ItemWarehouse=15)
   - SKU-003: Create ProductMaster (ItemWarehouse=19) + ProductRepair (ItemDamaged=1)

   Result:
   - 50 ProductMaster created (total 42 good items)
   - 3 ProductRepair created (total 3 damaged items)
   - Document status = "done"
   ```

---

## 📊 Perbandingan Keempat Proses

| Aspek              | Manual                  | Bulk                    | BAST               | SKU                  |
| ------------------ | ----------------------- | ----------------------- | ------------------ | -------------------- |
| **Input**          | Form individual         | File dengan mapping     | File + Scanner     | File Excel           |
| **Volume**         | 1-10 item               | 50-1000+ item           | 10-500 item        | 10-200 item          |
| **File Required**  | Tidak                   | Ya                      | Ya                 | Ya                   |
| **QC Validation**  | Manual saat input       | Otomatis kategori       | Per-item scanner   | Per-item crosscheck  |
| **Status Initial** | good/abnormal/damaged   | good                    | discrepancy        | good                 |
| **ProductMaster**  | Langsung created        | Langsung created        | Created saat scan  | Created saat finish  |
| **Location**       | staging_reguler/sticker | staging_reguler/sticker | Sesuai status scan | staging_sku          |
| **Special Field**  | -                       | -                       | DateScanned        | ItemGood/ItemDamaged |
| **Dokumen Type**   | manual                  | bulk                    | bast               | sku                  |

---

## 🔗 Hubungan Antar Entity

```
ProductDocument (type: manual|bulk|bast|sku)
    ├── 1:N → ProductPending (status: good|discrepancy|damaged|abnormal)
    │           ├── 1:1 → ProductMaster (location: staging_reguler|staging_sticker|damaged|abnormal|staging_sku|non)
    │           └── 0:1 → ProductRepair (untuk SKU damaged items)
    └── Metadata: HeaderMapping, Supplier, FileStats
```

---

## 🎯 Best Practices

1. **Inbound Manual**: Gunakan untuk input urgent atau volume kecil
2. **Inbound Bulk**: Optimal untuk supplier besar dengan data terstruktur
3. **Inbound BAST**: Best practice untuk QC ketat, per-item validation
4. **Inbound SKU**: Khusus untuk stok keeping unit dengan split quality

---

## 🔍 Debugging & Monitoring

### Query untuk melihat summary:

```sql
-- Bulk Summary
SELECT COUNT(*) as total_docs, SUM(file_item) as total_items, SUM(file_price) as total_price
FROM product_documents WHERE type='bulk';

-- BAST Summary
SELECT COUNT(*) as total_docs,
       SUM(CASE WHEN status='done' THEN 1 ELSE 0 END) as done_count,
       SUM(CASE WHEN status='progress' THEN 1 ELSE 0 END) as progress_count
FROM product_documents WHERE type='bast';

-- SKU Summary
SELECT COUNT(*) as total_docs,
       SUM(CASE WHEN status='done' THEN 1 ELSE 0 END) as done_count
FROM product_documents WHERE type='sku';

-- ProductPending by Status
SELECT status, COUNT(*) as total
FROM product_pendings
GROUP BY status;
```

---

## 📞 Support & Questions

Untuk pertanyaan atau clarification lebih lanjut, silakan refer ke:

- Controller files di `controller/`
- Service files di `services/`
- Model definitions di `models/`
- Routes di `routes/routes.go`
