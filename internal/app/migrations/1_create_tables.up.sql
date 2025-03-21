CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    email      TEXT UNIQUE                            NOT NULL,
    password   TEXT                                   NOT NULL,
    is_admin   BOOLEAN                  DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS categories
(
    id         SERIAL PRIMARY KEY,
    name       TEXT UNIQUE                            NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS books
(
    id          SERIAL PRIMARY KEY,
    title       TEXT                                   NOT NULL,
    author      TEXT                                   NOT NULL,
    year        INT                                    NOT NULL CHECK (books.year >= 0),
    price       FLOAT                                  NOT NULL CHECK (books.price >= 0),
    amount      INT                                    NOT NULL CHECK (books.amount >= 0),
    category_id INT                                    NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (category_id) REFERENCES categories (id),
    CONSTRAINT вввввввunique_author_title UNIQUE(author, title)
);

CREATE TABLE IF NOT EXISTS carts
(
    user_id    INTEGER                                NOT NULL PRIMARY KEY,
    book_ids   INTEGER[]                              NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS user_tokens
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);