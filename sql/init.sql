drop table if exists tasks;

create table tasks
(
    id          serial primary key,
    title       varchar(255) not null,
    description text,
    completed   boolean   default false,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
);


insert into tasks (title, description, completed)
values ('Изучить Go', 'Пройти базовый курс', false),
       ('Изучить REST API', 'Посмотреть это видео и написать самому', false),
       ('Зарелизить приложение', 'Нужно спросить у LLM', false)
