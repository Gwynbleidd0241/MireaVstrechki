-- +goose Up
ALTER TABLE events
    ADD COLUMN status      VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    ADD COLUMN location    TEXT        NOT NULL DEFAULT '',
    ADD COLUMN meeting_url TEXT        NOT NULL DEFAULT '';

ALTER TABLE event_participants
    ADD COLUMN rsvp_status VARCHAR(20) NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE events
    DROP COLUMN status,
    DROP COLUMN location,
    DROP COLUMN meeting_url;

ALTER TABLE event_participants
    DROP COLUMN rsvp_status;
