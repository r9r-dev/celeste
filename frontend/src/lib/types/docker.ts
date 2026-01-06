export interface Container {
	id: string;
	name: string;
	image: string;
	status: 'running' | 'stopped' | 'paused' | 'restarting';
	state: string;
	cpu: number;
	memory: number;
	memoryLimit: number;
	networkRx: number;
	networkTx: number;
	created: string;
	ports: Port[];
}

export interface Port {
	private: number;
	public: number;
	type: 'tcp' | 'udp';
}

export interface Stack {
	name: string;
	path: string;
	status: 'running' | 'partial' | 'stopped';
	containers: Container[];
	services: number;
	runningServices: number;
}

export interface Volume {
	name: string;
	driver: string;
	mountpoint: string;
	size: number;
	usedBy: string[];
}

export interface Network {
	id: string;
	name: string;
	driver: string;
	scope: string;
	containers: string[];
}

export interface SystemStats {
	cpuUsage: number;
	cpuCores: number;
	memoryUsed: number;
	memoryTotal: number;
	diskUsed: number;
	diskTotal: number;
	uptime: number;
	containersRunning: number;
	containersStopped: number;
	images: number;
	volumes: number;
	networks: number;
}
