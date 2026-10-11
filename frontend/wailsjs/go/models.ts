export namespace main {
	
	export class JSONStats {
	    keys: number;
	    nodes: number;
	    depth: number;
	    objects: number;
	    arrays: number;
	    strings: number;
	    numbers: number;
	    bools: number;
	    nulls: number;
	
	    static createFrom(source: any = {}) {
	        return new JSONStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keys = source["keys"];
	        this.nodes = source["nodes"];
	        this.depth = source["depth"];
	        this.objects = source["objects"];
	        this.arrays = source["arrays"];
	        this.strings = source["strings"];
	        this.numbers = source["numbers"];
	        this.bools = source["bools"];
	        this.nulls = source["nulls"];
	    }
	}
	export class AnalyzeResult {
	    valid: boolean;
	    message: string;
	    line: number;
	    col: number;
	    stats?: JSONStats;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.message = source["message"];
	        this.line = source["line"];
	        this.col = source["col"];
	        this.stats = this.convertValues(source["stats"], JSONStats);
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

