-- Enable the uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create locations table
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    street_address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    region VARCHAR(50) NOT NULL,
    longitude DECIMAL(10,8) NOT NULL,
    latitude DECIMAL(11,8) NOT NULL,
    landmark TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes to locations table
CREATE INDEX idx_locations_city ON locations(city);
CREATE INDEX idx_locations_region ON locations(region);
CREATE INDEX idx_locations_coordinates ON locations(longitude, latitude);

-- Create waakye_vendors table
CREATE TABLE waakye_vendors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    description TEXT,
    operating_hours TEXT,
    phone_number VARCHAR(20),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes to waakye_vendors table
CREATE INDEX idx_waakye_vendors_name ON waakye_vendors(name);
CREATE INDEX idx_waakye_vendors_location ON waakye_vendors(location_id);
CREATE INDEX idx_waakye_vendors_phone ON waakye_vendors(phone_number);
