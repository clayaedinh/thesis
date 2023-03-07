'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetKeyWorkload extends WorkloadModuleBase {
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
        var obscuredName = "5a6a59ea61c3e616cbe67281399d1991cf5efe32d253203698ad63fea1daf9ae"
    	let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'RetrieveUserRSAPubkey',
            contractArguments: [obscuredName],
            timeout: 30,
            readOnly: true
        };
        await this.sutAdapter.sendRequests(args);
    }
}
 

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new GetKeyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;