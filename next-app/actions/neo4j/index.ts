import type { Repository, ILanguageDistribution, IRepoAtDashboard } from './repos'
import type { UserInfo } from './user'
import fetcher from './fetcher'
import {
  getUserRepos,
  getUserStarsRelation,
  getUserStarsRelationRepos,
  getLanguageDistribution,
  getRepositoriesCount,
  getReposByKey,
} from './repos'
import { getUserInfo, getUserProviderInfo, updateInfo } from './user'

export { Repository, UserInfo, ILanguageDistribution, IRepoAtDashboard }
export { fetcher }
export {
  getUserRepos,
  getUserStarsRelation,
  getUserStarsRelationRepos,
  getLanguageDistribution,
  getRepositoriesCount,
  getReposByKey,
}
export { getUserInfo, getUserProviderInfo, updateInfo }
