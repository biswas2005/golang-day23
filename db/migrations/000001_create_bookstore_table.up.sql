CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL
);
CREATE TABLE authors (
    auid INT AUTO_INCREMENT PRIMARY KEY,
    aname VARCHAR(255) NOT NULL
);
ALTER TABLE books
ADD COLUMN author_id INT,
ADD CONSTRAINT fk_books_author
FOREIGN KEY (author_id) REFERENCES authors(auid);
CREATE TABLE orders (
    oid INT AUTO_INCREMENT PRIMARY KEY,
    book_id INT NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    CONSTRAINT fk_orders_book
    FOREIGN KEY (book_id) REFERENCES books(id)
);
