export interface Task {
  id: string;
  title: string;
  status: string;
  priority: string;
  assignee: string;
  tags: string[];
  links: string[];
  created: string;
  updated: string;
  body: string;
}

export interface ColumnData {
  status: string;
  tasks: Task[];
}

export interface ProjectConfig {
  name: string;
  prefix: string;
}

export interface Config {
  project: ProjectConfig;
  statuses: string[];
  priorities: string[];
  auto_commit: boolean;
}

export interface BoardData {
  columns: ColumnData[];
  config: Config;
  repoName: string;
}

export interface BoardFilter {
  priority: string;
  assignee: string;
  tag: string;
}

export interface RepoEntry {
  name: string;
  path: string;
  remote: string;
  dir: string;
}

export interface CreateTaskInput {
  title: string;
  priority: string;
  assignee: string;
  tags: string[];
  body: string;
}

export interface UpdateTaskInput {
  id: string;
  title: string;
  status: string;
  priority: string;
  assignee: string;
  tags: string[];
  body: string;
}

export interface Notification {
  id: number;
  type: 'success' | 'error' | 'info';
  message: string;
}
