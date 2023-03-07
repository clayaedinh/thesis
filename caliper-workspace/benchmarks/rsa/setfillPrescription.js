'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class SetfillPrescriptionWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
    }
     /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
    	this.txIndex++;
        var b64setfill = "Dv+DBAEC/4QAAQwBDAAA/gGc/4QAAUA1YTZhNTllYTYxYzNlNjE2Y2JlNjcyODEzOTlkMTk5MWNmNWVmZTMyZDI1MzIwMzY5OGFkNjNmZWExZGFmOWFl/gFUSW5zVTF0OTM5akNmM0xDNEVWcGRQcjk4MmNsNUhEWjRKT09rNjB1SDlubXlMOUp3ZzYvMWZUWUFseGF5YUtxc2hBRmlQcVc4M0pndlRtcVc2aHd6cnJsZWprQ0k3UkZPTTY3eiszc21aQkJYMGRsOTVZaUpmR0V5NnhzQ3ljVVFrejRzWXoxNTBvMUxjS1dBeE4ycUtLMlBSNisrSEtUYy9iOVZmRWhvQ0E1eU54Snd6cDgxbzF0NHFWK2diSTBPTUJZNEkxMDNKTllnLzIzTEVKZmRYWjJzdGdzVHZ2ajQrKzhOZWFRQUZjemlzVU5BTHdKMzhnK2JYa1ZYekF1VmxFc3JwSk9LYXBkRVZGa20zcy95Y1NLVDU0c1VuWERFaERSeTBhRUR2TkhwNEVuQkJCL1BwN3ZhTjZxdXpZQzcyMlNySlNiYWJ2cXZFZ1VEcm5rNA=="
    	var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
        let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'SetfillPrescription',
            contractArguments: [pid.toString(), b64setfill],
            timeout: 30
        };
        await this.sutAdapter.sendRequests(args);
    }
}
 

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new SetfillPrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;