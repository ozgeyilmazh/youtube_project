package quotesdiscovery

import "context"

type Repository interface {
	BulkInsertData(ctx context.Context, quotes []Quotes) error
}

// clickhouse kuracağız bunun için
