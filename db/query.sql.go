// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
)

const createBucket = `-- name: CreateBucket :one
INSERT INTO Bucket (
 name,envVariables) VALUES(
  $1,$2
 )
RETURNING bucket_id, name, user_id, organization_id, envvariables
`

type CreateBucketParams struct {
	Name         string
	Envvariables []byte
}

func (q *Queries) CreateBucket(ctx context.Context, arg CreateBucketParams) (Bucket, error) {
	row := q.db.QueryRow(ctx, createBucket, arg.Name, arg.Envvariables)
	var i Bucket
	err := row.Scan(
		&i.BucketID,
		&i.Name,
		&i.UserID,
		&i.OrganizationID,
		&i.Envvariables,
	)
	return i, err
}

const createOrganisation = `-- name: CreateOrganisation :one
INSERT INTO Organization (
 name) VALUES(
  $1
 )
RETURNING organization_id, name
`

func (q *Queries) CreateOrganisation(ctx context.Context, name string) (Organization, error) {
	row := q.db.QueryRow(ctx, createOrganisation, name)
	var i Organization
	err := row.Scan(&i.OrganizationID, &i.Name)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO Users (
 email,password) VALUES(
  $1,$2
 )
RETURNING user_id, organization_id, email, password, role
`

type CreateUserParams struct {
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.OrganizationID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const getBucket = `-- name: GetBucket :one
SELECT bucket_id, name, user_id, organization_id, envvariables FROM Bucket
WHERE bucket_id = $1 LIMIT 1
`

func (q *Queries) GetBucket(ctx context.Context, bucketID int32) (Bucket, error) {
	row := q.db.QueryRow(ctx, getBucket, bucketID)
	var i Bucket
	err := row.Scan(
		&i.BucketID,
		&i.Name,
		&i.UserID,
		&i.OrganizationID,
		&i.Envvariables,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT user_id, organization_id, email, password, role FROM Users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.OrganizationID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}
