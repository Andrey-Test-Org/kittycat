CREATE TABLE IF NOT EXISTS kitty_event (
    id           BIGSERIAL PRIMARY KEY,
    cat_id       TEXT NOT NULL,
    kind         TEXT NOT NULL,
    weight_g     INTEGER NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_kitty_event_cat_time
    ON kitty_event (cat_id, occurred_at DESC);

INSERT INTO kitty_event (cat_id, kind, weight_g, occurred_at) VALUES
    ('luna',  'nap',  4200, now() - interval '2 hours'),
    ('luna',  'play', 4180, now() - interval '90 minutes'),
    ('mochi', 'nap',  3300, now() - interval '4 hours'),
    ('mochi', 'snack',3320, now() - interval '30 minutes'),
    ('pixel', 'play', 5100, now() - interval '10 minutes');
