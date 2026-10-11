# Panduan Lanjutan Proyek E-Commerce

**Diperbarui: 3 Oktober 2026 · Versi ringkas 1.1**

Dokumen tunggal untuk melanjutkan pembelajaran di chat baru. Isinya menggabungkan kebutuhan aplikasi, progres, cara menjalankan, dan aturan pendampingan. Status berdasarkan percakapan serta hasil yang dikirim pengguna; folder proyek/GitHub belum diaudit langsung. Gunakan kode aktual sebagai acuan implementasi.

## 1. Cara AI mendampingi

- Pengguna menulis kode, menjalankan perintah, commit, dan push sendiri. AI hanya membaca source yang diberi akses, meninjau, menjelaskan, dan memberikan panduan.
- Jangan otomatis mengedit file proyek, memasang paket, menjalankan aplikasi/Docker, mengubah database, membuat PR, atau deploy. Permintaan membuat dokumen ini bukan izin mengubah aplikasi.
- Baca kode aktual sebelum memberi perubahan. Jika akses belum tersedia, minta potongan yang relevan; jangan menebak struktur.
- Jelaskan dalam bahasa Indonesia sederhana, dengan perbandingan Laravel bila membantu. Pengguna sudah belajar Laravel/PHP/React, tetapi masih mempelajari Go, API, Docker, dan DevOps.
- Berikan tujuan, lokasi perubahan, kode, kegunaan, dan cara menguji dalam langkah yang praktis. Hindari pemeriksaan berulang atau perluasan konfigurasi yang tidak diperlukan.
- Bedakan fitur berhasil, kode yang baru diberikan, dan rencana. Jangan menyatakan pengujian berhasil tanpa hasil.
- Jangan membaca atau menampilkan rahasia `.env`. Gunakan nama variabel atau contoh yang disamarkan.

## 2. Tujuan dan stack

Tujuan: membangun e-commerce latihan yang modular, memahami alur frontend–API–database, lalu mempelajari Git, CI/CD, dan deployment.

| Bagian | Pilihan saat ini |
|---|---|
| Frontend | React + TypeScript, Vite |
| Backend | Go, net/http, encoding/json |
| Database | PostgreSQL dalam Docker |
| Akses database | Raw SQL dengan pgx/pgxpool |
| Migration | golang-migrate |
| GUI database | DBX, untuk melihat/mengelola data |
| Lingkungan | Windows, PowerShell |

GORM pernah dibahas tetapi belum dipilih. Docker saat ini menjalankan database; Go dan frontend berjalan langsung di komputer. Pilihan hosting serta detail CI/CD belum ditetapkan.

## 3. Lokasi dan struktur

- Root proyek: `D:\WORK\ecom`.
- Backend: `D:\WORK\ecom\backend`; DevOps: `D:\WORK\ecom\devops`.
- Module Go: `example.com/ecom/backend`.
- Nama folder frontend dan kapitalisasi filename harus diperiksa dari source; belum dipastikan apakah `fe` atau `frontend`.

Struktur yang digunakan/dipandu, relatif terhadap root di atas:

```text
backend/
  cmd/api/main.go
  internal/
    database/postgres.go
    httpapi/router.go
    health/handler.go
    products/product.go
    products/repository.go
    products/handler.go
  migrations/
  seeds/products.sql
  go.mod, go.sum
<folder-frontend>/
  src/features/products/
    types.ts, api.ts, <halaman-produk>.tsx
  src/features/cart/
  package.json, package-lock.json
devops/
  compose.yaml, .env
```

Ini bukan hasil pemindaian folder. `CreateProductForm.tsx` baru diusulkan; keberadaan dan penerapannya belum dikonfirmasi.

## 4. Progres yang harus dipertahankan

| Bagian | Status terakhir |
|---|---|
| PostgreSQL, tabel products, migration dan data awal | Berhasil menurut pengguna |
| GET produk dari database sampai React | Sudah berhasil |
| Keranjang | Tambah/hapus barang, jumlah dan total sudah digunakan; masih di memori React |
| Repository.Create dan POST produk | Sudah diuji berhasil |
| Penolakan nama berisi spasi | Sudah diuji berhasil |
| Modul API frontend untuk POST | Pengguna menyatakan sudah mengisi perubahan; hasil build belum dikonfirmasi |
| Formulir tambah produk React | BELUM dikonfirmasi selesai; pengguna meminta mengulang tahap ini nanti |
| Update/delete produk, login, checkout, pesanan | Belum dikerjakan |
| CI/CD dan deployment | Belum dikerjakan; tetap target pembelajaran |

