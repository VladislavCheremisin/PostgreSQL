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
    genre VARCHAR(100) NOT NULL
);

CREATE TABLE authors (
     id INT primary key GENERATED ALWAYS AS IDENTITY,
     author VARCHAR(100) NOT NULL
);

CREATE TABLE characters (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    character VARCHAR(10) NOT NULL
);

CREATE TABLE users_password (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    password VARCHAR(30) NOT NULL,
    CHECK (length(password) > 3)
);

CREATE TABLE users (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    firstName VARCHAR(30) NOT NULL,
    lastName VARCHAR(30) NOT NULL,
    email VARCHAR(30) NOT NULL,
    loading_books INTEGER NOT NULL,
    registration_date DATE DEFAULT CURRENT_DATE,
    character_id INT NOT NULL REFERENCES characters(id),
    CHECK (length(firstName) > 3),
    CHECK (length(lastName) > 3)
);

CREATE TABLE books (
    id INT primary key GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(600) NOT NULL,
    loading_date DATE DEFAULT CURRENT_DATE,
    loading_id INT NOT NULL REFERENCES users(id),
    CHECK (length(title) > 5),
    CHECK (length(description) > 10)
);

CREATE TABLE books_authors (
    book_id INT NOT NULL,
    author_id INT NOT NULL
);

CREATE TABLE books_genres (
    book_id INT NOT NULL,
    genre_id INT NOT NULL
);

-- Indexes
CREATE INDEX concurrently authors_author_idx
    ON authors (author);

CREATE INDEX concurrently users_firstName_idx
    ON users (firstName);

CREATE INDEX concurrently users_email_idx
    ON users (email);