curl -iX POST localhost:8080/accounts \
     -H 'Content-Type: application/json' \
     -d '{"person_id": "01234567890", "first_name": "TEST", "last_name": "TESTS", "web_address": "test@test.com", "date_birth": "2022-12-31"}'

curl -iX PUT localhost:8080/accounts/01234567890 \
     -H 'Content-Type: application/json' \
     -d '{"person_id": "01234567890", "first_name": "TEST", "last_name": "TESTS UPDATE", "web_address": "test@test.com", "date_birth": "2023-01-01"}'

curl -iX PATCH localhost:8080/accounts/01234567890 \
     -H 'Content-Type: application/json' \
     -d '{"person_id": "01234567890", "first_name": "UPDATE TEST", "last_name": "TESTS UPDATE", "web_address": "tests@test.com", "date_birth": "2023-01-01"}'

curl -iX DELETE localhost:8080/accounts/01234567890 \
     -H 'Content-Type: application/json'

curl -iX GET localhost:8080/
curl -iX GET localhost:8080/version
curl -iX GET localhost:8080/health
curl -iX GET localhost:8080/accounts
curl -iX GET localhost:8080/accounts/01234567890