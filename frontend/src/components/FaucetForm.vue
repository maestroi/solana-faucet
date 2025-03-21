<template>
  <div class="w-full flex justify-center">
    <div class="w-full max-w-4xl bg-[#141414] overflow-hidden shadow-lg rounded-lg">
      <div class="p-8 bg-[#141414] border-b border-[#405045]">
        <div class="mb-8">
          <FaucetBalance />
        </div>
        
        <!-- Info Section -->
        <div class="mb-8 p-6 bg-[#1a1a1a] rounded-lg border border-[#405045]">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-[#00ffa3]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-medium text-[#00ffa3]">Important Information</h3>
              <div class="mt-2 text-[#405045]">
                <p class="text-base">• Maximum request: 0.1 SOL per 24 hours</p>
                <p class="text-base">• Requests are limited to one per wallet address</p>
                <p class="text-base">• Please ensure you're using a valid Solana testnet wallet address</p>
              </div>
            </div>
          </div>
        </div>
        
        <form @submit.prevent="requestFunds" class="space-y-8">
          <div>
            <label for="wallet" class="block text-sm font-medium text-[#405045]">Solana Wallet Address</label>
            <div class="mt-1 relative rounded-md shadow-sm">
              <input 
                id="wallet" 
                v-model="walletAddress" 
                type="text" 
                required
                class="block w-full px-4 py-3 bg-[#141414] border-2 border-[#405045] rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-[#00ffa3] focus:border-[#00ffa3] text-lg text-[#405045]"
                placeholder="Enter your Solana wallet address"
              />
            </div>
            <p v-if="validationError" class="mt-2 text-sm text-red-400">
              {{ validationError }}
            </p>
          </div>

          <div class="flex justify-center">
            <div
              class="cf-turnstile"
              data-sitekey="0x4AAAAAABB20aW1diwbHE5q"
              data-callback="onTurnstileSuccess"
            ></div>
          </div>
          
          <div class="flex items-center justify-between">
            <button 
              type="submit" 
              class="w-full flex justify-center py-3 px-6 border border-transparent rounded-md shadow-sm text-lg font-medium text-gray-900 bg-[#00ffa3] hover:bg-[#00e694] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#00ffa3] disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
              :disabled="isLoading || !turnstileToken"
            >
              <span v-if="isLoading" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-900" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Processing...
              </span>
              <span v-else>Request SOL</span>
            </button>
          </div>
        </form>
        
        <!-- Status Message -->
        <div v-if="statusMessage" class="mt-8 p-6 rounded-md" :class="statusClass">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg v-if="statusType === 'success'" class="h-6 w-6 text-[#00ffa3]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              <svg v-else-if="statusType === 'error'" class="h-6 w-6 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-lg font-medium" v-html="statusMessage"></p>
            </div>
          </div>
        </div>
        
        <!-- Next Claim Timer -->
        <div v-if="nextClaimTime" class="mt-8 p-6 bg-[#1b4e3f] bg-opacity-5 rounded-md">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-[#1b4e3f]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-medium text-[#1b4e3f]">Next Claim Available</h3>
              <p class="mt-1 text-lg text-[#1b4e3f]">{{ nextClaimTime }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Recent Transactions Section -->
      <div class="p-8 bg-[#1a1a1a] border-t border-[#405045]">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-2xl font-bold text-[#00ffa3]">Recent Transactions</h3>
          <button @click="fetchTransactions" class="text-[#405045] hover:text-[#00ffa3] transition-colors duration-200">
            <span class="flex items-center">
              <svg class="h-5 w-5 mr-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
              </svg>
              Refresh
            </span>
          </button>
        </div>
        <div v-if="transactions.length > 0" class="overflow-x-auto">
          <table class="min-w-full divide-y divide-[#405045]">
            <thead class="bg-[#1a1a1a]">
              <tr>
                <th class="px-6 py-4 text-left text-sm font-semibold text-[#405045] uppercase tracking-wider">Time</th>
                <th class="px-6 py-4 text-left text-sm font-semibold text-[#405045] uppercase tracking-wider">Wallet</th>
                <th class="px-6 py-4 text-left text-sm font-semibold text-[#405045] uppercase tracking-wider">Amount</th>
                <th class="px-6 py-4 text-left text-sm font-semibold text-[#405045] uppercase tracking-wider">Status</th>
              </tr>
            </thead>
            <tbody class="bg-[#1a1a1a] divide-y divide-[#405045]">
              <tr v-for="tx in transactions" :key="tx?.id || Math.random()" class="hover:bg-[#222222] transition-colors duration-200">
                <td class="px-6 py-4 whitespace-nowrap text-base text-[#405045]">{{ formatDate(tx?.timestamp) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-base text-[#405045]">
                  <a :href="`https://explorer.solana.com/address/${tx?.walletAddress || ''}?cluster=testnet`" target="_blank" class="text-[#00ffa3] hover:text-[#00e694] transition-colors duration-200">
                    {{ shortAddress(tx?.walletAddress) }}
                  </a>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-base text-[#405045]">{{ tx?.amount || 0 }} SOL</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full" :class="getStatusClass(tx?.status)">
                    {{ tx?.status || 'Unknown' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="text-center py-12">
          <svg class="mx-auto h-16 w-16 text-[#405045]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <p class="mt-4 text-lg text-[#405045]">No recent transactions</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { apiBaseUrl } from '../config'
import FaucetBalance from './FaucetBalance.vue'

// State
const walletAddress = ref('')
const validationError = ref('')
const isLoading = ref(false)
const statusMessage = ref('')
const statusType = ref('') // 'success', 'error', or 'info'
const nextClaimTime = ref('')
const transactions = ref([])
const turnstileToken = ref('')

// Computed properties
const statusClass = computed(() => {
  if (statusType.value === 'success') return 'bg-[#00ffa3] bg-opacity-10 text-[#00ffa3]'
  if (statusType.value === 'error') return 'bg-red-900 text-red-200'
  return 'bg-blue-900 text-blue-200'
})

// Methods
const validateWalletAddress = (address) => {
  if (!address) return 'Wallet address is required'
  if (address.length < 32) return 'Invalid Solana wallet address'
  return null
}

// Turnstile callback
const onTurnstileSuccess = (token) => {
  turnstileToken.value = token
}

// Add the callback to window so Turnstile can call it
window.onTurnstileSuccess = onTurnstileSuccess

const requestFunds = async () => {
  try {
    // Basic validation
    const error = validateWalletAddress(walletAddress.value)
    if (error) {
      validationError.value = error
      return
    }

    if (!turnstileToken.value) {
      validationError.value = 'Please complete the verification'
      return
    }

    validationError.value = ''
    isLoading.value = true
    statusMessage.value = ''
    statusType.value = ''

    // Log the request payload for debugging
    const payload = {
      wallet_address: walletAddress.value,
      cf_turnstile_response: turnstileToken.value
    }
    console.log('Sending request with payload:', payload)

    const response = await axios.post(`${apiBaseUrl}/api/request-funds`, payload, {
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (response.data.success) {
      statusType.value = 'success'
      statusMessage.value = `Successfully sent ${response.data.amount} SOL. Transaction: `
      + `<a href="https://explorer.solana.com/tx/${response.data.transaction_hash}?cluster=testnet" target="_blank" class="text-[#1b4e3f] hover:text-[#00ffa3] underline">`
      + `${response.data.transaction_hash.slice(0, 8)}...${response.data.transaction_hash.slice(-8)}</a>`
      // Reset form
      walletAddress.value = ''
      turnstileToken.value = ''
      // Reset Turnstile
      window.turnstile.reset()
      // Fetch updated transactions
      await fetchTransactions()
    } else {
      statusType.value = 'error'
      statusMessage.value = response.data.error || 'Failed to send SOL'
    }
  } catch (error) {
    console.error('Error requesting funds:', error)
    statusType.value = 'error'
    
    // Handle different error responses
    if (error.response?.data) {
      // If the error response has a data object with an error message
      statusMessage.value = error.response.data.error || error.response.data.message || 'Request failed'
      
      // If there's a next claim time, update the timer
      if (error.response.data.nextClaimTime) {
        const nextTime = new Date(error.response.data.nextClaimTime * 1000)
        nextClaimTime.value = formatDate(nextTime)
      }
    } else if (error.response?.status === 429) {
      // Handle rate limit error
      statusMessage.value = 'You have reached the request limit. Please try again later.'
    } else if (error.response?.status === 400) {
      statusMessage.value = 'Invalid request. Please check your wallet address and try again.'
    } else {
      statusMessage.value = 'An error occurred while processing your request'
    }
  } finally {
    isLoading.value = false
  }
}

const fetchTransactions = async () => {
  try {
    const response = await axios.get(`${apiBaseUrl}/api/transactions`)
    // Handle the correct response format where transactions are nested
    transactions.value = response.data?.transactions || []
  } catch (error) {
    console.error('Error fetching transactions:', error)
    transactions.value = [] // Set to empty array on error
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  // For next claim time, show relative time
  if (d > new Date()) {
    const diff = d - new Date()
    const hours = Math.floor(diff / 3600000)
    const minutes = Math.floor((diff % 3600000) / 60000)
    const seconds = Math.floor((diff % 60000) / 1000)
    
    // If more than 2 hours, show date and time
    if (hours >= 2) {
      return d.toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true
      })
    }
    
    // If more than 1 hour, show hours and minutes
    if (hours > 0) {
      if (minutes > 0) {
        return `${hours} hour${hours > 1 ? 's' : ''} and ${minutes} minute${minutes > 1 ? 's' : ''}`
      }
      return `${hours} hour${hours > 1 ? 's' : ''}`
    }
    
    // If more than 1 minute, show minutes and seconds
    if (minutes > 0) {
      if (seconds > 0) {
        return `${minutes} minute${minutes > 1 ? 's' : ''} and ${seconds} second${seconds !== 1 ? 's' : ''}`
      }
      return `${minutes} minute${minutes > 1 ? 's' : ''}`
    }
    
    // Less than 1 minute, show seconds
    return `${seconds} second${seconds !== 1 ? 's' : ''}`
  }

  // For transaction history, show absolute time
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  })
}

const shortAddress = (address) => {
  if (!address) return ''
  return `${address.slice(0, 4)}...${address.slice(-4)}`
}

const getStatusClass = (status) => {
  switch (status?.toLowerCase()) {
    case 'completed':
      return 'bg-[#00ffa3] bg-opacity-10 text-[#00ffa3]'
    case 'failed':
      return 'bg-red-900 text-red-200'
    default:
      return 'bg-gray-700 text-gray-300'
  }
}

onMounted(async () => {
  // Fetch initial transactions
  fetchTransactions()
})
</script> 