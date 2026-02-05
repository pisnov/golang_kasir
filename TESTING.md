# Testing Guide - Golang Kasir API

## Prerequisites
- Database Supabase sudah memiliki tables: `products`, `categories`, `transactions`, `transaction_details`
- Project sudah di-build dengan sukses
- Server berjalan di `http://localhost:8080`

---

## 🔍 Cara Mendapatkan Data Produk Anda

Sebelum melakukan testing, pastikan Anda tahu produk apa saja yang ada di database:

### Opsi 1: Via API (Recommended)
```bash
# Pastikan server berjalan terlebih dahulu
go run main.go

# Di terminal lain, jalankan:
curl http://localhost:8080/api/produk
```

### Opsi 2: Via SQL Query di Supabase Dashboard
```sql
SELECT p.id, p.name, p.price, p.stock, c.name as category_name
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.id;
```

### Opsi 3: Jalankan Script Helper
```bash
go run check_products.go
```

**⚠️ PENTING:** Ganti semua `product_id`, `nama_produk`, dan nilai lainnya dalam contoh testing dengan data produk ASLI dari database Anda!

---

## 1. Search Product by Name

### Get All Products
```bash
curl -X GET "http://localhost:8080/api/produk"
```

### Search by Name (case-insensitive)
```bash
# Ganti "nama_produk" dengan nama produk yang ada di database Anda
curl -X GET "http://localhost:8080/api/produk?name=nama_produk"
```

### Response Example
```json
[
  {
    "id": 1,
    "name": "Produk A",
    "price": 10000,
    "stock": 50,
    "category_name": "Kategori A"
  },
  {
    "id": 2,
    "name": "Produk B",
    "price": 15000,
    "stock": 30,
    "category_name": "Kategori A"
  }
]
```

### Notes
- Parameter `name` bersifat case-insensitive (menggunakan ILIKE)
- Pencarian partial (contoh: mencari "pro" akan menemukan "Produk A", "Produk B", dll)
- Jika tidak ada parameter `name`, akan menampilkan semua produk

---

## 2. Create Transaction (Checkout)

### Request
```bash
# PENTING: Ganti product_id dengan ID produk yang ada di database Anda
# Jalankan GET /api/produk terlebih dahulu untuk melihat ID produk yang tersedia

curl -X POST "http://localhost:8080/api/checkout" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

### Response (Success - 201)
```json
{
  "id": 1,
  "total_amount": 35000,
  "created_at": "2026-02-05T10:30:45.123Z",
  "details": [
    {
      "id": 0,
      "transaction_id": 1,
      "product_id": 1,
      "product_name": "Produk A",
      "quantity": 2,
      "subtotal": 20000
    },
    {
      "id": 0,
      "transaction_id": 1,
      "product_id": 2,
      "product_name": "Produk B",
      "quantity": 1,
      "subtotal": 15000
    }
  ]
}
```

### Error Responses

**Product not found (500)**
```json
"product id 999 not found"
```

**Insufficient stock (500)**
```json
"insufficient stock for product [Nama Produk]"
```

### How It Works
1. **Validasi Stock**: Mengecek apakah stock produk mencukupi
2. **Hitung Total**: Menghitung subtotal untuk setiap item dan total keseluruhan
3. **Database Transaction**: Menggunakan BEGIN/COMMIT/ROLLBACK untuk konsistensi data
4. **Update Stock**: Mengurangi stock produk dari database
5. **Insert Transaction**: Menyimpan transaksi dan detail ke database

---

## 3. Get Sales Summary - Today

### Request
```bash
curl -X GET "http://localhost:8080/api/report"
```

### Response Example
```json
{
  "total_revenue": 125000,
  "total_transaksi": 8,
  "produk_terlaris": {
    "nama": "Produk A",
    "qty_terjual": 25
  }
}
```

**Note:** Jika belum ada transaksi hari ini, `produk_terlaris` akan bernilai `null`

---

## 4. Get Sales Summary - By Date Range

### Request
```bash
curl -X GET "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-02-05"
```

### Response Example
```json
{
  "total_revenue": 450000,
  "total_transaksi": 52,
  "produk_terlaris": {
    "nama": "Produk A",
    "qty_terjual": 120
  }
}
```

### Notes
- Format date: `YYYY-MM-DD`
- Query aggregates data untuk range yang di-request
- Jika tidak ada transaksi dalam range, `produk_terlaris` akan `null`
- Data akan mencakup semua transaksi dari start_date sampai end_date (inclusive)

---

## Database Queries untuk Debugging

### Check transactions hari ini
```sql
SELECT * FROM transactions 
WHERE DATE(created_at) = CURRENT_DATE;
```

### Check transaction details
```sql
SELECT td.*, p.name 
FROM transaction_details td
JOIN products p ON td.product_id = p.id
WHERE td.transaction_id = 15;
```

### Check current product stock
```sql
SELECT id, name, price, stock 
FROM products 
ORDER BY id;
```

### Create sample data (jika diperlukan)
```sql
-- Buat kategori terlebih dahulu
INSERT INTO categories (name) VALUES
('Kategori A'),
('Kategori B'),
('Kategori C');

