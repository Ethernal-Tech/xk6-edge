import eth from 'k6/x/ethereum';
import exec from 'k6/execution';

export const options = {
    vus: 6,
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            duration: '1m',
            preAllocatedVUs: 6,
            rate: 500,
            timeUnit: '1s'
        }
    }
};

export function setup() {
    return eth.Premine()
}

var client
var gasPrice

export default function (data) {
    if (client == null) {
        client = new eth.Client(data[exec.vu.idInTest - 1].PrivateKey)
    }

    if (exec.vu.iterationInInstance == 0 || exec.vu.iterationInInstance % 500 == 0) {
        gasPrice = client.gasPrice()
    }

    const tx = {
        to: "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF",
        value: 100,
        gas_price: gasPrice * 1.2,
    };

    client.sendRawTransaction(tx)
}
