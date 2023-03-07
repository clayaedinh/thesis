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
        var b64Setfill = "/5f/gQMBAQxQcmVzY3JpcHRpb24B/4IAAQgBBUJyYW5kAQwAAQZEb3NhZ2UBDAABC1BhdGllbnROYW1lAQwAAQ5QYXRpZW50QWRkcmVzcwEMAAEOUHJlc2NyaWJlck5hbWUBDAABDFByZXNjcmliZXJObwEGAAELUGllY2VzVG90YWwBBgABDFBpZWNlc0ZpbGxlZAEGAAAABf+CCGQA"
    	var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
        let args = {
            contractId: 'basicb64',
            contractVersion: '1.0',
            contractFunction: 'SetfillPrescription',
            contractArguments: [pid.toString(), b64Setfill],
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
    return new SetfillPrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;