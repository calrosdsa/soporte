CREATE TABLE if not exists empresas (
    id serial primary key,
    slug VARCHAR UNIQUE NOT NULL,
    nombre VARCHAR UNIQUE NOT NULL,
    telefono VARCHAR UNIQUE NOT NULL    ,
    estado INT DEFAULT 0,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP,
    parent_id INT
);


create table if not exists areas (
    id serial primary key,
    nombre VARCHAR  NOT NULL,
    estado INT DEFAULT 0,
    empresa_id INT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    creador_id uuid NOT NULL
);

create table if not exists proyectos (
    id serial primary key,
    nombre VARCHAR  NOT NULL,
    parent_id int NOT NULL,
    estado INT DEFAULT 0,
    empresa_id INT NOT NULL,
    empresa_parent_id INT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    creador_id uuid NOT NULL
);

create table if not exists proyecto_duration (
    id serial primary key,
    proyecto_id int NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);

create table if not exists user_area (
    id serial primary key,
    user_id uuid NOT NULL,
    area_id int NOT NULL,
    estado int DEFAULT 0,
    created_on TIMESTAMP NOT NULL,  
    nombre_area VARCHAR NOT NULL
);