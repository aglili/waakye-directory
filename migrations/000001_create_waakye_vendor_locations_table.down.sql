-- Drop indexes
DROP INDEX IF EXISTS idx_locations_city;
DROP INDEX IF EXISTS idx_locations_region;
DROP INDEX IF EXISTS idx_locations_coordinates;
DROP INDEX IF EXISTS idx_waakye_vendors_name;
DROP INDEX IF EXISTS idx_waakye_vendors_location;
DROP INDEX IF EXISTS idx_waakye_vendors_phone;

-- Drop tables
DROP TABLE IF EXISTS waakye_vendors;
DROP TABLE IF EXISTS locations;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
