ALTER TABLE parcels
ADD COLUMN sender_description TEXT NULL,
ADD COLUMN motorbike_description TEXT NULL,
ALTER COLUMN pickup_address SET NOT NULL,
ALTER COLUMN dropoff_address SET NOT NULL,
ALTER COLUMN latitude SET NOT NULL,
ALTER COLUMN longitude SET NOT NULL;
