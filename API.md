## Run

You can start the application using the following command:

```shell
make run
```

The default API port is `3000`.

A PostgreSQL database pointing to the same seeded database will be builded.

## Endpoints

The following endpoints will be reachable through the API:

### Healthcheck

```
/healthcheck
```
Used to make sure that all dependencies from the application are reachable, like PostgreSQL.

### Available shifts from a worker

```
/v1/workers/:worker_id/available_shifts?end=:end_date&start=:start_date&limit=:limit&cursor=:cursor
```
Retrives all available shifts from the given worker, if any. Where:
- `:worker_id`: **Integer** worker ID to retrieve shifts
- `:start_date` & `:end_date`: **Date (ISO 8601)** start and end date to filter results
- `:limit`: **Integer** limits how many results will be retrieved
- `:cursor`: **String** used to go through result pages

## Testing

There are two types of tests in the application: unit and integration.

### Unit

This short version won't test external services like the PostgreSQL layer.

Use the following command to run:
```shell
make test-short
```

### Integration

This is a more long running process and will assert that connection with external services like PostgreSQL properly work.

To run the integration tests, a properly PostgreSQL database should be up and running, otherwise the tests will be skipped.

Note: this also runs unit tests.

Use the following command to run:
```shell
make test
```

## Improvements

- Indexes:
  - The main query should be analysed using `EXPLAIN ANALYSE` and a closer look on possible missing indexes.
    - Missing indexes makes query execution slower.
- PostgreSQL tuning:
  - Edit configs to a more reasonable values so the query execution could be speed up.
    - Increase max allocation memory;
    - Increase max parallel workers;
- Cache:
  - A cache layer could be implemented, using Redis or Elasticsearch.
    - Keeping a separated datasource that holds only the available shifts that workers can apply would speed up endpoint.
    - As soon as a shift would be applyed, this should be removed from this secondary datasource.
  - PgPool could be used as a transparent cache layer.
    - Cache invalidation would be easier.
    - More befenifits like load balancing out of the box.
- Data:
  - Create trigger to don't allow shifts with start and end with second and millisecond precision be inserted.
    - The amount of data would be lower.
