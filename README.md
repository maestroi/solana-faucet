# Solana Testnet Faucet

A simple faucet for distributing testnet SOL to developers. Built with Go and Vue.js.

## Features

- Request testnet SOL with a simple web interface
- Rate limiting and cooldown periods to prevent abuse
- Cloudflare Turnstile protection against bots
- Transaction history tracking
- Real-time faucet balance display
- Dark mode UI with modern design

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Docker and Docker Compose
- A Solana wallet with testnet SOL (for the faucet)

## Environment Variables

### Backend

```env
# Server Configuration
FAUCET_SERVER_ADDRESS=0.0.0.0
FAUCET_SERVER_PORT=8080

# Database Configuration
FAUCET_DB_PATH=/app/data/faucet.db

# Solana Configuration
FAUCET_SOLANA_RPC_URL=https://api.testnet.solana.com
FAUCET_WALLET_PATH=/app/data/wallet.json
FAUCET_AMOUNT_PER_REQUEST=0.1
FAUCET_NETWORK_TYPE=testnet
FAUCET_TRANSACTION_TIMEOUT=30

# Security Configuration
FAUCET_TURNSTILE_SECRET=your-turnstile-secret-key
FAUCET_TURNSTILE_SITE=your-turnstile-site-key
FAUCET_RATE_LIMIT_REQUESTS=5
FAUCET_RATE_LIMIT_DURATION=60
FAUCET_CLAIM_COOLDOWN=86400  # 24 hours in seconds

# CORS Configuration
FAUCET_CORS_ALLOWED_ORIGINS=http://localhost:3000,https://faucet.solana.com
```

### Frontend

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080  # Development
VITE_API_BASE_URL=https://api-sol-faucet.maestroi.cc  # Production
```

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/maestroi/solana-faucet.git
   cd solana-faucet
   ```

2. Set up the Solana wallet:
   ```bash
   # Create a new wallet
   solana-keygen new -o wallet.json --outfile wallet.json
   
   # Get the public key
   solana-keygen pubkey wallet.json
   
   # Request testnet SOL for the faucet wallet
   solana airdrop 10 <YOUR_WALLET_PUBKEY> --url https://api.testnet.solana.com
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Start the development environment:
   ```bash
   docker-compose up -d
   ```

5. Access the faucet at http://localhost:3000

## Production Deployment

1. Set up your production environment variables:
   ```bash
   cp .env.example .env.prod
   # Edit .env.prod with your production configuration
   ```

2. Build and deploy using Docker Compose:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. Access your faucet at your configured domain

## Security Considerations

- The faucet uses Cloudflare Turnstile for bot protection
- Rate limiting is implemented to prevent abuse
- A cooldown period is enforced between requests
- CORS is configured to allow only specific origins
- The faucet wallet should be kept secure and have limited funds

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 