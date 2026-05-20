# Inbound Process - Visual Diagrams

## 1. INBOUND MANUAL - Flow Diagram

```
┌─────────────────────┐
│  Form Input Manual  │
│ (user, item, price) │
└──────────┬──────────┘
           │
           ▼
┌────────────────────────────────┐
│   Validate Input               │
│   - Harga valid               │
│   - Item > 0                  │
└──────────┬─────────────────────┘
           │
           ▼
┌────────────────────────────────┐
│  Generate Barcode              │
│  BC-{timestamp}-{random}       │
└──────────┬─────────────────────┘
           │
           ▼
┌────────────────────────────────┐      ┌─────────────────────┐
│  Tentukan Kategori/Sticker     │      │  Logic Penentuan:   │
│  Berdasarkan Harga & Input     │◄────►│  • Harga >= 100rb:  │
└──────────┬─────────────────────┘      │    categories       │
           │                            │  • Harga < 100rb:   │
           │                            │    sticker + range  │
           ▼                            └─────────────────────┘
┌────────────────────────────────┐
│  Hitung Price Warehouse        │
│  = Price × (1 - Discount%)     │
│    atau Fixed Price Sticker    │
└──────────┬─────────────────────┘
           │
     ┌─────┴─────┐
     │           │
     ▼           ▼
┌──────────┐  ┌──────────────┐
│ INSERT   │  │ INSERT       │
│ Pending  │  │ ProductMaster│
│          │  │              │
│ barcode  │  │ barcode_warehouse
│ status   │  │ category_id  │
│ price    │  │ price_warehouse
│ item     │  │ location     │
└────┬─────┘  └──────┬───────┘
     │               │
     └───────┬───────┘
             │
             ▼
    ┌──────────────────┐
    │ Response         │
    │ Success: Created │
    └──────────────────┘
```

## 2. INBOUND BULK - Flow Diagram

```
┌─────────────────────────────┐
│  Upload File (CSV/XLSX/XLS) │
│  + supplier                 │
│  + type_product (reguler/   │
│    sticker)                 │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│  Parse File & Extract Rows  │
│  + Mapping header columns   │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐      ┌─────────────────┐
│  Hitung Total Item & Price  │      │  Result:        │
│  Sum dari kolom Qty & Price │      │  - TotalItem    │
└──────────┬──────────────────┘      │  - TotalPrice   │
           │                         └─────────────────┘
           ▼
┌─────────────────────────────┐
│  CREATE ProductDocument     │
│  type="bulk"                │
│  status="progress"          │
│  code=BULK-{timestamp}      │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│  LOOP setiap ROW dari file:         │
└──────────┬────────────────────────┬─┘
           │                        │
      ┌────▼──────┐           ┌─────▼───────┐
      │ Type:     │           │ Type:       │
      │ Reguler   │           │ Sticker     │
      │ Harga >=  │           │ Harga <     │
      │ 100rb     │           │ 100rb       │
      └────┬──────┘           └─────┬───────┘
           │                        │
           ▼                        ▼
    ┌─────────────┐       ┌──────────────┐
    │ Validate:   │       │ Validate:    │
    │ • Kategori  │       │ • Sticker    │
    │   exist     │       │   range      │
    │ • Harga >=  │       │   match      │
    │   100rb     │       │ • Harga <    │
    └────┬────────┘       │   100rb      │
         │                └──────┬───────┘
         │                       │
    ┌────▼────────┐        ┌────▼────────┐
    │ SKIP        │        │ SKIP        │
    │ Count++     │        │ Count++     │
    └─────────────┘        └─────────────┘
              │                    │
              └────────┬───────────┘
                       │
                       ▼
              ┌──────────────────┐
              │ INSERT Pending   │
              │ INSERT Master    │
              │ (calc price_wh)  │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │ Set Location:    │
              │ staging_reguler/ │
              │ staging_sticker  │
              └──────────────────┘
```

## 3. INBOUND BAST - Flow Diagram