Pada pengujian POST, input Sepatu Kanvas, harga 250000, stok 10 menghasilkan produk ID 4. Input nama spasi ditolak dengan pesan nama wajib diisi. Itu bukti dua kasus tersebut, bukan pengujian seluruh validasi.

Halaman produk terakhir yang dikirim memiliki prop `onAddToCart`, state products/loading/error, pemanggilan GET melalui useEffect + AbortController, format rupiah, dan tombol keranjang yang nonaktif saat stok <= 0. Belum ada formulir pada source yang dikirim.

Keranjang berbeda dari katalog: menghapus barang di keranjang tidak menghapus produk dari PostgreSQL. Refresh halaman saat ini menghilangkan keranjang, tetapi produk database tetap tersimpan.

## 5. Alur kode dan kontrak API

```text
Form React -> pemanggil API TypeScript -> router Go
           -> handler -> repository -> pgx -> PostgreSQL
```

- Modul API frontend memakai fetch untuk mengirim request; bukan tempat menjalankan SQL.
- Handler menerima HTTP, memvalidasi input, dan mengirim respons. Perannya mirip controller Laravel.
- Repository menjalankan SQL. Ia menampung pekerjaan penyimpanan yang di Laravel sering dilakukan lewat Eloquent/query builder; repository bukan ORM.
- POST frontend memanggil handler; handler memanggil repository. Satu fitur terbagi tanggung jawab, bukan menyimpan produk tiga kali.

| Endpoint | Perilaku |
|---|---|
| GET /api/health | JSON status/message backend |
| GET /api/products | Array produk dari database, urut id |
| POST /api/products | Simpan produk dan kembalikan objek beserta ID |

```json
{ "name": "Sepatu Kanvas", "price": 250000, "stock": 10 }
```

Kontrak yang dipandu: POST sukses 201 + JSON; input salah 400 + teks; kegagalan internal 500 + pesan umum. Handler memeriksa nama, integer harga/stok, nilai negatif, field hilang/null, JSON dan field tak dikenal; body dibatasi 1 MiB. Semua kasus selain dua pengujian di atas masih perlu diverifikasi dari implementasi.

Tabel products: id BIGINT identity primary key; name TEXT tidak kosong; price BIGINT >= 0; stock INTEGER >= 0. Stok maksimal mengikuti INTEGER PostgreSQL. Model Go memakai ID int, Name string, Price int64, Stock int.

Repository.Create memakai INSERT ... VALUES ($1, $2, $3) RETURNING id, name, price, stock. Input dikirim sebagai parameter, bukan digabung ke teks SQL.

Modul API TypeScript yang dipandu mengekspor `getProducts(signal?)` dan `createProduct(input)`. Gunakan `/api/products` dengan slash awal. POST mengirim JSON; error dibaca sebagai teks sesuai respons backend.

## 6. Menjalankan proyek kembali

Konfigurasi lokal yang disepakati adalah satu `.env` di `D:\WORK\ecom\devops`, berisi POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD. Jangan mengunggah nilainya ke GitHub.

Perbaikan startup terakhir: main.go membaca `../devops/.env` melalui godotenv.Read, menyusun URL dengan url.UserPassword, lalu memanggil database.Open. Alamat database lokal 127.0.0.1:5432; sslmode=disable khusus konfigurasi lokal ini. Godotenv v1.5.1 sudah dipasang.

Usulan `.env` backend tambahan dan internal/config sebelumnya dibatalkan. Jangan memperkenalkannya lagi sebagai syarat melanjutkan fitur.

1. Jalankan PostgreSQL dari Docker Desktop. Jika sudah hidup, langsung lanjut. Alternatif perintah pengguna:

```powershell
cd D:\WORK\ecom\devops
docker compose up -d db
```

2. Di terminal backend:

```powershell
cd D:\WORK\ecom\backend
go run ./cmd/api
```

Jalankan dari folder backend agar path relatif `.env` sesuai. Server menggunakan http://localhost:8080.

3. Di terminal lain, masuk folder frontend aktual yang berisi package.json, lalu:

```powershell
npm run dev
```

Buka URL Vite yang ditampilkan. Proxy pengembangan yang dipandu meneruskan /api ke 127.0.0.1:8080.

Tidak perlu mengulang instalasi, migration atau seed setiap membuka proyek. Seed awal bisa menggandakan data jika diulang. Pembacaan `.env` dalam Go tidak otomatis mengisi environment terminal untuk CLI migration.

## 7. GitHub: status terbaru dan masalah terakhir

