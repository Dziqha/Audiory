
# 🎵 Audiory

Audiory adalah sebuah proyek microservice yang dirancang untuk menyediakan layanan streaming musik dengan arsitektur yang scalable dan efisien. Proyek ini dibangun menggunakan berbagai teknologi modern untuk memastikan performa yang optimal dan kemudahan pengembangan.

## ✨ Fitur Utama
- 🏗️ **Microservices Architecture**: Audiory dibangun dengan pendekatan microservices, memungkinkan setiap bagian sistem untuk dikembangkan, diuji, dan dikelola secara terpisah.
- 🚀 **gRPC**: Komunikasi antar microservices dilakukan menggunakan gRPC, yang menawarkan performa tinggi dan efisiensi dalam pengiriman data.
- ⚡ **Caching dengan Redis**: Untuk meningkatkan performa, Audiory menggunakan Redis sebagai caching layer.
- 🔄 **Pub/Sub dengan Redis**: Sistem rekomendasi diimplementasikan menggunakan mekanisme publish/subscribe dari Redis, memungkinkan komunikasi yang efisien antar layanan.

## 💳 Sistem Pembayaran
Audiory menyediakan sistem pembayaran untuk layanan berlangganan (subscription) yang dikelola menggunakan Midtrans. Berikut adalah rincian tentang sistem pembayaran:

### Fitur Pembayaran
- **Berlangganan**: Pengguna dapat memilih paket berlangganan yang tersedia untuk mengakses layanan streaming musik.
- **Integrasi Midtrans**: Pembayaran diproses melalui Midtrans, yang menyediakan solusi pembayaran yang aman dan terpercaya.

### Konfigurasi Midtrans
Untuk mengintegrasikan Midtrans, Anda perlu:
1. Mendaftar di [Midtrans](https://midtrans.com) dan mendapatkan API key.
2. Mengonfigurasi API key di file konfigurasi layanan pembayaran.
3. Memastikan semua endpoint yang diperlukan untuk komunikasi dengan Midtrans telah diimplementasikan.

## 🗄️ Database
- 🧑‍💻 **User Service**: Menggunakan PostgreSQL untuk menyimpan data pengguna.
- 🎶 **Music Service**: Menggunakan MySQL untuk mengelola data musik.

## 🔧 Framework
- 🌀 **User Service**: Dibangun dengan menggunakan Fiber, framework Go yang ringan dan cepat.
- 🌐 **Music Service**: Dibangun menggunakan Express.js, framework Node.js yang populer.

## 📚 ORM
- 🌀 **User Service**: Menggunakan GORM, ORM yang kuat untuk Go.
- 🌐 **Music Service**: Menggunakan Prisma, ORM modern untuk Node.js.

## 🏛️ Arsitektur
Audiory terdiri dari beberapa komponen utama:
1. 🧑‍💻 **User Service**: Mengelola data pengguna, termasuk registrasi, dan autentikasi pengguna.
2. 🎶 **Music Service**: Mengelola data musik, termasuk informasi lagu, album, dan playlist.
3. 🔮 **Recommendation Service**: Memberikan rekomendasi musik berdasarkan preferensi pengguna menggunakan sistem pub/sub Redis.
4. 🔒 **Middleware**: Autentikasi dan keamanan sebelum permintaan API mencapai endpoint-nya.

## 🛠️ Teknologi yang Digunakan
- 🖥️ **Bahasa Pemrograman**: Go, JavaScript (Node.js)
- 🔧 **Framework**:
  - Go: Fiber
  - Node.js: Express.js
- 🗄️ **Database**:
  - PostgreSQL (User Service)
  - MySQL (Music Service)
- 📚 **ORM**:
  - GORM (Go)
  - Prisma (Node.js)
- ⚡ **Caching**: Redis
- 🚀 **Komunikasi**: gRPC

## ⚙️ Instalasi
Untuk menjalankan Audiory secara lokal, ikuti langkah-langkah berikut:

### 📥 Clone Repository:
```bash
git clone https://github.com/Dziqha/audiory.git
cd audiory
```

### 📦 Instalasi Dependencies:
#### Untuk User Service:
```bash
cd user-service
go mod tidy
```

#### Untuk Music Service:
```bash
cd music-service
npm install
```

### ⚙️ Konfigurasi Database:
1. Pastikan PostgreSQL dan MySQL telah terinstal dan berjalan.
2. Buat database untuk User Service dan Music Service.
3. Sesuaikan konfigurasi database pada file konfigurasi masing-masing service.

### ▶️ Menjalankan Layanan:
#### Untuk User Service menggunakan air:
```bash
cd user-service
air
```

#### Untuk Music Service menggunakan nodemon:
```bash
cd music-service
npm run dev
```

### 🔄 Menjalankan Redis:
Pastikan Redis telah terinstal dan berjalan.

## 📖 Penggunaan
Setelah semua layanan berjalan, Anda dapat mengakses API yang disediakan oleh masing-masing service. Pastikan untuk merujuk pada dokumentasi API untuk informasi lebih lanjut tentang endpoint yang tersedia.

## 🤝 Kontribusi
Kami sangat terbuka untuk kontribusi dari pengembang lain. Jika Anda ingin berkontribusi, silakan buka isu atau kirim pull request.

## 📜 Lisensi
Proyek ini dilisensikan di bawah MIT License. Silakan lihat file LICENSE untuk informasi lebih lanjut.
