CREATE TABLE restaurants (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    cuisine    TEXT NOT NULL,
    lat        DOUBLE PRECISION NOT NULL,
    lng        DOUBLE PRECISION NOT NULL,
    address    TEXT NOT NULL,
    is_open    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE menu_items (
    id            UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    name          TEXT NOT NULL,
    category      TEXT NOT NULL,
    price         BIGINT NOT NULL,
    currency      TEXT NOT NULL DEFAULT 'JPY',
    prep_time_sec INTEGER NOT NULL,
    available     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_menu_items_restaurant ON menu_items(restaurant_id);
