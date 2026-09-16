export namespace config {
	
	export class Config {
	    credentials_file_path: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.credentials_file_path = source["credentials_file_path"];
	    }
	}
	export class Credentials {
	    server: string;
	    username: string;
	    password: string;
	    group: number;
	    otp_secret: string;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.group = source["group"];
	        this.otp_secret = source["otp_secret"];
	    }
	}

}

export namespace updater {
	
	export class Release {
	    tag_name: string;
	    name: string;
	    body: string;
	    html_url: string;
	    published_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Release(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.name = source["name"];
	        this.body = source["body"];
	        this.html_url = source["html_url"];
	        this.published_at = source["published_at"];
	    }
	}

}

