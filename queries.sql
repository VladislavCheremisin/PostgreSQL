-- Выводим пользовотеля по частичному совпадению email
-- Display user by partial match of email
SELECT email, firstName, lastName
FROM users
WHERE email LIKE 'test%'
ORDER BY firstName ASC;

-- Выводим перечень книг и авторов этих книг
-- Display a list of books and authors of these books
SELECT books.title, authors.author
FROM books_authors, books, authors
WHERE books_authors.book_id = books.id AND books_authors.author_id = authors.id;

-- Выводим колличество книг каждого автора
-- Display the number of books by each author
SELECT authors.author, COUNT(authors.author) AS books
FROM books_authors, books, authors
WHERE books_authors.book_id = books.id AND books_authors.author_id = authors.id
GROUP BY authors.author;

