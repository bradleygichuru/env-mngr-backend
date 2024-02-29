-- name: CreateBucket :one
INSERT INTO Bucket (
 name,envVariables) VALUES(
  $1,$2
 )
RETURNING *;

-- name: CreateUser :one
INSERT INTO Users (
 email,password) VALUES(
  $1,$2
 )
RETURNING *;

-- name: CreateOrganisation :one
INSERT INTO Organization (
 name) VALUES(
  $1
 )
RETURNING *;

-- name: GetBucket :one
SELECT * FROM Bucket
WHERE bucket_id = $1 LIMIT 1;
