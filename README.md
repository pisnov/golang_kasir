# Golang Kasir — API Kasir (Layered Architecture)

Aplikasi API Kasir yang dibangun mengikuti pola **Layered Architecture** (Handler → Service → Repository → Database). Implementasi ini menggunakan PostgreSQL (Supabase) untuk penyimpanan data dan `viper` untuk konfigurasi environment.

## 🎯 Fitur Utama

### Product Management
- ✅ CRUD Produk (`CREATE`, `READ`, `UPDATE`, `DELETE`)
- ✅ Search Produk by Name (case-insensitive, partial matching)
- ✅ Manajemen Stock

### Category Management
- ✅ CRUD Kategori produk
- ✅ Link kategori dengan produk

### Transaction System
- ✅ Checkout/Create Transaksi multi-item
- ✅ Database transaction dengan rollback support
- ✅ Validasi stock otomatis
- ✅ Simpan detail transaksi

### Sales Report
- ✅ Sales summary hari ini
- ✅ Sales summary dengan date range
- ✅ Identifikasi produk terlaris

## 🏗️ Struktur Arsitektur

Implementasi mengikuti Layered Architecture:
- `models` — Struktur data (Product, Category, Transaction)
- `repositories` — Database queries & operations
- `services` — Business logic
- `handlers` — HTTP request/response
- `database` — Database connection

## 📋 Prerequisites

- **Go** 1.20+ (sudah terpasang)
- **PostgreSQL** (Supabase recommended)
- **Git** untuk version control

## ⚙️ Setup Awal

### 1. Clone & Setup Project

```bash
git clone <repository-url>
cd golang_kasir
```

### 2. Konfigurasi Environment

Buat atau update file `.env` di root project:

```dotenv
PORT=8080
DB_CONN=postgresql://postgres:PASSWORD@your-project.supabase.co:6543/postgres
```

Ganti `PASSWORD` dan host sesuai kredensial Supabase/PostgreSQL Anda.

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Build Project

```bash
go build -o kasir-api .
```

## 🗄️ Database Schema

Execute SQL berikut di PostgreSQL/Supabase:

```sql
-- Categories Table
CREATE TABLE IF NOT EXISTS categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products Table
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  price INT NOT NULL,
  stock INT NOT NULL,
  category_id INT REFERENCES categories(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  total_amount INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transaction Details Table
CREATE TABLE IF NOT EXISTS transaction_details (
  id SERIAL PRIMARY KEY,
  transaction_id INT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  product_id INT NOT NULL REFERENCES products(id),
  quantity INT NOT NULL,
  subtotal INT NOT NULL
);
```

### Sample Data (Optional)

Untuk testing, Anda bisa menambahkan sample data:

```sql
-- Insert sample categories
INSERT INTO categories (name) VALUES
('Kategori A'),
('Kategori B'),
('Kategori C');

-- Insert sample products
-- Sesuaikan dengan kebutuhan bisnis Anda
INSERT INTO products (name, price, stock, category_id) VALUES
('Produk 1', 10000, 100, 1),
('Produk 2', 15000, 80, 1),
('Produk 3', 20000, 150, 2),
('Produk 4', 25000, 120, 2),
('Produk 5', 30000, 50, 3);
```

## 🔍 Cara Melihat Data Produk Anda

Sebelum testing API, pastikan Anda tahu produk apa yang ada di database:

### Metode 1: Via API (Recommended)
```bash
# Jalankan server terlebih dahulu
go run main.go

# Di terminal lain:
curl http://localhost:8080/api/produk | jq
```

### Metode 2: Via Supabase Dashboard
- Buka Supabase Dashboard → Table Editor
- Pilih table `products`
- Lihat semua data produk Anda

### Metode 3: Via SQL Query
```sql
SELECT p.id, p.name, p.price, p.stock, c.name as category_name
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.id;
```

**💡 Tip:** Gunakan ID dan nama produk dari database Anda saat testing API endpoints!

## 🚀 Menjalankan Server

### Development Mode
```bash
go run main.go
```

### Production Mode (menggunakan binary)
```bash
./kasir-api
```

Server akan berjalan pada `http://0.0.0.0:<PORT>` (default `PORT=8080`)

## 📡 API Endpoints

