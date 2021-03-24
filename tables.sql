
CREATE TABLE users (
    uid integer NOT NULL PRIMARY KEY,
    name text,
    phone text,
    email text,
    photo text,
    mainAddress text references addresses(address) on delete cascade ,
    password text
);

CREATE TABLE sessions (
    session  integer NOT NULL PRIMARY KEY,
    uid INTEGER REFERENCES users(uid) ON DELETE CASCADE,
);

CREATE TABLE addresses (
    address text,
    user integer references users(uid) on delete cascade
);

CREATE TABLE restaurants (
    rid integer not null primary key,
    name text,
    deliveryCost integer,
    avgCheck integer,
    description text,
    rating float,
    avatar string,
)

CREATE TABLE dishes (
    did integer not null primary key,
    restaurant integer references restaurants(rid) on delete cascade
    name text,
    price integer,
    weight integer,
    description text,
    image text,
)