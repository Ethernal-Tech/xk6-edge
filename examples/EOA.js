import eth from 'k6/x/ethereum';
import exec from 'k6/execution';

export const options = {
  iterations: 21,
  VUS: 21,
};

export function setup() {
    eth.Premine()
}

export default function (data) {

}
