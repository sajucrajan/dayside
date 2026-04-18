export namespace detect {
	
	export class BrowserTab {
	    browserPid: number;
	    browserName: string;
	    title: string;
	    url: string;
	    incognito: boolean;
	    kind: string;
	    severity: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new BrowserTab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browserPid = source["browserPid"];
	        this.browserName = source["browserName"];
	        this.title = source["title"];
	        this.url = source["url"];
	        this.incognito = source["incognito"];
	        this.kind = source["kind"];
	        this.severity = source["severity"];
	        this.reason = source["reason"];
	    }
	}
	export class DeviceEntry {
	    name: string;
	    virtual: boolean;
	    severity: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.virtual = source["virtual"];
	        this.severity = source["severity"];
	        this.reason = source["reason"];
	    }
	}
	export class DeviceReport {
	    audio: DeviceEntry[];
	    video: DeviceEntry[];
	
	    static createFrom(source: any = {}) {
	        return new DeviceReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.audio = this.convertValues(source["audio"], DeviceEntry);
	        this.video = this.convertValues(source["video"], DeviceEntry);
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
	export class WindowEntry {
	    hwnd: number;
	    title: string;
	    topmost: boolean;
	    layered: boolean;
	    affinity: number;
	    cloaked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WindowEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hwnd = source["hwnd"];
	        this.title = source["title"];
	        this.topmost = source["topmost"];
	        this.layered = source["layered"];
	        this.affinity = source["affinity"];
	        this.cloaked = source["cloaked"];
	    }
	}
	export class ProcessInfo {
	    pid: number;
	    name: string;
	    path: string;
	    windows: WindowEntry[];
	    signed: string;
	    signer: string;
	    protected: boolean;
	    severity: string;
	    flags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.windows = this.convertValues(source["windows"], WindowEntry);
	        this.signed = source["signed"];
	        this.signer = source["signer"];
	        this.protected = source["protected"];
	        this.severity = source["severity"];
	        this.flags = source["flags"];
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
	export class SystemReport {
	    monitorCount: number;
	    remoteSession: boolean;
	    remoteTools: string[];
	    hostName: string;
	    userName: string;
	    osVersion: string;
	    platform: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.monitorCount = source["monitorCount"];
	        this.remoteSession = source["remoteSession"];
	        this.remoteTools = source["remoteTools"];
	        this.hostName = source["hostName"];
	        this.userName = source["userName"];
	        this.osVersion = source["osVersion"];
	        this.platform = source["platform"];
	    }
	}

}

export namespace main {
	
	export class ScanResult {
	    processes: detect.ProcessInfo[];
	    tabs: detect.BrowserTab[];
	    devices: detect.DeviceReport;
	    system: detect.SystemReport;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.processes = this.convertValues(source["processes"], detect.ProcessInfo);
	        this.tabs = this.convertValues(source["tabs"], detect.BrowserTab);
	        this.devices = this.convertValues(source["devices"], detect.DeviceReport);
	        this.system = this.convertValues(source["system"], detect.SystemReport);
	        this.warnings = source["warnings"];
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

