INSERT INTO users VALUES 
('0002','Ferris',current_timestamp)
ON CONFLICT
ON CONSTRAINT pk_users
DO NOTHING;