export namespace config {
	
	export class ProjectConfig {
	    name: string;
	    prefix: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.prefix = source["prefix"];
	    }
	}
	export class Config {
	    project: ProjectConfig;
	    statuses: string[];
	    priorities: string[];
	    auto_commit: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project = this.convertValues(source["project"], ProjectConfig);
	        this.statuses = source["statuses"];
	        this.priorities = source["priorities"];
	        this.auto_commit = source["auto_commit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class ColumnData {
	    status: string;
	    tasks: task.Task[];
	
	    static createFrom(source: any = {}) {
	        return new ColumnData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.tasks = this.convertValues(source["tasks"], task.Task);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BoardData {
	    columns: ColumnData[];
	    config?: config.Config;
	    repoName: string;
	
	    static createFrom(source: any = {}) {
	        return new BoardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columns = this.convertValues(source["columns"], ColumnData);
	        this.config = this.convertValues(source["config"], config.Config);
	        this.repoName = source["repoName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BoardFilter {
	    priority: string;
	    assignee: string;
	    tag: string;
	
	    static createFrom(source: any = {}) {
	        return new BoardFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.priority = source["priority"];
	        this.assignee = source["assignee"];
	        this.tag = source["tag"];
	    }
	}
	
	export class CreateTaskInput {
	    title: string;
	    priority: string;
	    assignee: string;
	    tags: string[];
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateTaskInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.priority = source["priority"];
	        this.assignee = source["assignee"];
	        this.tags = source["tags"];
	        this.body = source["body"];
	    }
	}
	export class RepoEntry {
	    name: string;
	    path: string;
	    remote: string;
	    dir: string;
	
	    static createFrom(source: any = {}) {
	        return new RepoEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.remote = source["remote"];
	        this.dir = source["dir"];
	    }
	}
	export class UpdateTaskInput {
	    id: string;
	    title: string;
	    status: string;
	    priority: string;
	    assignee: string;
	    tags: string[];
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTaskInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.assignee = source["assignee"];
	        this.tags = source["tags"];
	        this.body = source["body"];
	    }
	}

}

export namespace task {
	
	export class Task {
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
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.assignee = source["assignee"];
	        this.tags = source["tags"];
	        this.links = source["links"];
	        this.created = source["created"];
	        this.updated = source["updated"];
	        this.body = source["body"];
	    }
	}

}