```
┌──────────────────────────────────┐
│  1. UPLOAD FILE BAST             │
│     (CSV/XLSX/XLS)               │
│     + Supplier                   │
│     + Header Mapping             │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│  Parse File & Extract Rows       │
│  + Validate header mapping       │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│  CREATE ProductDocument          │
│  type="bast"                     │
│  status="progress"               │
│  code=BAST-{timestamp}           │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│  Loop setiap row → INSERT Pending    │
│  status="discrepancy" (belum disc)   │
│  date_scanned=NULL                   │
└──────────┬───────────────────────────┘
           │
           ▼
    ┌──────────────────────────────┐
    │ 2. GET DOCUMENT INFO         │
    │    (Scanner Interface)       │
    │    Display: unscanned count  │
    │             scanned count    │
    └──────────┬───────────────────┘
               │
    ┌──────────▼───────────────────┐
    │ 3. PER-ITEM SCANNER          │
    │    Scan Barcode              │
    │    ↓                         │
    │    GET Pending by Barcode    │
    │    (cek status != scanned)   │
    └──────────┬───────────────────┘
               │
               ▼
    ┌──────────────────────────────┐
    │ 4. QC INPUT                  │
    │    - Select category         │
    │    - Select status:          │
    │      good/abnormal/damaged   │
    │    - Input note              │
    └──────────┬───────────────────┘
               │
       ┌───────┼───────────────────────┐
       │       │                       │
       ▼       ▼                       ▼
    ┌──┐   ┌──┐                    ┌──┐
    │good  │abnormal               │damaged
    │      │                       │non
    └──┬┘  └──┬┘                   └──┬┘
       │      │                       │
       ▼      ▼                       ▼
    staging_ abnormal              damaged/
    reguler/                        non
    sticker
       │      │                       │
       └──────┴───────────┬──────────┘
                          │
                          ▼
              ┌─────────────────────────┐
              │ Calculate price_        │
              │ warehouse (dari         │
              │ kategori/sticker logic) │
              └───────────┬─────────────┘
                          │
                          ▼
              ┌─────────────────────────┐
              │ INSERT ProductMaster    │
              │ - barcode_warehouse     │
              │ - category_id (input)   │
              │ - price_warehouse       │
              │ - location (from status)│
              │ - product_pending_id    │
              └───────────┬─────────────┘
                          │
                          ▼
              ┌─────────────────────────┐
              │ UPDATE ProductPending   │
              │ - date_scanned=NOW()    │
              │ - status=scanned        │
              └───────────┬─────────────┘
                          │
                          ▼
            ┌──────────────────────────┐
            │ Response: Master Data    │
            │ + price_warehouse info   │
            └──────────────────────────┘
               │
    Continue loop atau
    next item → kembali ke
    GET Pending by Barcode
               │
               ▼
    ┌──────────────────────────────┐
    │ 5. FINISH DOCUMENT           │
    │    Update Document status    │
    │    = "done"                  │
    └──────────────────────────────┘
```

## 4. INBOUND SKU - Flow Diagram

