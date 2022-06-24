INSERT INTO users VALUES 
('0001','gopher1',current_timestamp)
ON CONFLICT
ON CONSTRAINT pk_users
DO NOTHING;