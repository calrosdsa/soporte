DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS invitaciones;
DROP TABLE IF EXISTS clientes;
DROP TABLE IF EXISTS areas;
DROP TABLE IF EXISTS funcionarios;
DROP TABLE IF EXISTS sub_areas;
DROP TABLE IF EXISTS user_area;           
DROP TABLE IF EXISTS empresas;    
DROP TABLE IF EXISTS messages;    
DROP TABLE IF EXISTS casos;   
DROP TABLE IF EXISTS recursos;
DROP TABLE IF EXISTS conversations;
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE if not exists casos  (
    id uuid DEFAULT gen_random_uuid(),
    titulo VARCHAR UNIQUE NOT NULL,
    descripcion VARCHAR NOT NULL,
    detalles_de_finalizacion TEXT,
    empresa INT NOT NULL,
    area INT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP,
    fecha_inicio TIMESTAMP,
    fecha_fin TIMESTAMP,
    prioridad INT,
    estado INT DEFAULT 0,
    client_id uuid NOT NULL,
    funcionario_id uuid,
    superior_id uuid,
    status INT DEFAULT 0,
    rol INT,
    PRIMARY KEY (id)
);




create table if not exists recursos (
    id serial primary key,
    file_url TEXT NOT NULL,
    ext VARCHAR,
    descripcion VARCHAR,
    caso_id uuid,
    created_on TIMESTAMP NOT NULL
);


create table if not exists conversations (
    caso_id uuid,
    client_id uuid,
    funcionario_id uuid
);


create table if not exists messages (
    id serial primary key,
    caso_id uuid,
    from_user uuid,
    to_user uuid,
    media_url  text[],
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_on TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT false
);



-- insert into empresas (id,slug,nombre,telefono,created_on) values (6,'teclu','Teclu','75390560',now());
-- insert into users (user_id,email,username,password,created_on) values ('8ba62445-3c62-45ad-aab5-e15bc5d68323','jorgemiranda0180@gmail.com','Jorge',crypt('12ab34cd56ef',gen_salt('bf')),now());
-- insert into funcionarios(email,nombre,apellido,rol,empresa_id,user_id,superior_id,created_on) values('jorgemiranda0180@gmail.com','Jorge','Miranda',3,6,'8ba62445-3c62-45ad-aab5-e15bc5d68323','8ba62445-3c62-45ad-aab5-e15bc5d68323',now());


insert into users (user_id,email,username,password,created_on) values ('8ba62445-3c62-45ad-aab5-e15bc5d68321','marca@yopmail','henry_marca',crypt('201120',gen_salt('bf')),now());
insert into funcionarios(email,nombre,apellido,rol,empresa_id,user_id,superior_id,created_on) values('marca@yopmail','Henry','Marca',3,6,'8ba62445-3c62-45ad-aab5-e15bc5d68321','8ba62445-3c62-45ad-aab5-e15bc5d68321',now());

-- select nombre,apellido,titulo from clientes inner join casos on clientes.client_id = casos.client_id;

-- select client_id,nombre,apellido,profile_photo,user_area.estado from clientes inner join user_area on clientes.client_id = user_area.user_id;


select titulo,clientes.nombre,funcionarios.nombre from casos inner join clientes on clientes.client_id = casos.client_id left join funcionarios on funcionarios.funcionario_id = casos.funcionario_id;