```
┌──────────────────────────────┐
│  1. UPLOAD EXCEL SKU FILE    │
│     + Supplier               │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Parse Excel File            │
│  Extract columns:            │
│  [0]=barcode                 │
│  [1]=name                    │
│  [3]=price                   │
│  [4]=item (initial qty)      │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  CREATE ProductDocument      │
│  type="sku"                  │
│  status="pending"            │
│  code=SKU-{timestamp}        │
│  supplier={input}            │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│  Loop setiap row →               │
│  CREATE ProductPending           │
│  is_sku=true                     │
│  item=file_qty (initial)         │
│  item_good=NULL (will fill)      │
│  item_damaged=NULL (will fill)   │
└──────────┬───────────────────────┘
           │
           ▼
    ┌────────────────────────────────┐
    │ 2. QC STAFF CROSSCHECK         │
    │    (per-item split validation) │
    └────────┬─────────────────────┬─┘
             │                     │
        ┌────▼──────────┐    ┌─────▼─────┐
        │ GET Document  │    │ Per Pending│
        │ Info          │    │ Input:    │
        │               │    │ item_good │
        └───────────────┘    │ + damaged │
                             └─────┬─────┘
                                   │
                                   ▼
                        ┌──────────────────────┐
                        │ VALIDATE:            │
                        │ item_good +          │
                        │ item_damaged         │
                        │ <= total item        │
                        └──────────┬───────────┘
                                   │
                        ┌──────────▼────────┐
                        │ UPDATE Pending:   │
                        │ item_good = input │
                        │ item_damaged =    │
                        │ input             │
                        └──────────┬────────┘
                                   │
        Continue loop next pending  │
        atau finish jika semua done │
                                   │
                                   ▼
                    ┌───────────────────────┐
                    │ 3. FINISH INBOUND SKU │
                    │    (auto process)     │
                    └─────────┬─────────────┘
                              │
                 ┌────────────▼─────────────┐
                 │ LOOP all Pendings:      │
                 └────────────┬─────────────┘
                              │
                    ┌─────────▼──────────┐
                    │ Jika item_good > 0 │
                    └─────────┬──────────┘
                              │
                    ┌─────────▼──────────────┐
                    │ CREATE ProductMaster   │
                    │ item_warehouse=good    │
                    │ location="staging_sku" │
                    │ is_sku=true            │
                    └──────────────────────┘
                              │
                              ▼
                    ┌──────────────────────┐
                    │ Jika damaged > 0     │
                    └──────────┬───────────┘
                               │
                    ┌──────────▼────────────┐
                    │ CREATE ProductRepair  │
                    │ item_damaged=input    │
                    │ product_pending_id ref│
                    └──────────┬────────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │ UPDATE Document      │
                    │ status="done"        │
                    └──────────────────────┘
```

## 5. ENTITY RELATIONSHIP - Data Model

```
┌─────────────────────────────┐
│  ProductDocument            │
├─────────────────────────────┤
│ id (PK)                     │
│ code                        │
│ type: manual|bulk|bast|sku  │ ◄── Defines Inbound Type
│ status: progress|done       │
│ supplier                    │
│ file_item, file_price       │
│ header_* (mapping)          │
└──────────────┬──────────────┘
               │ 1:N
               │
        ┌──────▼──────────────┐         ┌──────────────────┐
        │ ProductPending      │         │ ProductMaster    │
        ├─────────────────────┤         ├──────────────────┤
        │ id (PK)             │1:1     │ id (PK)          │
        │ document_id (FK) ────┼───────► document_id (FK) │
        │ barcode             │        │ barcode_warehouse│
        │ name, item, price   │        │ category_id (FK) │
        │ status: good|       │        │ sticker_id (FK)  │
        │   discrepancy|      │        │ price_warehouse  │
        │   damaged|abnormal  │        │ location         │
        │ is_sku              │        │ is_sku           │
        │ item_good (SKU)     │        │ item_warehouse   │
        │ item_damaged (SKU)  │        │ product_pending  │
        │ date_scanned (BAST) │        │ _id (FK)         │
        └──────┬──────────────┘        └────────┬─────────┘
               │                              │
               │ 0:1 (SKU only)              │
               │                             │
        ┌──────▼──────────────────┐         │
        │ ProductRepair           │         │
        ├─────────────────────────┤         │
        │ id (PK)                 │         │
        │ product_pending_id (FK) ├────┐    │
        │ item_damaged            │    │    │
        │ (for damaged items)     │    │    │
        └─────────────────────────┘    │    │
                                       │    │
        ┌──────────────────────────┐   │    │
        │ Category                 │◄──┼────┤
        ├──────────────────────────┤   │    │
        │ id (PK)                  │   │    │
        │ name                     │   │    │
        │ discount (%)             │   │    │
        └──────────────────────────┘   │    │
                                       │    │
        ┌──────────────────────────┐   │    │
        │ Sticker                  │◄──┼────┤
        ├──────────────────────────┤   │    │
        │ id (PK)                  │   │    │
        │ name                     │   │    │
        │ min_price, max_price     │   │    │
        │ fixed_price              │   │    │
        └──────────────────────────┘   │    │
                                       │    │
                                       └────┘
```

