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
        var b64update = "Dv+DBAEC/4QAAQwBDAAA/gLw/4QAAUA1YTZhNTllYTYxYzNlNjE2Y2JlNjcyODEzOTlkMTk5MWNmNWVmZTMyZDI1MzIwMzY5OGFkNjNmZWExZGFmOWFl/gKoSC8wNkdkTlhiMGNEdi9NQ3pzOGFRd1NYMllzT1RlSGJ1V0pmM0h3bGFpRDV1MWpvZWI4eGw1czdKTXNQQnd4T1creGMvMFphbUhXbEoxTHQwMXdDODg1OFV0Z3FreXk4RGhwNGJBMS9xOWlGRnNrVWZuNUZYelJqRlU2VUFycWd1UFlJSHVBd3Zsb3JyT2w5ZzVZbTRUSHVRYzhtMTdXL1cvRlJrcEFnTDJzZFByS3Q1a2hySUxNbEl3VVExUkZzcUIydC9mRXE2TnhHVGxqcWIyTHZJd0xRL1JvZTNRVXBJejVrdWR2NmlVRFpxYllQNytTSnEybkpFSSt0VnF5aXNMcUx1bUcvb0hYZzFOMllYMVIzc3ZJTXYrbkZHUGgxdUdUNzZDeVBXakZwQi9ZSDIvNzBQanpYKzFnSmh5TXpUUWlQVFhadWU2UWt6QjdMNEUzZElqczdVbVlxMUtRd0JKN0cxVnBvTDc1VmJXUE9INXhUL0R4dE5uOVdiRjRTUlcwVE5nZjNxeTZJVVpJRmVXaVpyQzhXdlNUUEsxbTNsUU5sUENjTjJCbjc4NEQ0a3o3N2NQMHQ0c0tkcXl6ZUZMbThycHNVTDF6ODNBYUJNM2MwbWJHL2Z4YTJKMzNac042aHk5TFpLRnF6bnpQZ1Y1ZWRIM3RaMkQyV0czVVhPMVcrL1gyZ0VuWmEzMUJJMHNLaVI4UXdNUG9XS1ZjYWxIN1VXUkZ6cW0vL1N5Y1hUbUxZMFk0UWF0cHl2aE1qUU54RjhBWVd2bUFsQ2RvTWk4MjJTVUx4WXZNOW1Oc2tEVExsUmh3VGxzT3hza2lhYlhXK3ZySHNEWExxb0FRa2RNR1RiU0dmamJ0eWFEZk5NZGdDYkVFN0ZWZ00xZVBBYmFnUWxZYkM="
    	var pid = this.txIndex + (this.workerIndex * this.totalWorkers);
        let args = {
            contractId: 'rsa',
            contractVersion: '1.0',
            contractFunction: 'UpdatePrescription',
            contractArguments: [pid.toString(), b64update],
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
    return new UpdatePrescriptionWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;