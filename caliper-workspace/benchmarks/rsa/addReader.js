'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class AddReaderWorkload extends WorkloadModuleBase {
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
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'RegisterMeAsReportReader',
            contractArguments: [],
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
    return new AddReaderWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;