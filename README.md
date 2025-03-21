# Solana Testnet Faucet

A self-hosted faucet for distributing Solana testnet tokens with rate limiting, recaptcha protection, and claim history tracking.

## Features

- Request testnet SOL with a simple user interface
- IP and wallet-based rate limiting (once per day per wallet)
- reCAPTCHA protection against bots
- Transaction history tracking
- Multiple deployment options (local, Nginx, or Cloudflare Tunnel)

## Tech Stack

- **Backend**: Golang
- **Frontend**: Vue 3 with Vite and Tailwind CSS 3
- **Database**: SQLite
- **Deployment**: Docker & Docker Compose with multiple configuration options

## Prerequisites

- Docker and Docker Compose
- A Solana wallet with testnet SOL (for funding the faucet)
- reCAPTCHA v2 site and secret keys (optional)
- Cloudflare account for Tunnel (optional for production deployment)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/solana-faucet.git
   cd solana-faucet
   ```

2. Set up your Solana wallet:
   ```bash
   # Create a new wallet.json file or copy an existing one
   # Format should be:
   # {
   #   "privateKey": "base64_encoded_private_key",
   #   "publicKey": "your_public_key"
   # }
   ```

3. Configure the application in `config.json` or use environment variables.

4. Create a `.env` file based on the provided `.env.example`:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

## Deployment Options

The project includes three Docker Compose configurations for different deployment scenarios, all conveniently accessible via the Makefile:

### Quick Start with Makefile

```bash
# Show all available commands
make help

# Run local development setup
make local

# Run development setup with Nginx
make dev

# Run production setup with Cloudflare Tunnel
make prod

# Stop all containers
make down

# Remove all containers, networks, and volumes
make clean
```

### 1. Local Development (docker-compose.local.yml)

Simplest setup with direct port exposure for local development:

```bash
make local
```
or
```bash
docker-compose -f docker-compose.local.yml up -d
```

Access:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080

### 2. Nginx Proxy (docker-compose.yml)

Uses Nginx as a reverse proxy with optional SSL support for testing or simple production environments:

```bash
make dev
```
or
```bash
docker-compose up -d
```

Access:
- Frontend: http://localhost:80 (or https://localhost:443 if SSL is configured)
- Backend API: Proxied through the same domain via paths

To configure SSL:
1. Place your SSL certificates in the `nginx/certs` directory
2. Uncomment the HTTPS server section in `nginx/conf.d/default.conf`

### 3. Cloudflare Tunnel (docker-compose.prod.yml)

Secure production deployment using Cloudflare Tunnel:

```bash
make prod
```
or
```bash
docker-compose -f docker-compose.prod.yml up -d
```

Access:
- Via your configured domain(s) in Cloudflare Dashboard

## Cloudflare Tunnel Setup

Cloudflare Tunnel provides several advantages over traditional port forwarding or direct internet exposure:

1. **No Public IP or Port Forwarding Required**: Eliminates the need to expose ports to the internet
2. **TLS Encryption**: Automatic HTTPS with Cloudflare-managed certificates
3. **DDoS Protection**: Built-in Cloudflare protection against attacks
4. **Access Controls**: Can be integrated with Cloudflare Access for authentication
5. **DNS Management**: Automatic DNS record management

### Tunnel Setup Steps

1. Install the Cloudflare CLI locally:
   ```bash
   # On Linux
   curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
   sudo dpkg -i cloudflared.deb
   
   # Or with Docker
   docker pull cloudflare/cloudflared:latest
   ```

2. Authenticate with Cloudflare:
   ```bash
   cloudflared login
   ```

3. Create a tunnel:
   ```bash
   cloudflared tunnel create solana-faucet
   ```

4. Get your tunnel token and add it to your `.env` file:
   ```
   CLOUDFLARE_TUNNEL_TOKEN=your-tunnel-token
   ```

5. Configure DNS records in the Cloudflare dashboard to point to your tunnel

6. Configure your tunnel:
   - Option 1: Using Cloudflare Zero Trust dashboard (Recommended)
     - Configure public hostnames in the Cloudflare dashboard for your frontend and API

   - Option 2: Using config file
     - Copy the example config file: `cp cloudflared/config.yml.example cloudflared/config.yml`
     - Edit with your specific tunnel ID and domain names

For more detailed instructions, refer to the [Cloudflare Tunnel documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/).

## Environment Variables

You can configure the application using environment variables in the `.env` file:

- `FAUCET_SERVER_ADDRESS`: Server address (default: 0.0.0.0)
- `FAUCET_SERVER_PORT`: Server port (default: 8080)
- `FAUCET_DB_PATH`: Path to SQLite database file (default: /data/faucet.db)
- `FAUCET_SOLANA_RPC_URL`: Solana RPC URL (default: https://api.testnet.solana.com)
- `FAUCET_WALLET_PATH`: Path to wallet file (default: /data/wallet.json)
- `RECAPTCHA_SECRET_KEY`: reCAPTCHA secret key
- `RECAPTCHA_SITE_KEY`: reCAPTCHA site key
- `CLOUDFLARE_TUNNEL_TOKEN`: Cloudflare Tunnel token (for production deployment)

## Security Considerations

- In production, set proper reCAPTCHA keys
- Secure your wallet private key
- Consider running the faucet only on testnets
- Adjust rate limits in config.json based on your needs
- When using Cloudflare Tunnel, ensure proper authentication is set up if needed

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details. 