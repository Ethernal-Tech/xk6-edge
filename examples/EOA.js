import eth from 'k6/x/ethereum';
import exec from 'k6/execution';

export const options = {
    iterations: 10,
    VUS: 10,
};

export function setup() {
    return eth.Premine()
}

var client

export default function (data) {
    if (client == null) {
        client = new eth.Client(data[exec.vu.idInTest - 1].PrivateKey)
    }
}
