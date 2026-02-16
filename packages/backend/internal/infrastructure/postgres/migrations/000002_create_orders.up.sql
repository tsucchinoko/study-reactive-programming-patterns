CREATE TABLE orders (
    id            UUID PRIMARY KEY,
    customer_id   UUID NOT NULL,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    status        TEXT NOT NULL DEFAULT 'CREATED',
    total_amount  BIGINT NOT NULL,
    currency      TEXT NOT NULL DEFAULT 'JPY',
    placed_at     TIMESTAMPTZ NOT NULL,
    confirmed_at  TIMESTAMPTZ,
    delivered_at  TIMESTAMPTZ,
    cancelled_at  TIMESTAMPTZ,
    cancel_reason TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id             UUID NOT NULL REFERENCES orders(id),
    menu_item_id         UUID NOT NULL,
    name                 TEXT NOT NULL,
    quantity             INTEGER NOT NULL,
    unit_price           BIGINT NOT NULL,
    currency             TEXT NOT NULL DEFAULT 'JPY',
    special_instructions TEXT NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_restaurant ON orders(restaurant_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order ON order_items(order_id);
