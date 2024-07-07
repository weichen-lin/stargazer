import type { Repository, ILanguageDistribution, IRepoAtDashboard, ISearchKey, IRepoDetail } from './repos'
import type { IUserConfig, IUserCrontab } from './user'
import { fetcher } from './fetcher'
import {
  getUserRepos,
  getUserStarsRelation,
  getUserStarsRelationRepos,
  getLanguageDistribution,
  getRepositoriesCount,
  getReposByKey,
  getRepoDetail,
  deleteRepo,
  getFullTextSearch,
} from './repos'
import { getUserInfo, updateInfo, getCrontabInfo } from './user'
import { getTags, getTagsByRepo, createTag, deleteTagByRepo } from './tags'

export { Repository, IUserConfig, ILanguageDistribution, IRepoAtDashboard, ISearchKey, IUserCrontab, IRepoDetail }
export { fetcher }
export {
  getUserRepos,
  getUserStarsRelation,
  getUserStarsRelationRepos,
  getLanguageDistribution,
  getRepositoriesCount,
  getReposByKey,
  getRepoDetail,
  deleteRepo,
  getFullTextSearch,
}
export { getTags, getTagsByRepo, createTag, deleteTagByRepo }
export { getUserInfo, updateInfo, getCrontabInfo }
