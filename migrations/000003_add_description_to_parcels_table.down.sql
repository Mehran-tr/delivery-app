ALTER TABLE parcels
DROP COLUMN sender_description,
DROP COLUMN motorbike_description,
ALTER COLUMN pickup_address DROP NOT NULL,
ALTER COLUMN dropoff_address DROP NOT NULL,
ALTER COLUMN latitude DROP NOT NULL,
ALTER COLUMN longitude DROP NOT NULL;
