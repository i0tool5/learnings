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
    ev_type,
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
