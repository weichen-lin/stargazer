interface Repository {
  id: number;
  name: string;
  owner_name: string;
  avatar_url: string;
  html_url: string;
  homepage: string;
  description: string;
  created_at: string;
  updated_at: string;
  synced_at: string;
  watchers: number;
  forks: number;
  open_issues: number;
  language: string;
  archived: boolean;
  topics: string[];
}

export type { Repository };
