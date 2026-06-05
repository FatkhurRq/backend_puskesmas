# Puskesmas Queue Management System - Backend

Backend API untuk sistem antrean puskesmas menggunakan Go, Gin, PostgreSQL, JWT Authentication, dan GORM.

## Tech Stack

* Golang
* Gin Framework
* PostgreSQL
* GORM
* JWT Authentication

---

# Setup Project

## Clone Repository

```bash
git clone <repository-url>
cd puskesmas-backend
```

## Install Dependencies

```bash
go mod tidy
```

## Create Environment File

Buat file `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=puskesmas_queue

JWT_SECRET=your_secret_key
```

## Run Application

```bash
go run main.go
```

Server akan berjalan di:

```text
http://localhost:8080
```

---

# Authentication

Backend menggunakan JWT.

Login akan mengembalikan token:

```json
{
  "token": "jwt_token_here"
}
```

Untuk endpoint yang membutuhkan login:

```http
Authorization: Bearer <token>
```

---

# User Roles

## Patient

Dapat:

* Melihat poli
* Mendaftar antrean
* Melihat antrean sendiri
* Melihat riwayat pendaftaran

## Admin

Dapat:

* Mengakses dashboard
* Melihat seluruh antrean
* Melihat seluruh pasien
* Memanggil antrean
* Menyelesaikan antrean

---

# API Endpoints

## Authentication

### Register

```http
POST /api/register
```

Body:

```json
{
  "nik": "3578123456789012",
  "name": "Fath",
  "phone": "08123456789",
  "password": "123456"
}
```

---

### Login

```http
POST /api/login
```

Body:

```json
{
  "nik": "3578123456789012",
  "password": "123456"
}
```

---

### Get Profile

```http
GET /api/me
```

Header:

```http
Authorization: Bearer TOKEN
```

---

# Poli

## Get All Polis

```http
GET /api/polis
```

Response:

```json
[
  {
    "id": "...",
    "name": "Poli Umum"
  }
]
```

---

## Create Poli (Admin)

```http
POST /api/polis
```

Body:

```json
{
  "name": "Poli Gigi"
}
```

---

# Registration

## Create Registration

```http
POST /api/registrations
```

Header:

```http
Authorization: Bearer TOKEN
```

Body:

```json
{
  "poli_id": "uuid-poli"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "queue_number": 4
  }
}
```

---

## Get My Registrations

```http
GET /api/registrations/me
```

---

## Get My Queue

```http
GET /api/queue/me
```

Response:

```json
{
  "success": true,
  "data": {
    "queue_number": 4,
    "patient_name": "Fath",
    "poli": "Poli Umum",
    "status": "waiting",
    "registration_date": "2026-06-05T10:00:00Z"
  }
}
```

---

# Queue Monitoring

## Current Queue By Poli

```http
GET /api/polis/:id/queues
```

Response:

```json
{
  "success": true,
  "data": {
    "queue_number": 4,
    "patient_name": "Fath",
    "poli": "Poli Umum",
    "status": "called"
  }
}
```

Frontend dapat melakukan polling setiap 5-10 detik untuk memperbarui antrean secara otomatis.

---

# Admin Endpoints

## Dashboard

```http
GET /api/admin/dashboard
```

Response:

```json
{
  "success": true,
  "data": {
    "total_users": 4,
    "total_poli": 2,
    "total_waiting": 5,
    "total_called": 1,
    "total_completed": 3
  }
}
```

---

## Get All Patients

```http
GET /api/admin/patients
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": "...",
      "nik": "3578123456789012",
      "name": "Fath",
      "phone": "08123456789"
    }
  ]
}
```

---

## Get All Queues

```http
GET /api/admin/queues
```

Response berisi:

* Queue Number
* Queue Status
* Patient Data
* Poli Data

---

## Call Queue

```http
PUT /api/admin/queues/:id/call
```

Response:

```json
{
  "status": "called"
}
```

---

## Complete Queue

```http
PUT /api/admin/queues/:id/complete
```

Response:

```json
{
  "status": "completed"
}
```

---

# Suggested Frontend Pages

## Patient

### Login Page

* Login menggunakan NIK dan password

### Dashboard Page

Menampilkan:

* Nama pasien
* Nomor antrean
* Poli
* Status antrean

Menggunakan:

```http
GET /api/queue/me
```

---

### Registration Page

Menampilkan daftar poli:

```http
GET /api/polis
```

Daftar antrean:

```http
POST /api/registrations
```

---

## Admin

### Dashboard

Menggunakan:

```http
GET /api/admin/dashboard
```

Menampilkan statistik.

---

### Queue Management

Menggunakan:

```http
GET /api/admin/queues
```

Tombol:

* Call
* Complete

---

### Patients Page

Menggunakan:

```http
GET /api/admin/patients
```

---

# Notes For Frontend Developer

* Simpan JWT token di localStorage.
* Sertakan token pada semua request yang membutuhkan autentikasi.
* Gunakan Axios atau Fetch API.
* Polling antrean setiap 5-10 detik pada halaman monitoring.
* Backend menggunakan format response:

```json
{
  "success": true,
  "message": "Success message",
  "data": {}
}
```

Kecuali beberapa endpoint lama yang belum dimigrasikan ke format DTO.
