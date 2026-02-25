CREATE TABLE drivers (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'AVAILABLE',
    latitude    DOUBLE PRECISION NOT NULL DEFAULT 35.6762,
    longitude   DOUBLE PRECISION NOT NULL DEFAULT 139.6503,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE delivery_assignments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id        UUID NOT NULL REFERENCES orders(id),
    driver_id       UUID NOT NULL REFERENCES drivers(id),
    status          TEXT NOT NULL DEFAULT 'ASSIGNED',
    assigned_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    picked_up_at    TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_assignments_order ON delivery_assignments(order_id);
CREATE INDEX idx_assignments_driver ON delivery_assignments(driver_id);
CREATE INDEX idx_drivers_status ON drivers(status);
