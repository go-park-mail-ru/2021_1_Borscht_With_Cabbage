DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS restaurants CASCADE;
DROP TABLE IF EXISTS dishes;

CREATE TABLE users (
                       uid SERIAL PRIMARY KEY,
                       name text,
                       phone text,
                       email text,
                       photo text,
    -- mainAddress text references addresses(address) on delete cascade ,
                       password text
);

CREATE TABLE sessions (
                          session text NOT NULL PRIMARY KEY,
                          uid INTEGER REFERENCES users(uid) ON DELETE CASCADE
);

CREATE TABLE addresses (
                           address text,
                           "user" integer references users(uid) on delete cascade
);

CREATE TABLE restaurants (
                             rid SERIAL PRIMARY KEY,
                             name text,
                             deliveryCost integer,
                             avgCheck integer,
                             description text,
                             rating float,
                             avatar text
);

CREATE TABLE dishes (
                        did SERIAL PRIMARY KEY,
                        restaurant integer references restaurants(rid) on delete cascade,
                        name text,
                        price integer,
                        weight integer,
                        description text,
                        image text
)