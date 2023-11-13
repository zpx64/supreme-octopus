CREATE TABLE IF NOT EXISTS peoples (
	people_id  SERIAL       PRIMARY KEY,
	name       VARCHAR(256) NOT NULL,
	surname    VARCHAR(256) NOT NULL,
	patronymic VARCHAR(256),
	age        INTEGER      NOT NULL
);

