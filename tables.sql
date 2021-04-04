DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS restaurants CASCADE;
DROP TABLE IF EXISTS dishes CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS baskets CASCADE;
DROP TABLE IF EXISTS baskets_food;

CREATE TABLE users (
                       uid SERIAL PRIMARY KEY,
                       name TEXT,
                       phone TEXT,
                       email TEXT,
                       photo TEXT,
    -- mainAddress text references addresses(address) on delete cascade ,
                       password TEXT

);

CREATE TABLE addresses (
                           address TEXT,
                           "user" INTEGER REFERENCES users(uid) ON DELETE CASCADE
);

CREATE TABLE restaurants (
                             rid SERIAL PRIMARY KEY,
                             name TEXT UNIQUE,
                             adminEmail TEXT,
                             adminPhone TEXT,
                             adminPassword TEXT,
                             deliveryCost INTEGER,
                             avgCheck INTEGER,
                             description TEXT,
                             rating FLOAT,
                             avatar TEXT
);

CREATE TABLE dishes (
                        did SERIAL PRIMARY KEY,
                        restaurant TEXT REFERENCES restaurants(name) ON DELETE CASCADE,
                        name TEXT,
                        price INTEGER,
                        weight INTEGER,
                        description TEXT,
                        image TEXT
);

CREATE TABLE orders (
                        oid SERIAL PRIMARY KEY,
                        restaurant TEXT REFERENCES restaurants(name) ON DELETE CASCADE,
                        userID INTEGER REFERENCES  users(uid) ON DELETE CASCADE,
                        orderTime TIMESTAMP,
                        address TEXT,
                        deliveryCost INTEGER,
                        sum INTEGER,
                        status TEXT,
                        deliveryTime TIME
);

CREATE TABLE baskets (
                         bid SERIAL PRIMARY KEY,
                         restaurant TEXT REFERENCES restaurants(name) ON DELETE CASCADE,
                         userID INTEGER REFERENCES  users(uid) ON DELETE CASCADE, -- корзина либо принадлежит юзеру и пока не оформлена
                         orderID INTEGER REFERENCES orders(oid) ON DELETE CASCADE -- любо уже сформированному заказу
);


CREATE TABLE baskets_food (
                              basket INTEGER REFERENCES baskets(bid) ON DELETE CASCADE,
                              dish INTEGER REFERENCES dishes(did) ON DELETE CASCADE
)
