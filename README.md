
# Scrape Scheduler
Quotes Scraping işlemlerini Temporal workflow ile zamanlayan servis. cron tabanlı tetikleme yapar ve Scraper Gatewaye HTTP isteği gönderir.

##  Mimari 
[Client] -- Cron tetiklenir --> [Temporal Workflow] --> [Activity] --> [Scraper Gateway HTTP] | [Worker] <-- Workflow dinler 

POST /scrape/quotes/start/{start}/end/{end}