-- Buat produk dengan kategori
-- Sesuaikan nama produk, harga, dan stock dengan kebutuhan Anda
INSERT INTO products (name, price, stock, category_id) VALUES
('Produk A', 10000, 100, 1),
('Produk B', 15000, 80, 1),
('Produk C', 20000, 150, 2),
('Produk D', 25000, 120, 2),
('Produk E', 30000, 30, 3);
```

---

## Testing Checklist

- [ ] Search product tanpa parameter (tampil semua)
- [ ] Search product dengan parameter name (case insensitive)
- [ ] Checkout dengan 1 item
- [ ] Checkout dengan multiple items
- [ ] Checkout dengan product yang tidak ada (error)
- [ ] Checkout dengan quantity melebihi stock (error)
- [ ] View sales summary hari ini
- [ ] View sales summary dengan date range
- [ ] Verify stock berkurang setelah checkout
- [ ] Verify transaction & transaction_details tersimpan di database

---

## Troubleshooting

### Error: "product id X not found"
- Pastikan product ID ada di database
- Gunakan `/api/produk` untuk melihat daftar product yang tersedia

### Error: "insufficient stock for product X"
- Check current stock menggunakan query di atas
- Kurangi quantity yang di-checkout

### Error: Database connection failed
- Verify DB_CONN environment variable
- Check Supabase connection string
- Pastikan tables sudah di-create

### No response dari server
- Verify server berjalan: `curl http://localhost:8080/api/produk`
- Check port tidak digunakan oleh aplikasi lain
- Lihat logs untuk error messages

---

## 💡 Best Practices untuk Testing

### 1. Selalu Cek Produk Yang Ada Terlebih Dahulu
```bash
# Sebelum testing apapun, jalankan ini dulu
curl http://localhost:8080/api/produk | jq
```

### 2. Simpan Response untuk Reference
```bash
# Simpan daftar produk ke file
curl http://localhost:8080/api/produk > my_products.json
```

### 3. Buat Script Testing Sendiri
Buat file `test_my_data.sh`:
```bash
#!/bin/bash

# Ambil produk pertama dan kedua
PRODUCT_1_ID=1  # Ganti dengan ID produk Anda
PRODUCT_2_ID=2  # Ganti dengan ID produk Anda

echo "Testing checkout dengan produk ID $PRODUCT_1_ID dan $PRODUCT_2_ID"
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d "{
    \"items\": [
      {\"product_id\": $PRODUCT_1_ID, \"quantity\": 1},
      {\"product_id\": $PRODUCT_2_ID, \"quantity\": 1}
    ]
  }"
```

### 4. Testing Flow yang Disarankan
1. **GET** semua produk → catat ID yang ada
2. **Search** by name → test pencarian
3. **POST** checkout dengan ID yang valid → test transaksi
4. **GET** report → lihat hasil transaksi
5. **GET** produk lagi → verify stock berkurang

### 5. Gunakan Postman atau Thunder Client
Untuk testing yang lebih mudah, import collection ini atau buat sendiri di VS Code extension Thunder Client.

---

## 📝 Template Request dengan Data Anda

Copy template ini dan sesuaikan dengan data Anda:

```bash
# 1. CEK PRODUK YANG ADA
curl http://localhost:8080/api/produk

# 2. SEARCH SPESIFIK (ganti "nama_produk_anda")
curl "http://localhost:8080/api/produk?name=nama_produk_anda"

# 3. CHECKOUT (ganti ID dengan ID produk Anda)
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": GANTI_DENGAN_ID_ANDA, "quantity": 1}
    ]
  }'

# 4. LIHAT REPORT
curl http://localhost:8080/api/report
```

**Tip:** Save template di atas ke file `my_testing_commands.txt` untuk referensi cepat!
