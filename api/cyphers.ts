export const getUserRepos = `
MATCH (n:User {name: $username})-[:STARS]->(r:Repository)
WITH count(r) as total
MATCH (n:User {name: $username})-[:STARS]->(r:Repository)
WITH total, r
ORDER BY r.created_at DESC
SKIP $limit * ($page - 1)
LIMIT $limit
WITH total, collect(r) as limitedRepositories
RETURN total, limitedRepositories;
`
