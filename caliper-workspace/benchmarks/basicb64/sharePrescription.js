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
    	let args = {
            contractId: 'basicb64',
            contractVersion: '1.0',
            contractFunction: 'SharePrescription',
            contractArguments: ["0", "Admin"],
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