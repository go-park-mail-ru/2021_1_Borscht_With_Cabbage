DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS restaurants CASCADE;
DROP TABLE IF EXISTS dishes;

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
    name TEXT,
    adminEmail TEXT,
    adminPhone TEXT,
    adminPassword TEXT,
    deliveryCost INTEGER,
    avgCheck INTEGER,
    description TEXT,
    rating FLOAT,
    avatar TEXT
);

CREATE TABLE sections (
    sid SERIAL PRIMARY KEY,
    restaurant INTEGER REFERENCES restaurants(rid) ON DELETE CASCADE,
    name TEXT
);

CREATE TABLE dishes (
    did SERIAL PRIMARY KEY,
    restaurant INTEGER REFERENCES restaurants(rid) ON DELETE CASCADE,
    section INTEGER REFERENCES sections(sid) ON DELETE CASCADE,
    name TEXT,
    price INTEGER,
    weight INTEGER,
    description TEXT,
    image TEXT
);


-- GRANT ALL PRIVILEGES ON TABLE users TO delivery;
-- GRANT ALL PRIVILEGES ON TABLE addresses TO delivery;
-- GRANT ALL PRIVILEGES ON TABLE restaurants TO delivery;
-- GRANT ALL PRIVILEGES ON TABLE sections TO delivery;
-- GRANT ALL PRIVILEGES ON TABLE dishes TO delivery;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO delivery;
