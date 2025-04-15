


-- to add image_url column to waakye_vendors table
ALTER TABLE waakye_vendors
ADD COLUMN image_url VARCHAR(255) DEFAULT 'https://example.com/images/default-waakye-vendor.jpg';