
---

## `README.md`

```markdown
# Fleet Management System (Backend)

Backend service untuk **Sistem Manajemen Armada** Transjakarta.  
Dibangun dengan **Golang**, menggunakan **MQTT, PostgreSQL, RabbitMQ, dan Docker**.

Video Link : https://drive.google.com/file/d/1LUBBdXmauS7cDbZWN5vWN-QOPr95KRKi/view?usp=sharing
---

## ⚙️ Fitur Utama
- **Terima data lokasi kendaraan via MQTT** (`/fleet/vehicle/{vehicle_id}/location`)
- **Simpan data lokasi ke PostgreSQL**
- **REST API**:
  - `GET /vehicles/{vehicle_id}/location` → lokasi terakhir
  - `GET /vehicles/{vehicle_id}/history?start=...&end=...` → riwayat lokasi
- **Geofence detection** → trigger event jika kendaraan masuk radius **50m**
- **Publish event ke RabbitMQ** (`geofence_alerts`)
- **Script mock publisher** → kirim data dummy ke MQTT setiap 2 detik
- **Docker Compose** → semua service (App, PostgreSQL, RabbitMQ, MQTT) jalan otomatis

---

## 📂 Struktur Project
```

.
├── main.go
├── config/
│   └── db.go
├── model/
│   └── vehicle\_location.go
├── repository/
│   └── vehicle\_location\_repository.go
├── service/
│   ├── vehicle\_location\_service.go
│   └── geofence\_service.go
├── handler/
│   └── vehicle\_handler.go
├── mqtt/
│   └── subscriber.go
├── rabbitmq/
│   └── publisher.go
├── scripts/
│   └── publisher.go
├── Dockerfile
├── docker-compose.yml
└── README.md

````

---

##Cara Menjalankan

### 1. Clone repo
```bash
git clone https://github.com/ardifx01/Sistem-Manajemen-Armada.git
cd fleet-management
````

### 2. Jalankan dengan Docker Compose

```bash
docker-compose up --build
```

Service yang akan jalan:

* **App (Go backend)** → [http://localhost:8080](http://localhost:8080)
* **PostgreSQL** → localhost:5432
* **RabbitMQ** (management UI) → [http://localhost:15672](http://localhost:15672) (user: guest / pass: guest)
* **MQTT Broker (Mosquitto)** → localhost:1883

### 3. Jalankan script mock publisher

Di terminal lain:

```bash
go run scripts/publisher.go
```

Script ini akan mengirim data lokasi kendaraan setiap 2 detik ke MQTT.

---

## 📡 API Endpoints

### Lokasi Terakhir Kendaraan

```http
GET /vehicles/{vehicle_id}/location
```

**Response:**

```json
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": "2024-09-07T15:44:16Z"
}
```

### Riwayat Lokasi Kendaraan

```http
GET /vehicles/{vehicle_id}/history?start=1715000000&end=1715009999
```

**Response:**

```json
[
  {
    "vehicle_id": "B1234XYZ",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "timestamp": "2024-09-07T15:44:16Z"
  },
  {
    "vehicle_id": "B1234XYZ",
    "latitude": -6.2090,
    "longitude": 106.8460,
    "timestamp": "2024-09-07T15:46:16Z"
  }
]
```

---

## 🛰️ Geofence Event

Jika kendaraan masuk radius **50 meter** dari titik Monas (`-6.2088, 106.8456`), aplikasi akan publish event ke RabbitMQ:

**Exchange:** `fleet.events`
**Queue:** `geofence_alerts`

**Message:**

```json
{
  "vehicle_id": "B1234XYZ",
  "event": "geofence_entry",
  "location": {
    "latitude": -6.2088,
    "longitude": 106.8456
  },
  "timestamp": 1715003456
}
```

---

## Testing dengan Postman

* Import **Postman Collection** (`postman_collection.json`) yang tersedia di repo
* Jalankan request:

  * `GET /vehicles/{id}/location`
  * `GET /vehicles/{id}/history?start=...&end=...`

---

## Demo

Untuk demo, jalankan:

1. `docker-compose up --build`
2. `go run scripts/publisher.go`
3. Akses API di `http://localhost:8080`
4. Lihat event terkirim di RabbitMQ (`http://localhost:15672`)

---

## Teknologi

* [Go](https://go.dev/) (Fiber, Paho MQTT, pgx, amqp)
* [PostgreSQL](https://www.postgresql.org/)
* [RabbitMQ](https://www.rabbitmq.com/)
* [MQTT Mosquitto](https://mosquitto.org/)
* [Docker](https://www.docker.com/)

```

docker exec -it fleet_rabbitmq bash
rabbitmqctl list_queues
rabbitmqadmin get queue=geofence_alerts requeue=false
