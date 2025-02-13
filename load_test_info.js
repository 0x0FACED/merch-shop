import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 1000,  
  duration: '60s', 
  rps: 1000,
  thresholds: {
	http_req_duration: ['p(95)<50'],
	checks: ['rate>0.9999'],
	http_req_failed: ['rate<0.0001'],
  },
};

let token = '';

export function setup() {
  let authRes = http.post('http://localhost:8080/api/auth', JSON.stringify({ 
	username: "loader1", password: "loader1" 
  }), { headers: { 'Content-Type': 'application/json' } });

  check(authRes, { 'Auth success': (r) => r.status === 200 });
  token = JSON.parse(authRes.body).token;

  return { token };
}

export default function (data) {
  let headers = { headers: { 'Authorization': `Bearer ${data.token}` } };

  let infoRes = http.get('http://localhost:8080/api/info', headers);
  check(infoRes, { 
	'Info success': (r) => [200, 401].includes(r.status) 
  });

  sleep(1);
}
