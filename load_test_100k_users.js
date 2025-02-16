import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

const users = new SharedArray('users', function () {
    return open('users.txt')
        .split('\n')
        .map(line => {
            let [username, password] = line.split(':');
            return { username, password };
        });
});

export let options = {
    vus: 500,
    duration: '60s',
    rps: 1000,
    thresholds: {
        http_req_duration: ['p(95)<50'],
        checks: ['rate>0.9999'],
    },
};

export default function () {
    let user = users[Math.floor(Math.random() * users.length)];

    let authRes = http.post('http://localhost:8080/api/auth', JSON.stringify({
        username: user.username, password: user.password
    }), { headers: { 'Content-Type': 'application/json' } });

    check(authRes, { 'Auth success': (r) => r.status === 200 });

    if (authRes.status !== 200) return;

    let token = JSON.parse(authRes.body).token;
    let headers = { headers: { 'Authorization': `Bearer ${token}` } };

    let infoRes = http.get('http://localhost:8080/api/info', headers);
    check(infoRes, { 'Info success': (r) => [200, 401].includes(r.status) });

    let buyRes = http.get('http://localhost:8080/api/buy/pen', headers);
    check(buyRes, { 'Buy success': (r) => [200, 400].includes(r.status) });

	let user2 = users[Math.floor(Math.random() * users.length)];

    let sendCoinRes = http.post('http://localhost:8080/api/sendCoin', JSON.stringify({
        toUser: user2.username, amount: 1
    }), { headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` } });

    check(sendCoinRes, { 'SendCoin success': (r) => [200, 400].includes(r.status) });

    sleep(1);
}
