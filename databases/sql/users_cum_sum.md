# Users cumulaative sum
Let's imagine that we have "time-series" DB, which collects events. In this case `users` with fields
- `ts` - timestamp
- `user_cnt` - number of users in this event
- `ev_type` - type of event

```sql
--- DDL
CREATE TEMP TABLE IF NOT EXISTS
    users(
        ts bigint,
        user_cnt bigint NOT NULL CHECK(user_cnt <> 0),
        ev_type smallint NOT NULL);
```

The task is: **`find timestamp when simultaneus number of users was the higest`**.

*Query to solve this task*:

```sql
--- DQL
WITH user_traffic_with_cumulative_amount AS (
  SELECT
    ts,
    user_cnt,
    SUM(CASE WHEN ev_type = 1 THEN user_cnt ELSE -user_cnt END)
      OVER (ORDER BY ts) AS cumulative_users
  FROM users
)
SELECT DISTINCT
  ts,
  cumulative_users
FROM user_traffic_with_cumulative_amount
WHERE cumulative_users = (
  SELECT MAX(cumulative_users) FROM user_traffic_with_cumulative_amount
);
```
Another approach is to drop "case when" and use the magic of math:
```sql
WITH user_traffic_with_cumulative_amount AS (
    SELECT
      ts, 
      SUM(user_cnt*(2*ev_type-1)) OVER (ORDER BY ts)
    FROM users
)
--- as in previous example
--- ...
```
Here, event type is numeric, so:
- *2 \* OUT - 1 = -1* (ev_type OUT = 0)
- *2 * IN -1 = 1* (ev_type OUT = 0)

Which variant to use is up to you. EXPLAIN ANALYZE shows no difference on this queries, but the amount of rows is small (100), so performance can differ on higher amount.
