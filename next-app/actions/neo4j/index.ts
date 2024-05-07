import type { Repository, ILanguageDistribution } from './repos'
import type { UserInfo } from './user'
import fetcher from './fetcher'
import { getUserRepos, getUserStarsRelation, getUserStarsRelationRepos, getLanguageDistribution } from './repos'
import { getUserInfo, getUserProviderInfo, updateInfo } from './user'

export { Repository, UserInfo, ILanguageDistribution }
export { fetcher }
export { getUserRepos, getUserStarsRelation, getUserStarsRelationRepos, getLanguageDistribution }
export { getUserInfo, getUserProviderInfo, updateInfo }
