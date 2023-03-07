'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class StoreKeyWorkload extends WorkloadModuleBase {
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
        var b64key = "MIIBHzANBgkqhkiG9w0BAQEFAAOCAQwAMIIBBwKB/yr1Isw6i2l0TV41ldVEe1RRIJeL4DnGusWdYRDojV60DVmjnLytbyPdJh00pVU+ZfzTq+dAToD3fGA5mhk3Pt9vhnygN1uFjxJB7x0jM9q1mbxhd488Yzk2lm+RUGjxa+cJITN66Z+1wxMiftZQBb15fPvTOv8tl+JiR83bjuTFLjMIkiCHwH7Z2XwueNQY8qpX4uzhDLCjdjG2w1oLiuO9uCX1tbyfpkxe0599TtNJwvVoisaF2ZiMo73nJIoSvg+B+pFoe9uByzp239fiMDNxnOQhTKhGTwXG6Z8lt1a3i+zn8GiED9nxQcWqC9oqyrHTnLn3/Lc5Nsb/xcAEowIDAQAB"
    	let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'StoreUserRSAPubkey',
            contractArguments: [obscuredName, b64key],
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
    return new StoreKeyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;