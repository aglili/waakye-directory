
-- This migration file is used to remove the image_url column from the waakye_vendors table.
ALTER TABLE waakye_vendors
DROP COLUMN image_url;