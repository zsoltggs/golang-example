# Validation Service

## How to run it

```
make build-docker
docker-compose up
```
Service will be exposed on port 11000

to run curl E2E curl requests:
- `cd services/validation-service`
- `./test.sh`

# Users

## How to run it

```
make build-docker
docker-compose up
```

For tests to work you need to run `docker-compose up`, it depends on mongo

## How to make better

- Add better unit testing
  - e.g. Repository level tests
- Define proper indexes in mongo
- Use a library for health
  - make it generic to check for multiple services
- Generalize Makefile
- Currently notifier service just logs things
  - Would be better to actually publish a message to a queue e.g. Kafka
- Currently Create performs an upsert
  - it still notifies users about user change
- Better testing e.g. proper mocking for svc layer, e2e testing


# Ports

- TODO

## Assumptions/Notes

- Not enough time to fix parsing to make sure it reads the input in batches
- Not enough time to make sure docker actually builds the code, so currently local machine builds the code and copies it there
- Because of privacy concerns I won't share this code on github even though this is a requirement, all the code that I write is confidential
- More testing is needed but not enough time to add them
  - e.g. full E2E integration tests