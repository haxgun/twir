-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

UPDATE "channels_overlays_layers" SET settings_html_html = convert_from(decode(settings_html_html, 'base64'), 'UTF-8') WHERE settings_html_html IS NOT NULL;
UPDATE "channels_overlays_layers" SET settings_html_css = convert_from(decode(settings_html_css, 'base64'), 'UTF-8') WHERE settings_html_css IS NOT NULL;
UPDATE "channels_overlays_layers" SET settings_html_js = convert_from(decode(settings_html_js, 'base64'), 'UTF-8') WHERE settings_html_js IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
