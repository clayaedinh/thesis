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
        var b64Share = "DoTid34kFhdyWHewc6WwgIxdW3SSwFvaNcN8fNRquh/Gbvy7PorsJorkWt1tMIdPjznOzX6AxRANRY4nwyelEVg4b6CJDh2/qi1VuRz1hS20fQAIcOCFNPJoUZPhpuHBXbjtcd+yeq9ooQSbOXK92Y6ipsj38xtBfNqcGCSzQB6nwwiGWo4bHgQFG/TdMHKShpJf6E4q4JX8H1UikHt+HCB+2UjrXfZ0L9t6qvrfA7VvffKMl+thAdJ4lyos2TNYNn0tWsABOIHdLCiOu20/tSbzY+WNDCYH0Xmxqa41pnbTQ/nYNap6Vst41KwdLt/Nv8utSKx+3io2HQqwyor8"
    	var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
        var obscuredName = "5a6a59ea61c3e616cbe67281399d1991cf5efe32d253203698ad63fea1daf9ae"
        let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'SharePrescription',
            contractArguments: [pid.toString(), obscuredName, b64Share],
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
    return new SharePrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;