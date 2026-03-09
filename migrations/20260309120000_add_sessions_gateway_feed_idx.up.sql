CREATE INDEX CONCURRENTLY IF NOT EXISTS sessions_gateway_feed_idx
    ON sessions (gateway_feed)
    WHERE removed_at IS NULL;
