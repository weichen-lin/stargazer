-- name: ListStarredRepositoriesByUser :many
SELECT
    -- Repository details
    r.id,
    r.name,
    r.html_url,
    r.homepage,
    r.description,
    r.watchers,
    r.forks,
    r.open_issues,
    r.language,
    r.archived,
    r.created_at,
    r.updated_at,
    -- Owner details (aliased to avoid potential conflicts)
    o.name AS owner_name,
    o.avatar_url AS owner_avatar_url,
    -- Aggregate topic names into an array and CAST the result
    COALESCE(
        ARRAY_AGG(t.name ORDER BY t.name) FILTER (WHERE t.id IS NOT NULL),
        ARRAY[]::VARCHAR[] -- Keep the empty array part consistent or also cast if needed
    )::TEXT[] AS topics -- <<< CAST the final result to TEXT[]
FROM
    stars s                     -- Start with the stars table
JOIN
    repositories r ON s.repo_id = r.id -- Join to get repository details
JOIN
    owner o ON r.owner_id = o.id       -- Join to get owner details
LEFT JOIN
    repository_topics rt ON r.id = rt.repo_id -- Left join to include repos with no topics
LEFT JOIN
    topics t ON rt.topic_id = t.id           -- Left join to get topic names
WHERE
    s.user_id = '0c28555f-5844-4c9d-9d13-f6d874d4d0c7'              -- Filter by the specific user's ID
AND
    s.is_delete = false         -- Only include active stars
GROUP BY
    r.id,                       -- Group by repository to aggregate topics
    o.id                        -- Include owner ID in group by (as its details are selected)
ORDER BY
    r.updated_at DESC;  

-- name: GetStarredLanguageDistributionByUser :many
SELECT
	r.language,
	COUNT(r.id) AS COUNT
FROM
	stars s
	JOIN repositories r ON s.repo_id = r.id
WHERE
	s.user_id = $1
	AND s.is_delete = FALSE
GROUP BY
	r.language
ORDER BY
	COUNT DESC;


-- name: GetRepositoriesByLanguagesWithPagination :many
WITH filtered AS (
    SELECT
        r.*
    FROM
        stars s
    JOIN
        repositories r ON s.repo_id = r.id
    WHERE
        s.user_id = $1
        AND s.is_delete = FALSE
        AND r.language = ANY($2::text[])
)
SELECT
    *,
    (SELECT COUNT(*) FROM filtered) AS total_count
FROM
    filtered
ORDER BY
    created_at DESC
LIMIT $3
OFFSET $4;
