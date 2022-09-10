/* Проект "Открытая библиотека"
    В данном проекте пользователи могут скачивать и загружать книги. Загруженные книги доступны к скачиванию всем зарегистри
   ровавшимся пользователям, к просмотру они доступны без регистрации.
   Так же есть администраторы, которые проверяют контент, могут добавлять, редактировать, удалять книги. */

/* Project "Open Library"
    In this project, users can download and upload books. Downloaded books are available for download to all registered
   for experienced users, they are available for viewing without registration.
   There are also administrators who check the content, can add, edit, delete books.*/

CREATE TABLE genres (
    id INT  primary key GENERATED ALWAYS AS IDENTITY,
    genre VARCHAR(100)
);

CREATE TABLE authors (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    author VARCHAR(100)
);

CREATE TABLE characters (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    character VARCHAR(10)
);

CREATE TABLE users (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    firstName VARCHAR(30),
    lastName VARCHAR(30),
    email VARCHAR(30),
    password VARCHAR(30),
    age INT,
    loading_books INTEGER,
    registration_date DATE DEFAULT CURRENT_DATE,
    character_id INT NOT NULL REFERENCES characters(id)
);

CREATE TABLE books (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(100),
    description VARCHAR(200),
    loading_date DATE DEFAULT CURRENT_DATE,
    author_id INT NOT NULL REFERENCES authors(id),
    genre_id INT NOT NULL REFERENCES genres(id),
    loading_id INT NOT NULL REFERENCES users(id)
);