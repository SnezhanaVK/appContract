import http from 'k6/http';
import { check } from 'k6';

// Конфигурация теста
export const options = {
  stages: [
    { duration: '30s', target: 50 },  
    { duration: '1m', target: 50 },    
    { duration: '30s', target: 0 },    
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],    
    http_req_duration: ['p(95)<500'],  
  },
};

export default function () {
  const url = 'http://localhost:8080/api/contracts';
  
  const res = http.get(url);
  
  // Проверяем, что статус 200
  check(res, {
    'Status is 200': (r) => r.status === 200,
  });
}