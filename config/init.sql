CREATE TABLE device_models (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sensor_models (
    id UUID PRIMARY KEY,
    device_id UUID REFERENCES device_models(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(255) NOT NULL,
    config JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sensor_readings_models (
    id UUID PRIMARY KEY,
    sensor_id UUID REFERENCES sensors(id) ON DELETE CASCADE,
    device_id UUID REFERENCES device_models(id) ON DELETE CASCADE,
    type VARCHAR(255),
    value FLOAT NOT NULL,
    unit VARCHAR(50),
    timestamp TIMESTAMP NOT NULL,
    meta JSONB
);