> **⚠️ PENTING:** Semua contoh di bawah menggunakan data template. 
> Ganti `product_id`, `nama_produk`, dan nilai lainnya dengan data ASLI dari database Anda.
> Lihat section [Cara Melihat Data Produk Anda](#-cara-melihat-data-produk-anda) untuk mendapatkan data produk Anda.

### Product Management

#### Get All Products (dengan search optional)
```bash
# Ambil semua produk
GET /api/produk

# Search produk by name (case-insensitive)
# Ganti "nama_produk" dengan produk yang ada di database Anda
GET /api/produk?name=nama_produk
```

**Response Example:**
```json
[
  {
    "id": 1,
    "name": "Nama Produk Anda",
    "price": 10000,
    "stock": 50,
    "category_name": "Nama Kategori"
  },
  {
    "id": 2,
    "name": "Produk Lainnya",
    "price": 15000,
    "stock": 30,
    "category_name": "Nama Kategori"
  }
]
```

#### Create Product
```bash
POST /api/produk
Content-Type: application/json

{
  "name": "Nama Produk Baru",
  "price": 10000,
  "stock": 50,
  "category_id": 1
}
```

#### Get Product by ID
```bash
GET /api/produk/{id}
```

#### Update Product
```bash
# Ganti {id} dengan ID produk yang ingin diupdate
PUT /api/produk/{id}
Content-Type: application/json

{
  "name": "Nama Produk Updated",
  "price": 15000,
  "stock": 20,
  "category_id": 1
}
```

#### Delete Product
```bash
DELETE /api/produk/{id}
```

### Category Management

#### Get All Categories
```bash
GET /api/kategori
```

#### Create Category
```bash
POST /api/kategori
Content-Type: application/json

{
  "name": "Nama Kategori Baru"
}
```

#### Get Category by ID
```bash
GET /api/kategori/{id}
```

#### Update Category
```bash
# Ganti {id} dengan ID kategori yang ingin diupdate
PUT /api/kategori/{id}
Content-Type: application/json

{
  "name": "Nama Kategori Updated"
}
```

#### Delete Category
```bash
DELETE /api/kategori/{id}
```

### Transaction/Checkout

#### Create Transaction (Checkout)
```bash
# PENTING: Ganti product_id dengan ID produk yang ada di database Anda
# Jalankan GET /api/produk untuk melihat ID yang tersedia
POST /api/checkout
Content-Type: application/json

{
  "items": [
    {"product_id": 1, "quantity": 2},
    {"product_id": 2, "quantity": 1}
  ]
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "total_amount": 35000,
  "created_at": "2026-02-05T10:30:45Z",
  "details": [
    {
      "id": 0,
      "transaction_id": 1,
      "product_id": 1,
      "product_name": "Nama Produk A",
      "quantity": 2,
      "subtotal": 20000
    },
    {
      "id": 0,
      "transaction_id": 1,
      "product_id": 2,
      "product_name": "Nama Produk B",
      "quantity": 1,
      "subtotal": 15000
    }
  ]
}
```

### Sales Report

#### Get Today's Sales Summary
```bash
GET /api/report
```

#### Get Sales Summary by Date Range
```bash
GET /api/report?start_date=2026-01-01&end_date=2026-02-05
```

**Response Example:**
```json
{
  "total_revenue": 125000,
  "total_transaksi": 8,
  "produk_terlaris": {
    "nama": "Nama Produk Terlaris Anda",
    "qty_terjual": 25
  }
}
```

## 📂 Struktur File Project

```
golang_kasir/
├── main.go                          # Entry point, routing, konfigurasi
├── go.mod                          # Go module definition
├── go.sum                          # Go dependencies lock file
├── .env                            # Environment variables (ignore di git)
├── README.md                       # Dokumentasi
├── IMPLEMENTATION_NOTES.md         # Detail implementasi fitur
├── TESTING.md                      # Panduan testing API
│
├── database/
│   └── database.go                 # Database connection & initialization
│
├── models/
│   ├── product.go                  # Product struct
│   ├── category.go                 # Category struct
│   └── transaction.go              # Transaction & TransactionDetail structs
│
├── repositories/
│   ├── product_repository.go       # Product CRUD & queries
│   ├── category_repository.go      # Category CRUD & queries
│   └── transaction_repository.go   # Transaction creation & reports
│
├── services/
│   ├── product_service.go          # Product business logic
│   ├── category_service.go         # Category business logic
│   └── transaction_service.go      # Transaction business logic
│
└── handlers/
    ├── product_handler.go          # HTTP handlers untuk product
    ├── category_handler.go         # HTTP handlers untuk category
    ├── transaction_handler.go      # HTTP handlers untuk checkout
    └── report_handler.go           # HTTP handlers untuk sales report
```

## 💡 Tips & Notes

- **Database Connection**: Pastikan `DB_CONN` format sesuai dengan provider Anda
  - Supabase: gunakan connection string dari dashboard
  - Local PostgreSQL: `postgresql://user:password@localhost:5432/dbname`
  
- **Search Feature**: Search produk menggunakan `ILIKE` yang bersifat case-insensitive dan support partial matching
  
- **Transaction Safety**: Checkout menggunakan SQL transaction untuk memastikan konsistensi data (all or nothing)
  
- **Stock Management**: Stock otomatis berkurang saat checkout, dengan validasi untuk mencegah overselling
  
- **Date Format**: Format tanggal untuk report adalah `YYYY-MM-DD`

## 📖 Dokumentasi Tambahan

- [IMPLEMENTATION_NOTES.md](IMPLEMENTATION_NOTES.md) — Detail implementasi fitur
- [TESTING.md](TESTING.md) — Panduan testing API dengan curl examples

## 🔄 Development Workflow

```bash
# 1. Update kode
# 2. Build & test
go run main.go

# 3. Commit & push
git add .
git commit -m "feat: deskripsi fitur"
git push
```

## 📝 License

MIT

