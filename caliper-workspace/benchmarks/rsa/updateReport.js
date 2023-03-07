'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class UpdateReportWorkload extends WorkloadModuleBase {
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
        var b64report = "Dv+DBAEC/4QAAQwBDAAABP+EAAA="
    	var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
        let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'UpdateReport',
            contractArguments: [pid.toString(), b64report],
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
    return new UpdateReportWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;