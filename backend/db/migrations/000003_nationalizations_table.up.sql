CREATE TABLE IF NOT EXISTS nationalizations (
	nationalization_id SERIAL     PRIMARY KEY,
	country_code       VARCHAR(4) NOT NULL,
	probability        NUMERIC    NOT NULL
);

CREATE TABLE IF NOT EXISTS people_nationalizations (
	id                 SERIAL  PRIMARY KEY,
	people_id          INTEGER NOT NULL REFERENCES peoples(people_id),
	nationalization_id INTEGER NOT NULL REFERENCES nationalizations(nationalization_id),
	UNIQUE(people_id, nationalization_id)
);

