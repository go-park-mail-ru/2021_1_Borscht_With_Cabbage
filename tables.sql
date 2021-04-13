DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS restaurants CASCADE;
DROP TABLE IF EXISTS sections CASCADE;
DROP TABLE IF EXISTS dishes CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS baskets CASCADE;
DROP TABLE IF EXISTS baskets_food;
DROP TABLE IF EXISTS basket_users;
DROP TABLE IF EXISTS basket_orders;

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
                            deliveryCost INTEGER DEFAULT 0,
                            avgCheck INTEGER DEFAULT 0,
                            description TEXT,
                            rating FLOAT DEFAULT 0,
                            avatar TEXT
);

CREATE TABLE sections (
                            sid SERIAL PRIMARY KEY,
                            restaurant INTEGER REFERENCES restaurants(rid) ON DELETE CASCADE,
                            name TEXT
);

CREATE TABLE dishes (
                        did SERIAL PRIMARY KEY,
                        restaurant TEXT REFERENCES restaurants(name) ON DELETE CASCADE,
                        restaurantId INTEGER REFERENCES restaurants(rid) ON DELETE CASCADE,
                        section INTEGER REFERENCES sections(sid) ON DELETE CASCADE,
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
                         sum INTEGER,
                         restaurant TEXT REFERENCES restaurants(name) ON DELETE CASCADE
);


CREATE TABLE baskets_food (
                              basket INTEGER REFERENCES baskets(bid) ON DELETE CASCADE,
                              dish INTEGER REFERENCES dishes(did) ON DELETE CASCADE
);

CREATE TABLE basket_users (
                              basketID INTEGER REFERENCES baskets(bid) ON DELETE CASCADE,
                              userID INTEGER REFERENCES  users(uid) ON DELETE CASCADE
);

CREATE TABLE basket_orders(
                              basketID INTEGER REFERENCES baskets(bid) ON DELETE CASCADE,
                              orderID INTEGER REFERENCES orders(oid) ON DELETE CASCADE -- любо уже сформированному заказу
);


GRANT ALL PRIVILEGES ON TABLE users TO delivery;
GRANT ALL PRIVILEGES ON TABLE addresses TO delivery;
GRANT ALL PRIVILEGES ON TABLE restaurants TO delivery;
GRANT ALL PRIVILEGES ON TABLE sections TO delivery;
GRANT ALL PRIVILEGES ON TABLE dishes TO delivery;
GRANT ALL PRIVILEGES ON TABLE orders TO delivery;
GRANT ALL PRIVILEGES ON TABLE baskets TO delivery;
GRANT ALL PRIVILEGES ON TABLE baskets_food TO delivery;
GRANT ALL PRIVILEGES ON TABLE basket_users TO delivery;
GRANT ALL PRIVILEGES ON TABLE basket_orders TO delivery;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO delivery;
