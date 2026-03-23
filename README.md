# Scrape Scheduler

Quotes Scraping İşlemlerini Temproal Workflow ile zamanlayan servis.Cron tabanblı tetikleme yapar ve Scraper Gatewaye HTTP isteği gönderir.

## Mimari 
[Client] -- Cron tetiklenir --> [Temporal Workflow] --> [Activity] --> [Scraper Gateway HTTP] | [Worker] <-- Workflow dinnler

POST /scrape/quotes/start/{start}/end/{end}









Kurulacak programlar

Temporal Server -> https://docs.temporal.io/

macos için
brew install temporal
temporal server start-dev