import eth from 'k6/x/ethereum';
import exec from 'k6/execution';

export const options = {
    vus: 5,
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            duration: '1m',
            preAllocatedVUs: 5,
            rate: 300,
            timeUnit: '1s'
        }
    }
};

export function setup() {
    return eth.Premine()
}

var client

export default function (data) {
    if (client == null) {
        client = new eth.Client(data[exec.vu.idInTest - 1].PrivateKey)
    }

    const tx = {
        to: "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF",
        value: 100,
        gas_price: client.gasPrice() * 1.2,
    };

    client.sendRawTransaction(tx)
}
