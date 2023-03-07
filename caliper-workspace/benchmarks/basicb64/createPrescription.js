'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class CreatePrescriptionWorkload extends WorkloadModuleBase {
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
        var b64Create = "/5f/gQMBAQxQcmVzY3JpcHRpb24B/4IAAQgBBUJyYW5kAQwAAQZEb3NhZ2UBDAABC1BhdGllbnROYW1lAQwAAQ5QYXRpZW50QWRkcmVzcwEMAAEOUHJlc2NyaWJlck5hbWUBDAABDFByZXNjcmliZXJObwEGAAELUGllY2VzVG90YWwBBgABDFBpZWNlc0ZpbGxlZAEGAAAAA/+CAA=="
    	let args = {
            contractId: 'basicb64',
            contractVersion: '1.0',
            contractFunction: 'CreatePrescription',
            contractArguments: [b64Create],
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
    return new CreatePrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;