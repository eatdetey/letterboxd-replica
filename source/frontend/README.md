# Vue 3 + TypeScript + Vite

## Development

1. Install dependencies:
   ```bash
   npm install
   ```
2. Create local env file:
   ```bash
   cp .env.example .env.development
   ```
3. In `.env.development` set backend gateway address:
   ```env
   VITE_API_GATEWAY_URL=http://localhost:8080
   ```
4. Run frontend in dev mode:
   ```bash
   npm run dev
   ```

In dev mode frontend calls `/api/*`, and Vite proxies these requests to `VITE_API_GATEWAY_URL`.

## Production (nginx)

Frontend still calls `/api/*`, and nginx proxies these requests to api-gateway via `API_GATEWAY_URL`.
Default value inside container:

```env
API_GATEWAY_URL=http://api-gateway:8080
```

Build:

```bash
npm run build
```
