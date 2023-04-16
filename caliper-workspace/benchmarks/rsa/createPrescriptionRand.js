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
        var b64Create = "HJHGAxq584fF3avQJWS57j0aDuEKCn4rRa37MyahQKPiuPccUUkNg0R+EPF4eo52kXhF8VMxzTOT0/wGH12Dg8g+tafBjwqgQRZwFVBBg+n9vd276+GRbfKAm2/r2IrdTdQTiGqH6eMLLY0hSHv/7WXkT161+v0ajBv7TEAFyTv55AVV7ORSFm/GMk8NjGsHLTSq+MP2RsYkjM7DdfVlCLtaikITh3Ve1pMcB1k/a8jONTG0fu1yHaWvX/7E/yqPIBya/puh8bsju0RqybnsB0TC0GamJd4T0CE6dtZxC47sJbTfBQTxWUfL/gwlkNHq1zGbrv4kedkm3+x7Bej+"
    	let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'CreatePrescriptionRand',
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