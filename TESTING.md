# Testing Guide - Golang Kasir API

## Prerequisites
- Database Supabase sudah memiliki tables: `products`, `categories`, `transactions`, `transaction_details`
- Project sudah di-build dengan sukses
- Server berjalan di `http://localhost:8080`

---

## 1. Search Product by Name

### Request
```bash
curl -X GET "http://localhost:8080/api/produk?name=indomie"
```

### Response
```json
[
  {
    "id": 1,
    "name": "Indomie Goreng",
    "price": 3500,
    "stock": 50,
    "category_name": "Mie"
  },
  {
    "id": 4,
    "name": "Indomie Kuah",
    "price": 3500,
    "stock": 35,
    "category_name": "Mie"
  }
]
```

### Notes
- Parameter `name` bersifat case-insensitive (menggunakan ILIKE)
- Pencarian partial (contoh: "indom" akan menampilkan semua produk dengan "indom" di nama)
- Jika tidak ada parameter `name`, akan menampilkan semua produk

---

## 2. Create Transaction (Checkout)

### Request
```bash
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
  "id": 15,
  "total_amount": 15500,
  "created_at": "2026-02-05T10:30:45.123Z",
  "details": [
    {
      "id": 0,
      "transaction_id": 15,
      "product_id": 1,
      "product_name": "Indomie Goreng",
      "quantity": 2,
      "subtotal": 7000
    },
    {
      "id": 0,
      "transaction_id": 15,
      "product_id": 2,
      "product_name": "Coca Cola",
      "quantity": 1,
      "subtotal": 8500
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
"insufficient stock for product Indomie Goreng"
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

### Response
```json
{
  "total_revenue": 125000,
  "total_transaksi": 8,
  "produk_terlaris": {
    "nama": "Indomie Goreng",
    "qty_terjual": 25
  }
}
```

---

## 4. Get Sales Summary - By Date Range

### Request
```bash
curl -X GET "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-02-05"
```

### Response
```json
{
  "total_revenue": 450000,
  "total_transaksi": 52,
  "produk_terlaris": {
    "nama": "Indomie Goreng",
    "qty_terjual": 120
  }
}
```

### Notes
- Format date: `YYYY-MM-DD`
- Query aggregates data untuk range yang di-request
- Jika tidak ada transaksi dalam range, `produk_terlaris` akan `null`

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
INSERT INTO products (name, price, stock, category_id) VALUES
('Indomie Goreng', 3500, 100, 1),
('Indomie Kuah', 3500, 80, 1),
('Coca Cola', 8500, 150, 2),
('Sprite', 8500, 120, 2),
('Roti Tawar', 25000, 30, 3);
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

