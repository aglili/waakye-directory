-- Create ratings table
CREATE TABLE vendor_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES waakye_vendors(id) ON DELETE CASCADE,
    hygiene_rating INTEGER NOT NULL CHECK (hygiene_rating BETWEEN 1 AND 5),
    value_rating INTEGER NOT NULL CHECK (value_rating BETWEEN 1 AND 5),
    service_rating INTEGER NOT NULL CHECK (service_rating BETWEEN 1 AND 5),
    taste_rating INTEGER NOT NULL CHECK (taste_rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_ratings_vendor_id ON vendor_ratings(vendor_id);