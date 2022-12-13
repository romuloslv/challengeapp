-- name: CreateAccount :one
INSERT INTO accounts (person_id, first_name, last_name, web_address, date_birth)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE person_id = $1
LIMIT 1;

-- name: UpdateAccount :one
UPDATE accounts
SET first_name  = $2,
    last_name   = $3,
    web_address = $4,
    date_birth  = $5
WHERE person_id = $1
RETURNING *;

-- name: PartialUpdateAccount :one
UPDATE accounts
SET first_name  = CASE WHEN @update_first_name::boolean THEN @first_name::VARCHAR(30) ELSE first_name END,
    last_name   = CASE WHEN @update_last_name::boolean THEN @last_name::VARCHAR(20) ELSE last_name END,
    web_address = CASE WHEN @update_web_address::boolean THEN @web_address::VARCHAR(50) ELSE web_address END,
    date_birth  = CASE WHEN @update_date_birth::boolean THEN @date_birth::VARCHAR(10) ELSE date_birth END
WHERE person_id = @person_id
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE person_id = $1;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY first_name;