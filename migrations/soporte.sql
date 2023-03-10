-- DROP TABLE IF EXISTS users;
-- DROP TABLE IF EXISTS invitaciones;
-- DROP TABLE IF EXISTS clientes;
-- DROP TABLE IF EXISTS areas;
-- DROP TABLE IF EXISTS funcionarios;
-- DROP TABLE IF EXISTS             



CREATE TABLE users (
    user_id uuid DEFAULT uuid_generate_v4 (),
    username VARCHAR ( 50 ) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    last_login TIMESTAMP,
    created_on TIMESTAMP NOT NULL,
    email VARCHAR ( 50 ) UNIQUE NOT NULL,
    estado INT DEFAULT 0,
    PRIMARY KEY (user_id)
);

CREATE TABLE if not exists invitaciones(
    id uuid DEFAULT uuid_generate_v4 (),
    email VARCHAR UNIQUE NOT NULL,
    pendiente BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    creador_id VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE if not exists casos  (
    id uuid DEFAULT uuid_generate_v4 (),
    titulo VARCHAR UNIQUE NOT NULL,
    descripcion VARCHAR NOT NULL,
    detalles_de_finalizacion TEXT,
    empresa INT NOT NULL,
    area INT NOT NULL,
    cliente_name VARCHAR NOT NULL,
    funcionario_name VARCHAR,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP,
    fecha_inicio TIMESTAMP,
    fecha_fin TIMESTAMP,
    prioridad INT,
    estado INT DEFAULT 0,
    client_id VARCHAR NOT NULL,
    funcionario_id VARCHAR,
    superior_id VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE if not exists empresas (
    id serial primary key,
    slug VARCHAR UNIQUE NOT NULL,
    nombre VARCHAR UNIQUE NOT NULL,
    telefono VARCHAR UNIQUE NOT NULL,
    estado INT DEFAULT 0,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP
);

create table if not exists areas (
    id serial primary key,
    nombre VARCHAR  NOT NULL,
    -- descripcion TEXT,
    estado INT DEFAULT 0,
    empresa_id INT NOT NULL,
    -- empresa_name INT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    creador_id VARCHAR NOT NULL
);

create table if not exists clientes (
    client_id uuid DEFAULT uuid_generate_v4(),
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
    areas INT[],
    estado INT DEFAULT 0,
    profile_photo TEXT,
    is_admin BOOLEAN DEFAULT false,
    rol INT,
    PRIMARY KEY (client_id)
);


create table if not exists funcionarios (
    funcionario_id uuid DEFAULT uuid_generate_v4(),
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
    areas INT[],
    profile_photo TEXT,
    rol INT,
    PRIMARY KEY (funcionario_id)
);

create table if not exists user_area (
    user_id VARCHAR NOT NULL,
    area_id int NOT NULL,
    estado int DEFAULT 0,
    nombre_user VARCHAR NOT NULL,
    nombre_area VARCHAR NOT NULL
);

-- create table if not exists user_area (
--     client_id VARCHAR references clientes (client_id) on update cascade,
--     area_id int references areas (id) on update cascade,
--     estado int DEFAULT 0,
--     nombre VARCHAR NOT NULL,
-- );



create table if not exists recursos (
    id serial primary key,
    file_url TEXT NOT NULL,
    ext VARCHAR,
    descripcion VARCHAR,
    caso_id VARCHAR NOT NULL,
    created_on TIMESTAMP NOT NULL
);


create table if not exists conversations (
    -- id serial primary key,
    caso_id VARCHAR primary key,
    client_id VARCHAR NOT NULL,
    funcionario_id VARCHAR,
);


create table if not exists messages (
    id serial primary key,
    caso_id uuid,
    client_id uuid NOT NULL,
    funcionario_id uuid,
    client_name VARCHAR,
    funcionario_name VARCHAR,
    media_url VARCHAR,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_on TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT false
);

-- UPDATE clientes SET areas = areas || '{1}' WHERE nombre = 'alejandro' ;

-- insert into users (user_id,email,username,password,created_on) values (
--     '8ba62445-3c62-45ad-aab5-e15bc5d68337','alejandro@gmail.com','alejandro',crypt('12ab34cd56ef',gen_salt('bf')),now());

-- insert into clientes (email,nombre,rol,empresa_id,is_admin,user_id,superior_id,created_on) 
-- values ('alejandro@gmail.com','alejandro',2,5,true,'8ba62445-3c62-45ad-aab5-e15bc5d68337',
-- '8ba62445-3c62-45ad-aab5-e15bc5d68337',now());


insert into users (user_id,email,username,password,created_on) values (
    '8ba62445-3c62-45ad-aab5-e15bc5d68323','jorgemiranda0180@gmail.com','jorgedm',crypt('12ab34cd56ef',gen_salt('bf')),now());
insert into funcionarios(email,nombre,rol,empresa_id,user_id,superior_id,created_on) values('jorgemiranda0180@gmail.com','jorgedm',
3,6,'8ba62445-3c62-45ad-aab5-e15bc5d68323','8ba62445-3c62-45ad-aab5-e15bc5d68323',now());