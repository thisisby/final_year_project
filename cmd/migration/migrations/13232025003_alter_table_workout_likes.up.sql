ALTER TABLE IF EXISTS workout_likes
DROP CONSTRAINT IF EXISTS unique_workout_user_like;

ALTER TABLE IF EXISTS workout_likes
    ADD CONSTRAINT unique_workout_user_like
    UNIQUE (workout_id, user_id, deleted_at);