export namespace novels {
	
	export class Novel {
	    Title: string;
	    Debug: string[];
	
	    static createFrom(source: any = {}) {
	        return new Novel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Debug = source["Debug"];
	    }
	}

}

