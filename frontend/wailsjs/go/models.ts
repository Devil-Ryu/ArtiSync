export namespace artlog {
	
	export class LoggerResult {
	    datetime: string;
	    level: string;
	    tag: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LoggerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.datetime = source["datetime"];
	        this.level = source["level"];
	        this.tag = source["tag"];
	        this.message = source["message"];
	    }
	}

}

