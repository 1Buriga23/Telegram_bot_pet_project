CREATE TABLE categories (

    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL

);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,

    category_id INTEGER NOT NULL,

    FOREIGN KEY(category_id)
        REFERENCES categories(id)
);

CREATE TABLE videos (

    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,

    file_id TEXT NOT NULL,

    item_id INTEGER NOT NULL,

    FOREIGN KEY(item_id)
        REFERENCES items(id)
        ON DELETE CASCADE

);
