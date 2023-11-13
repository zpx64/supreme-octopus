CREATE TABLE IF NOT EXISTS genders (
	people_id   SERIAL  NOT NULL REFERENCES peoples(people_id),
	gender      INTEGER NOT NULL,
	probability NUMERIC NOT NULL
);
