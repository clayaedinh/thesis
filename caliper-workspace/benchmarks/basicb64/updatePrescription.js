'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class UpdatePrescriptionWorkload extends WorkloadModuleBase {
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
        var b64Update = "/5f/gQMBAQxQcmVzY3JpcHRpb24B/4IAAQgBBUJyYW5kAQwAAQZEb3NhZ2UBDAABC1BhdGllbnROYW1lAQwAAQ5QYXRpZW50QWRkcmVzcwEMAAEOUHJlc2NyaWJlck5hbWUBDAABDFByZXNjcmliZXJObwEGAAELUGllY2VzVG90YWwBBgABDFBpZWNlc0ZpbGxlZAEGAAAAS/+CAQpEUlVHIEJSQU5EAQtEUlVHIERPU0FHRQEMUEFUSUVOVCBOQU1FAQxQQVRJRU5UIEFERFIBClBSRVNDIE5BTUUB/RLWhwFkAA=="
    	let args = {
            contractId: 'basicb64',
            contractVersion: '1.0',
            contractFunction: 'UpdatePrescription',
            contractArguments: ["0", b64Update],
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
    return new UpdatePrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;