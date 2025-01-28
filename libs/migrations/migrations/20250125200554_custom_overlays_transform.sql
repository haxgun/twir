-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_overlays_layers ADD COLUMN transform_string text NOT NULL DEFAULT 'translate(0px, 0px)';

UPDATE channels_overlays_layers SET transform_string = format('translate(%spx,%spx)', pos_x, pos_y);

ALTER TABLE channels_overlays_layers DROP COLUMN pos_x;
ALTER TABLE channels_overlays_layers DROP COLUMN pos_y;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
