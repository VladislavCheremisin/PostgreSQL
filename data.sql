-- Insert data

-- genres
INSERT INTO genres (genre)
VALUES
('fantasy'),
('comedy'),
('horror')

-- authors
INSERT INTO authors (author)
VALUES
('Pushkin'),
('Dostoevskij'),
('Gogol')

-- characters
INSERT INTO characters (character)
VALUES
('Admin'),
('User')

-- users_password
INSERT INTO users_password (password)
VALUES
('1234556'),
('7654321')

-- users
INSERT INTO users (firstName, lastName, email, loading_books, character_id)
VALUES
('Vasya', 'Pupkin', 'test@email.ru', 7, 1),
('Petya', 'Nikonov', 'testNikon@email.ru', 2, 2)

-- error users
INSERT INTO users (firstName, lastName, email, loading_books, character_id)
VALUES
('Va', 'Pu', 'test@email.ru', 7, 1),
('Petya', 'Nikonov', 'testNikon@email.ru', 2, 2)

-- books
INSERT INTO books (title, description, loading_id)
VALUES
('Горе от ума', 'Очень интересная книга', 1),
('Тень Чернобыля', 'О сталкерах и зоне отчуждения', 2)

-- books_authors
INSERT INTO books_authors (book_id, author_id)
VALUES
(1, 1),
(2, 2)

-- books_genres
INSERT INTO books_genres (book_id, genre_id)
VALUES
(1, 1),
(1, 2),
(2, 1),
(2, 2);