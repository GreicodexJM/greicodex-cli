import request from 'supertest';
import express from 'express';
import serverless from 'serverless-http';

const app = express();
app.get('/', (req, res) => {
  res.send({ "Hello": "World" });
});
const handler = serverless(app);


describe('GET /', () => {
  it('responds with json', async () => {
    const response = await request(app).get('/');
    expect(response.statusCode).toBe(200);
    expect(response.body).toEqual({ Hello: 'World' });
  });
});