## 6. LOCATION MAPPING - Product Journey

```
┌────────────────────────────────────────────────────────────────┐
│                    INBOUND PROCESS                             │
└────────────────────────────────────────────────────────────────┘

┌─────────────────┐
│   Manual Input  │ → status=good      → location=staging_reguler/sticker
│   or Upload     │ → status=abnormal  → location=abnormal
│   (all types)   │ → status=damaged   → location=damaged
└─────────────────┘ → status=non       → location=non

          ↓

┌─────────────────────────────────────────────────────────────────┐
│             Categorization / Type Determination                 │
│                                                                 │
│  Price >= 100rb → Category → Discount Applied → staging_reguler  │
│  Price < 100rb  → Sticker  → Fixed Price      → staging_sticker │
│                                                                 │
│  (OR explicit status overrides → abnormal/damaged/non)         │
└─────────────────────────────────────────────────────────────────┘

          ↓

┌─────────────────────────────────────────────────────────────────┐
│              PRODUCT READY FOR WAREHOUSE OPERATIONS             │
│                                                                 │
│  Location: staging_reguler ──► [can rack display]              │
│  Location: staging_sticker  ──► [can bag/wholesale]            │
│  Location: staging_sku      ──► [SKU specific area]            │
│  Location: abnormal         ──► [QC area]                      │
│  Location: damaged          ──► [Return/Repair area]           │
│  Location: non              ──► [Non-sellable area]            │
└─────────────────────────────────────────────────────────────────┘
```

## 7. STATUS FLOW - Document & Product Status

```
DOCUMENT STATUS FLOW:

              ┌──────────────┐
              │   Created    │
              └──────┬───────┘
                     │
    ┌────────────────┼────────────────┐
    ▼                ▼                ▼
┌─────────┐      ┌─────────┐      ┌─────────┐
│ Manual  │      │ Bulk    │      │ BAST    │ ← Scanner
│ type    │      │ type    │      │ type    │   Per-item
│ DONE    │      │ DONE    │      │ Progress│
│ always  │      │ always  │      │ ──→Done │
└─────────┘      └─────────┘      └─────────┘

                     ▼
              ┌──────────────┐
              │   Archived   │
              │   (history)  │
              └──────────────┘


PRODUCT STATUS FLOW (in Pending):

              ┌──────────────┐
              │   Created    │
              │ (discrepancy │  ← for BAST
              │  or good)    │
              └──────┬───────┘
                     │
    ┌────────┬───────┼─────────┬────────┐
    │        │       │         │        │
    ▼        ▼       ▼         ▼        ▼
┌─────┐  ┌──────┐ ┌─────┐ ┌────────┐ ┌────┐
│good │  │abnorm│ │damag│ │discrepan│ │non │
│     │  │al    │ │ed   │ │cy (BAST)│ │    │
└─────┘  └──────┘ └─────┘ │──┬─────┘ └────┘
   │        │        │       │
   │        │        │       │ (after scan)
   │        │        │       ▼
   │        │        │    ┌─────┐
   │        │        │    │other│  ← (good/abnormal/damaged/non)
   │        │        │    └─────┘
   └────────┴────────┴────┬────┘
                          │
                          ▼
                  ProductMaster
                  + LocationSet
```

---

## 8. Priority & Decision Tree

```
INPUT BARCODE/NAME/PRICE
│
├─ Harga >= 100.000
│  ├─ Ada Category Input?
│  │  ├─ Yes → Use Category (get discount)
│  │  └─ No  → Skip (need category)
│  └─ Location = staging_reguler
│
├─ Harga < 100.000
│  ├─ Ada Sticker Input?
│  │  ├─ Yes → Use Sticker (get fixed price)
│  │  └─ No  → Auto-match by range
│  └─ Location = staging_sticker
│
└─ Status Override
   ├─ Status = damaged → location = damaged
   ├─ Status = abnormal → location = abnormal
   └─ Status = non → location = non
```

