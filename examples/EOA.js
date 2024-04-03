import eth from 'k6/x/ethereum';
import exec from 'k6/execution';

export const options = {
    duration: '1m',
    VUS: 2,
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
