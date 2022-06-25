CREATE TABLE todo (
    id bigserial NOT NULL,
    title varchar(255) NOT NULL,
    note text,
    finished boolean NOT NULL,
    due_date timestamp,
    CONSTRAINT todo_pkc PRIMARY KEY(id)
);