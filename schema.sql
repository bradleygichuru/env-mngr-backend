CREATE TABLE Bucket (
	bucket_id INT GENERATED ALWAYS AS IDENTITY,
	name text not null,
	user_id integer,
	organization_id integer,
	envVariables json not null,
	PRIMARY KEY(bucket_id),
	CONSTRAINT fk_organization FOREIGN KEY(organization_id) REFERENCES Organization(organization_id)  ON DELETE SET NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE SET NULL

);

CREATE TYPE valid_roles AS ENUM ('Admin', 'User');

CREATE TABLE Users (
	user_id INT GENERATED ALWAYS AS IDENTITY,
	organization_id integer,
	email text not null,
	password text not null,
	role valid_roles,
	PRIMARY KEY(user_id),
	CONSTRAINT fk_organization FOREIGN KEY(organization_id) REFERENCES Organization(organization_id) ON DELETE SET NULL

);

CREATE TABLE Organization (
	organization_id INT GENERATED ALWAYS AS IDENTITY,
	name text not null,
	PRIMARY KEY(organization_id)
);
