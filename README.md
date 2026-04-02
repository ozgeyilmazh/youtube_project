
# Scrape Scheduler

Quotes scraping işlemlerini Temporal workflow ile zamanlayan servis. Cron tabanlı tetikleme yapar ve Scraper Gateway'e HTTP isteği gönderir.

## Mimari

[Client] -- Cron tetiklenir --> [Temporal Workflow] --> [Activity] --> [Scraper Gateway HTTP] | [Worker] <-- Workflow dinler -- [Temporal] v POST /scrape/quotes/start/{start}/end/{end}


## Bileşenler

| Bileşen | Açıklama |
|---------|----------|
| **client** | Cron schedule ile workflow'u başlatır, sürekli çalışır |
| **worker** | Temporal task queue dinler, workflow ve activity'leri çalıştırır |
| **Workflow** | `ScheduleQuotesScrape` – Activity çağrısını koordine eder |
| **Activity** | `BeginScrape` – Scraper Gateway'e HTTP request gönderir |

## Ön Koşullar

- Go 1.25+
- [Temporal Server](https://docs.temporal.io/self-hosted-guide) (local: `temporal server start-dev`)

## Ortam Değişkenleri

`.env.local` / `.env.dev` / `.env` dosyasında:

| Değişken | Açıklama | Örnek |
|----------|----------|-------|
| `SCHEDULER_CRON` | Cron expression (saatlik vb.) | `0 * * * *` |
| `SCRAPER_GATEWAY_HOST` | Scraper Gateway base URL | `http://localhost:3000` |
| `TEMPORAL_HOST_PORT` | Temporal server adresi | `localhost:7233` |
| `SCHEDULER_TASK_QUEUE` | Temporal task queue adı | `quote-discovery-scheduler-task-queue` |

## Çalıştırma

```bash
# Kurmak için : 
brew install temporal
# Temporal başlat (ayrı terminalde)
temporal server start-dev

# Worker (workflow ve activity dinler)
just run-worker
# veya: ENV=local go run worker/main.go

# Client (cron ile workflow tetikler)
just run-client
# veya: ENV=local go run client/main.go
```

Not: Hem worker hem client aynı anda çalışmalıdır. Client workflow başlatır, worker da çalıştırır.

Health Endpoints
Endpoint	Port	Açıklama
/healthz	8080 (worker), 8081 (client)	Health check
/livez	8080 (worker)	Liveness probe
/readyz	8080 (worker)	Readiness probe

Test
just test local
# veya: ENV=local go test ./... -v

Not: workflow_test.go gerçek Temporal server gerektirir. Test öncesinde temporal server start-dev çalışıyor olmalıdır.

Proje Yapısı

scrape-scheduler/
├── client/main.go        # Cron client – workflow tetikleyici
├── worker/main.go        # Temporal worker
├── workflow.go           # ScheduleQuotesScrape workflow
├── activity.go           # BeginScrape activity
├── scraper_gateway_client.go
├── config.go
├── healthz.go
└── justfile              # test, run-worker, run-client, lint


just lint
# veya: golangci-lint run


# Kontrol
which mockgen

# Yoksa, shell config'e ekle (~/.zshrc veya ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

go generate ./...

https://docs.pact.io/implementation_guides/go/readme

brew install just