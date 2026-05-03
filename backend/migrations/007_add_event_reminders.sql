-- +goose Up
CREATE TABLE event_reminders (
    id         BIGSERIAL PRIMARY KEY,
    event_id   BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    sent_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT event_reminders_event_id_unique UNIQUE (event_id)
);

-- +goose Down
DROP TABLE IF EXISTS event_reminders;
