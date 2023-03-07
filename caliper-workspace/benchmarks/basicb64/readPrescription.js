'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class ReadPrescriptionWorkload extends WorkloadModuleBase {
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
            contractFunction: 'ReadPrescription',
            contractArguments: ["0"],
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
    return new ReadPrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;