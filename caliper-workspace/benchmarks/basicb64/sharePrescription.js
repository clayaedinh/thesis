'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class SharePrescriptionWorkload extends WorkloadModuleBase {
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
        var caliperAdminObscured = "5a6a59ea61c3e616cbe67281399d1991cf5efe32d253203698ad63fea1daf9ae"
        var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
    	let args = {
            contractId: 'basicb64',
            contractVersion: '1.0',
            contractFunction: 'SharePrescription',
            contractArguments: [pid.toString(), caliperAdminObscured],
            timeout: 30,
        };
        await this.sutAdapter.sendRequests(args);
    }
}
 

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new SharePrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;