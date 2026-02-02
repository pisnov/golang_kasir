go run main.go
# Golang Kasir — API Kasir (Layered Architecture)

Proyek sederhana API kasir yang mengikuti tutorial Layered Architecture (Handler → Service → Repository → Database). Implementasi ini menggunakan PostgreSQL (mis. Supabase) untuk penyimpanan data, dan `viper` untuk konfigurasi.

Fitur utama
- CRUD produk (`products` table)
- Struktur berlapis: `models`, `repositories`, `services`, `handlers`, `database`

Prerequisites
- Go 1.20+ terpasang
- Database PostgreSQL (Supabase recommended)

Persiapan & Konfigurasi

1. Salin repository dan masuk ke folder project:

```bash
git clone <repo>
cd golang_kasir
```

2. Isi file `.env` di root project (contoh sudah ada `.env`):

```dotenv
PORT=8080
DB_CONN=postgresql://postgres:PASSWORD@your-project.supabase.co:6543/postgres
```

Ganti `PASSWORD` dan host sesuai kredensial Supabase / Postgres Anda.

3. Install dependency dan build:

```bash
go mod tidy
go build ./...
```

Database

Contoh skema table `products`:

```sql
CREATE TABLE products (
  id serial PRIMARY KEY,
  name varchar NOT NULL,
  price int NOT NULL,
  stock int NOT NULL
);
```

Jalankan SQL di atas di database Supabase Anda.

Menjalankan server

```bash
go run main.go
# atau
./<binary>
```

Server akan berjalan pada `0.0.0.0:<PORT>`.

API Endpoints

- GET /api/produk
  - Mengambil semua produk
  - contoh:

```bash
curl http://localhost:8080/api/produk
```

- POST /api/produk
  - Membuat produk baru. Body JSON:

```json
{
  "name": "Nasi Goreng",
  "price": 25000,
  "stock": 10
}
```

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"Nasi Goreng","price":25000,"stock":10}' http://localhost:8080/api/produk
```

- GET /api/produk/{id}
  - Ambil produk berdasarkan ID

- PUT /api/produk/{id}
  - Update produk (kirim JSON body seperti POST)

- DELETE /api/produk/{id}
  - Hapus produk

Struktur Project

- `main.go` — inisialisasi, routing, dan konfigurasi `viper`
- `database/database.go` — koneksi DB
- `models/product.go` — tipe `Product`
- `repositories/product_repository.go` — query SQL CRUD
- `services/product_service.go` — logika bisnis
- `handlers/product_handler.go` — HTTP handlers

Notes / Tips
- Pastikan `DB_CONN` format sesuai Supabase (gunakan transaction pooler URL bila perlu).
- Untuk deploy di Railway / Railway-like platforms: set env vars `PORT` dan `DB_CONN` di settings deployment.

Jika mau, saya bisa membuat skrip kecil untuk smoke-test endpoint atau menambahkan dokumentasi Postman/k6 untuk load-test.

License
- MIT

