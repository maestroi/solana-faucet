const config = {
    development: {
        apiBaseUrl: 'http://localhost:8080'
    },
    production: {
        apiBaseUrl: 'https://api-sol-faucet.maestroi.cc'
    }
}

const env = import.meta.env.MODE || 'development'
export const apiBaseUrl = config[env].apiBaseUrl 