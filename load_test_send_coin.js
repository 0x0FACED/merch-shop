import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 500,  
  duration: '60s', 
  rps: 1000,
  thresholds: {
    http_req_duration: ['p(95)<50'],
    http_req_failed: ['rate<0.0001'],
    checks: ['rate>0.9999'],
  },
};

export function setup() {
  let authLoader = http.post('http://localhost:8080/api/auth', JSON.stringify({ 
    username: "loadtest2user1", password: "loadtest2user1" 
  }), { headers: { 'Content-Type': 'application/json' } });

  check(authLoader, { 'Auth loadtest2user1 success': (r) => r.status === 200 });
  let loaderToken = JSON.parse(authLoader.body).token;

  let authUser = http.post('http://localhost:8080/api/auth', JSON.stringify({ 
    username: "loadtest2user2", password: "loadtest2user2" 
  }), { headers: { 'Content-Type': 'application/json' } });

  check(authUser, { 'Auth loadtest2user2 success': (r) => r.status === 200 });
  let userToken = JSON.parse(authUser.body).token;

  return { loaderToken, userToken };
}

export default function (data) {
  let sender, recipient;

  if (__ITER % 2 === 0) {
    sender = data.loaderToken;
    recipient = "loadtest2user2";
  } else {
    sender = data.userToken;
    recipient = "loadtest2user1";
  }

  let headers = { headers: { 
    'Authorization': `Bearer ${sender}`,
    'Content-Type': 'application/json'
  }};

  let sendCoinRes = http.post('http://localhost:8080/api/sendCoin', JSON.stringify({
    toUser: recipient, amount: 1
  }), headers);

  check(sendCoinRes, { 
    'SendCoin success': (r) => [200].includes(r.status) 
  });

  sleep(1);
}