- Git terpasang: 2.44.0.windows.1.
- Repository lokal sekarang sudah berada pada branch main; laporan terakhir menunjukkan working tree clean dan tidak ada perubahan baru untuk commit.
- Remote tujuan: https://github.com/dyarodyh/ecom.git.
- Push terakhir DITOLAK dengan 403: Permission to dyarodyh/ecom.git denied to imarod.
- Pengguna ingin memakai akun dyarodyh. Perbaikan berikut sudah diberikan, tetapi hasilnya belum dikonfirmasi:

```powershell
cd D:\WORK\ecom
git remote set-url origin https://dyarodyh@github.com/dyarodyh/ecom.git
git push -u origin main
```

Pengguna menyelesaikan login sebagai dyarodyh jika diminta. Jangan menjalankan push atas nama pengguna. Jangan menghapus semua kredensial atau mengubah konfigurasi global tanpa kebutuhan yang jelas.

`nothing to commit` adalah status lokal, bukan bukti upload berhasil. `origin already exists` berarti remote sudah terdaftar; ubah dengan set-url, bukan menambah origin lagi. Commit hash dan isi repository remote belum diperiksa.

Aturan ignore yang telah dipandu mencakup `.env`, `.env.*` dengan pengecualian `.env.example`, node_modules, dist, coverage, log, dan hasil build Go. Penerapan aktual harus dicocokkan dengan repository; contoh environment tidak boleh berisi rahasia.

AI boleh membaca folder lokal yang diberikan akses untuk melihat perubahan terbaru tanpa menunggu push. Jika sumbernya GitHub, AI hanya melihat versi yang sudah diunggah; repository privat memerlukan akses yang sesuai.

## 8. Langkah berikutnya dan target selesai

1. Gunakan folder lokal atau repository yang sudah dapat diakses sebagai sumber review; periksa nama frontend, file API, halaman produk, dan form jika ada.
2. Cocokkan status dengan source. Jangan mengulang backend POST yang sudah berhasil tanpa alasan.
3. Ulangi panduan formulir React dari bagian yang benar-benar belum ditulis.
4. Form memanggil createProduct; setelah sukses, berikan Product hasil backend kepada halaman melalui callback onCreated.
5. Halaman menambahkan produk ke state menggunakan pembaruan array; pertahankan integrasi onAddToCart.

**Formulir selesai bila:** produk dapat dibuat dari browser, pesan menyimpan/berhasil/error tersedia, input salah ditolak, produk tetap muncul setelah refresh, produk baru bisa masuk keranjang, dan build TypeScript berhasil menurut hasil pengguna.

Setelah itu, lanjut bertahap: update/delete produk -> CI -> autentikasi/hak admin -> checkout/pesanan -> deployment dan CD. Urutan bisa disesuaikan. Detail login, persistensi keranjang, hosting, dan pembayaran belum dipilih; jangan memasang teknologi tambahan sekaligus.

Masalah lama yang sudah ditangani: migration up/down tertukar, typo INDENTIFY menjadi IDENTITY, kode startup di luar fungsi run, self-import package products, dan import source/stub yang tidak diperlukan. Jangan mengulangi reset migration/database berdasarkan riwayat lama.

## 9. Prompt untuk chat baru

```text
Saya sedang belajar e-commerce React + TypeScript, Go, dan PostgreSQL.
Baca dokumen ECOM_SRS_HANDOFF.md ini dan source proyek yang saya berikan.
Kamu hanya tutor dan peninjau: baca kode, jelaskan, dan berikan panduan.
Jangan otomatis mengedit file, menjalankan proyek, mengubah database,
memasang paket, melakukan commit/push/PR, atau deploy. Saya mengerjakannya sendiri.

GET/POST produk sudah berhasil diuji dan keranjang dasar sudah ada.
Modul API frontend sudah diisi; formulir tambah produk belum dikonfirmasi selesai.
Push GitHub terakhir gagal karena memakai imarod; akun yang dituju dyarodyh.
Periksa source aktual, jelaskan titik lanjutnya, lalu pandu secara praktis.
Bedakan kode yang sudah diterapkan, hasil pengujian, dan rencana.

Sumber: [folder lokal atau URL repository dan branch/commit]
```

Simpan dokumen ini sebagai `D:\WORK\ecom\docs\ECOM_SRS_HANDOFF.md` bila ingin memasukkannya ke proyek. Dokumen ini menjadi acuan ringkas tunggal untuk menggantikan catatan lanjutan lama. Pengguna menyalin dan melakukan commit/push sendiri.

Sesudah tahap berhasil, perbarui tanggal, status fitur, hasil uji, branch/commit bila diketahui, dan langkah berikutnya. Selalu dahulukan source aktual untuk menilai keadaan implementasi.
