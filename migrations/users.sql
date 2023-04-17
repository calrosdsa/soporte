
create table if not exists clientes (
    client_id uuid DEFAULT gen_random_uuid(),
    nombre VARCHAR,
    apellido VARCHAR,
    celular VARCHAR,
    email VARCHAR ( 50 ),
    superior_id VARCHAR,
    empresa_id INT NOT NULL,
    telefono VARCHAR,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP,
    user_id VARCHAR NOT NULL,
    estado INT DEFAULT 0,
    profile_photo TEXT,
    is_admin BOOLEAN DEFAULT false,
    rol INT,
    PRIMARY KEY (client_id)
);


create table if not exists funcionarios (
    funcionario_id uuid DEFAULT gen_random_uuid(),
    nombre VARCHAR,
    apellido VARCHAR,
    celular VARCHAR,
    email VARCHAR ( 50 ) UNIQUE NOT NULL,
    superior_id VARCHAR,
    empresa_id INT NOT NULL,
    telefono VARCHAR,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP,
    user_id VARCHAR NOT NULL,
    estado INT DEFAULT 0,
    profile_photo TEXT,
    is_admin BOOLEAN DEFAULT false,
    rol INT
);


CREATE TABLE users (
    user_id uuid DEFAULT uuid_generate_v4 (),
    username VARCHAR ( 50 ) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP,
    email VARCHAR ( 50 ) UNIQUE NOT NULL,
    estado INT DEFAULT 0,
    PRIMARY KEY (user_id)
);

CREATE TABLE if not exists invitaciones(
    id uuid DEFAULT uuid_generate_v4 (),
    email VARCHAR UNIQUE NOT NULL,
    pendiente BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    creador_id uuid,
    send_on TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);



-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -- Create the table with a UUID primary key
-- CREATE TABLE mytable (
--   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   column1 TEXT,
--   column2 INTEGER,
--   column3 TIMESTAMP
-- );

-- -- Insert 1 million rows
-- INSERT INTO mytable (column1, column2, column3)
-- SELECT 
--     md5(random()::text), -- generate a random string for column1
--     random() * 100, -- generate a random number between 0 and 100 for column2
--     current_timestamp - random() * interval '365 days' -- generate a random date within the last year for column3
-- FROM 
--     generate_series(1, 1000000);



SELECT Count(*)
FROM mytable
-- ORDER BY some_column -- specify an ordering column to make the result predictable
OFFSET 9 -- get the 10th row (offset is 0-indexed)
LIMIT 1; -- limit the result to 1 row


select * from mytable where id = '9e781c87-161b-45fa-933d-a881fe4ded75';
