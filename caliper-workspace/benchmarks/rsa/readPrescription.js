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
        var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
    	let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'ReadPrescription',
            contractArguments: [pid.toString()],
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